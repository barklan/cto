package logserver

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/barklan/cto/pkg/core/logserver/types"
	"github.com/barklan/cto/pkg/core/storage"
	"github.com/gofrs/uuid"
)

func similarErrorExists(
	data *storage.Data,
	logData *LogMetadata,
	sessData *SessionData,
	recordStr string,
) bool {
	var maxSimilarity float64
	var maxSimilarityIndex string
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
		if kErr, ok := sessData.KnownErrors[maxSimilarityIndex]; ok {
			kErr.Counter++
			kErr.LastSeen = time.Now()
		} else {
			data.Log.Warn("expected to find knownError by key but none found", zap.String("key", maxSimilarityIndex))
		}

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
	sessData.Mutex.Lock()

	recordStr := fmt.Sprint(record)

	if similarErrorExists(data, logData, sessData, recordStr) {
		sessData.Mutex.Unlock()
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
	uid4, err := uuid.NewV4()
	if err != nil {
		data.Log.Error("failed to generate uuid for new error", zap.Error(err))
		return
	}
	u4 := uid4.String()
	sessData.KnownErrors[u4] = newError
	data.Log.Info("added new issue", zap.String("project", projectName))
	sessData.Mutex.Unlock()

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
