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
	flagError    = "err"
	flagCritical = "crit"
	flagNone     = "none"
)

type RawLogRecord map[string]interface{}

// TODO add more fields (may be specific to service - handle default values gracefully)
// TODO is this used anywhere?
type MetaLogData struct {
	Hostname    string `json:"fluentd_hostname"`
	ServiceName string `json:"service_name"`
	Timestamp   string `json:"fluentd_time"`
	LogString   string `json:"log"`
}

// TODO benchmark this.
func processLogRecord(
	data *storage.Data,
	record RawLogRecord,
	sessDataMap map[string]*SessionData,
	// TODO google: passing one-way channels to functions
	reportChan chan LogRecordReport,
) {
	// TODO logData should be struct instead
	logData := map[string]string{
		"service_name":     "",
		"fluentd_hostname": "",
		"fluentd_time":     "",
		"log":              "",
	}

	for key := range logData {
		if val, ok := record[key]; ok {
			logData[key] = fmt.Sprint(val)
		} else {
			log.Printf("Record missing required field: %v", key)
			return
		}
	}

	projectName, ok := data.Config.EnvToProjectName[logData["fluentd_hostname"]]
	if !ok {
		log.Println("Did not find project associated with", logData["fluentd_hostname"])
		return
	}

	// TODO it is heavy to do it on every log record - should decide randomly instead
	knownServices := querying.GetKnownServices(data, logData["fluentd_hostname"])
	if _, ok := knownServices[logData["fluentd_hostname"]]; !ok {
		knownServices[logData["service_name"]] = struct{}{}
		querying.SetKnownServices(data, logData["fluentd_hostname"], knownServices)
	}

	flag := assignFlag(logData["log"])
	logData["flag"] = flag

	// Save log record
	badgerKey, _ := constructBadgerKey(data, logData, projectName)
	retentionDuration := time.Duration(data.Config.Internal.Log.RetentionHours) * time.Hour

	delete(record, "log") // HACK - this should be done on fluentd side after it has parsed this to json

	data.SetLog(badgerKey, record, retentionDuration)

	// log.Printf("Log added with key: %q", badgerKey)

	// Send message if error
	if (logData["flag"] == flagError) || (logData["flag"] == flagCritical) {
		projectSessData := sessDataMap[projectName]
		_ = handleErrorRecordInteractive(data, logData, record, badgerKey, projectSessData, projectName)
	}

	reportChan <- LogRecordReport{ProjectName: projectName}
}

func constructBadgerKey(
	data *storage.Data,
	logData map[string]string,
	projectName string,
) (string, error) {
	randID := strconv.FormatInt(rand.Int63(), 10)
	fluentdTimeArr := strings.Fields(logData["fluentd_time"])
	dateString, timeString := fluentdTimeArr[0], fluentdTimeArr[1]
	badgerKey := strings.Join(
		[]string{
			projectName,
			logData["fluentd_hostname"],
			logData["service_name"],
			dateString,
			timeString,
			logData["flag"],
			randID,
		}, " ",
	)
	return badgerKey, nil
}

// TODO more flags!
func assignFlag(str string) string {
	if strings.Contains(str, "ERROR") {
		return flagError
	} else if strings.Contains(str, "CRITICAL") {
		return flagCritical
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
