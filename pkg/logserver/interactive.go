package logserver

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/barklan/cto/pkg/logserver/types"
	"github.com/barklan/cto/pkg/storage"
)

func similarErrorExists(
	data *storage.Data,
	logData *LogMetadata,
	sessData *SessionData,
	recordStr string,
) bool {
	var maxSimilarity float64
	var maxSimilarityIndex int
	for i, knownErr := range sessData.KnownErrors {
		if knownErr.Hostname != logData.Hostname || knownErr.Service != logData.Service {
			continue
		}

		jacMetic := metrics.NewJaccard()
		similarity := strutil.Similarity(knownErr.LogStr, recordStr, jacMetic)
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			maxSimilarityIndex = i
		}
	}

	data.Log.Info("max similarity with previous error", zap.Float64("maxSimilarity", maxSimilarity))

	if maxSimilarity > data.Config.Internal.Log.SimilarityThreshold {
		sessData.KnownErrors[maxSimilarityIndex].Counter++
		sessData.KnownErrors[maxSimilarityIndex].LastSeen = time.Now()
		return true
	}

	return false
}

func handleErrorRecord(
	data *storage.Data,
	logData *LogMetadata,
	record RawLogRecord,
	badgerKey string,
	sessData *SessionData,
	projectName string,
) {
	sessData.KnownErrorsMutex.Lock()

	recordStr := fmt.Sprint(record)

	if similarErrorExists(data, logData, sessData, recordStr) {
		sessData.KnownErrorsMutex.Unlock()
		return
	}

	newError := types.KnownError{
		Hostname:        logData.Hostname,
		Service:         logData.Service,
		OriginBadgerKey: badgerKey,
		LogStr:          recordStr,
		Counter:         1,
		LastSeen:        time.Now(),
	}
	sessData.KnownErrors = append(sessData.KnownErrors, newError)
	data.Log.Info("added new issue", zap.String("project", projectName))
	sessData.KnownErrorsMutex.Unlock()

	raw := data.GetLog(badgerKey)
	if err := data.Cache.Set(badgerKey, raw, 336*time.Hour); err != nil {
		warn := "failed to set new error in cache"
		data.Log.Warn(warn, zap.String("badgerKey", badgerKey))
		data.InternalAlert(warn)
	}

	data.NewIssue(
		projectName,
		logData.Hostname,
		logData.Service,
		logData.Timestamp,
		badgerKey,
		logData.Flag,
	)
}
