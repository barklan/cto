package logserver

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/barklan/cto/pkg/logserver/querying"
	"github.com/barklan/cto/pkg/storage"
)

const (
	flagError     = "err"
	flagCritical  = "crit"
	flagFatal     = "fatal"
	flagEmergency = "emerg"
	flagNone      = "none"
)

type RawLogRecord map[string]interface{}

// TODO add more fields (may be specific to service - handle default values gracefully)
// TODO is this used anywhere?
type LogMetadata struct {
	Hostname  string `json:"fluentd_hostname"`
	Service   string `json:"service_name"`
	Timestamp string `json:"fluentd_time"`
	Flag      string
}

func constructMetadata(record RawLogRecord) *LogMetadata {
	logMetadata := LogMetadata{}

	key := "fluentd_hostname"
	if val, ok := record[key]; ok {
		if v, okType := val.(string); okType {
			logMetadata.Hostname = v
		}
	} else {
		log.Printf("Record missing required field: %v", key)
	}

	key = "service_name"
	if val, ok := record[key]; ok {
		if v, okType := val.(string); okType {
			logMetadata.Service = v
		}
	} else {
		log.Printf("Record missing required field: %v", key)
	}

	key = "fluentd_time"
	if val, ok := record[key]; ok {
		if v, okType := val.(string); okType {
			logMetadata.Timestamp = v
		}
	} else {
		log.Printf("Record missing required field: %v", key)
	}

	return &logMetadata
}

// TODO benchmark this.
func processLogRecord(
	data *storage.Data,
	record RawLogRecord,
	sessDataMap map[string]*SessionData,
	// TODO google: passing one-way channels to functions
	reportChan chan LogRecordReport,
) {
	logData := constructMetadata(record)

	projectName, ok := data.Config.EnvToProjectName[logData.Hostname]
	if !ok {
		log.Println("Did not find project associated with", logData.Hostname)
		return
	}

	// TODO it is heavy to do it on every log record - should decide randomly instead
	knownServices := querying.GetKnownServices(data, logData.Hostname)
	if _, ok := knownServices[logData.Hostname]; !ok {
		knownServices[logData.Service] = struct{}{}
		querying.SetKnownServices(data, logData.Hostname, knownServices)
	}

	flag := assignFlag(fmt.Sprint(record))
	logData.Flag = flag

	// Save log record
	badgerKey, _ := constructBadgerKey(data, logData, projectName)
	retentionDuration := time.Duration(data.Config.Internal.Log.RetentionHours) * time.Hour

	data.SetLog(badgerKey, record, retentionDuration)

	// log.Printf("Log added with key: %q", badgerKey)

	// Send message if error
	if logData.Flag != flagNone {
		projectSessData := sessDataMap[projectName]
		_ = handleErrorRecord(data, logData, record, badgerKey, projectSessData, projectName)
	}

	reportChan <- LogRecordReport{ProjectName: projectName}
}

func constructBadgerKey(
	data *storage.Data,
	logData *LogMetadata,
	projectName string,
) (string, error) {
	randID := strconv.FormatInt(rand.Int63(), 10)
	fluentdTimeArr := strings.Fields(logData.Timestamp)
	dateString, timeString := fluentdTimeArr[0], fluentdTimeArr[1]
	badgerKey := strings.Join(
		[]string{
			projectName,
			logData.Hostname,
			logData.Service,
			dateString,
			timeString,
			logData.Flag,
			randID,
		}, " ",
	)
	return badgerKey, nil
}

// TODO more flags!
// TODO flags for traefik
func assignFlag(str string) string {
	if strings.Contains(str, "ERROR") || strings.Contains(str, " [error] ") {
		return flagError
	}

	if strings.Contains(str, "CRITICAL") {
		return flagCritical
	}

	if strings.Contains(str, "FATAL") {
		return flagFatal
	}

	if strings.Contains(str, " [emerg] ") {
		return flagEmergency
	}

	return flagNone
}

// Deprecated hack - do not use.
func findTimeStringByRegex(str string) (string, bool) {
	r, _ := regexp.Compile(`(?:\D|\b)((?:[01]\d|2[0-3]):(?:[0-5]\d):(?:[0-5]\d))(?:\D|\b)`)
	timeStr := r.FindStringSubmatch(str)
	if len(timeStr) >= 2 {
		return timeStr[1], true
	}
	return "", false
}
