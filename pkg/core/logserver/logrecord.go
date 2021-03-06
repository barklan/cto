package logserver

import (
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/barklan/cto/pkg/core/namespaces"
	"github.com/barklan/cto/pkg/core/storage"
	"go.uber.org/zap"
)

const (
	flagError     = "err"
	flagCritical  = "crit"
	flagFatal     = "fatal"
	flagEmergency = "emerg"
	flagNone      = "none"
)

type RawLogRecord map[string]interface{}

type LogMetadata struct {
	Hostname  string `json:"fluentd_hostname"`
	Service   string `json:"service_name"`
	Timestamp string `json:"fluentd_time"`
	Flag      string
}

func serviceFromContainer(container string) string {
	service := strings.SplitN(container, ".", 2)[0]
	return service[1:]
}

func constructMetadata(record RawLogRecord) (*LogMetadata, error) {
	logMetadata := LogMetadata{}

	key := "fluentd_hostname"
	if val, ok := record[key]; ok {
		if v, okType := val.(string); okType {
			logMetadata.Hostname = v
		}
	}
	if logMetadata.Hostname == "" {
		logMetadata.Hostname = "undefined"
	}

	key = "container_name"
	if val, ok := record[key]; ok {
		if v, okType := val.(string); okType {
			logMetadata.Service = serviceFromContainer(v)
		}
	}
	if logMetadata.Service == "" {
		logMetadata.Service = "undefined"
	}

	key = "fluentd_time"
	if val, ok := record[key]; ok {
		if v, okType := val.(string); okType {
			logMetadata.Timestamp = v
		}
	}
	if logMetadata.Timestamp == "" {
		// FIXME reasonable defaults here and ingest msg
	}

	return &logMetadata, nil
}

func processLogRecord(
	data *storage.Data,
	pid string,
	record RawLogRecord,
	sessData *SessionData,
	reportChan chan LogRecordReport,
) {
	logData, err := constructMetadata(record)
	if err != nil {
		data.Log.Error("failed to construct metadata", zap.String("project", pid), zap.Error(err))
		log.Printf("%s\n", record)
		return
	}

	// FIXME hardcoded, but still better to have some small number than nothing
	if rand.Intn(1) == 0 {
		sessData.Mutex.Lock()

		sinceLast := namespaces.GetLastRefresh(data, pid)
		// FIXME maybe need more than that
		if sinceLast > 6*time.Hour {
			namespaces.Clear(data, pid)
			namespaces.SetLastRefresh(data, pid)
		}

		knownEnvs := namespaces.GetKnownEnvs(data, pid)
		if _, ok := knownEnvs[logData.Hostname]; !ok {
			knownEnvs[logData.Hostname] = struct{}{}
			namespaces.SetKnownEnvs(data, pid, knownEnvs)
		}
		// This is for querying purposes.
		// TODO it is heavy to do it on every log record - should decide randomly instead.
		knownServices := namespaces.GetKnownServices(data, pid, logData.Hostname)
		if _, ok := knownServices[logData.Service]; !ok {
			knownServices[logData.Service] = struct{}{}
			namespaces.SetKnownServices(data, pid, logData.Hostname, knownServices)
		}
		sessData.Mutex.Unlock()
	}

	// flag := assignFlag(fmt.Sprint(record))
	// logData.Flag = flag
	// FIXME no error reporting for now
	logData.Flag = flagNone

	// Save log record
	badgerKey, _ := constructBadgerKey(logData, pid)
	retentionDuration := time.Duration(data.Config.Internal.Log.RetentionHours) * time.Hour

	data.SetLog(badgerKey, record, retentionDuration)

	if logData.Flag != flagNone {
		handleErrorRecord(data, logData, record, badgerKey, sessData, pid)
	}

	reportChan <- LogRecordReport{ProjectName: pid}
}

func constructBadgerKey(
	logData *LogMetadata,
	pid string,
) (string, error) {
	randID := strconv.FormatInt(rand.Int63(), 10)
	fluentdTimeArr := strings.Fields(logData.Timestamp)
	dateString, timeString := fluentdTimeArr[0], fluentdTimeArr[1]
	badgerKey := strings.Join(
		[]string{
			pid,
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
