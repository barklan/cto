package logserver

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/adrg/strutil"
	"github.com/adrg/strutil/metrics"
	"github.com/barklan/cto/pkg/logserver/types"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func notifyAboutError(
	data *storage.Data,
	projectName string,
	logData *LogMetadata,
	badgerKey string,
) {
	queryString := url.QueryEscape(badgerKey)
	exactLogURL := fmt.Sprintf(
		"%s/log/exact?key=%s",
		data.Config.Internal.Log.ServiceHostname,
		queryString,
	)
	selector := &tb.ReplyMarkup{}
	btnURL := selector.URL("View Log", exactLogURL)
	selector.Inline(
		selector.Row(btnURL),
	)

	// Get elapsed tiem
	timeStr := logData.Timestamp
	var extraTimeStr string
	fluentdTime, err := time.Parse("2006-01-02 15:04:05 -0700", timeStr)
	if err == nil {
		now := time.Now()
		elapsed := now.Sub(fluentdTime).Round(time.Second)
		extraTimeStr = fmt.Sprintf(" (%s ago)", elapsed)
	}

	upperFlag := strings.ToUpper(logData.Flag)
	message := fmt.Sprintf(
		"*%s* in %s (%s) at %s%s.",
		upperFlag,
		logData.Service,
		logData.Hostname,
		timeStr,
		extraTimeStr,
	)

	data.PSend(projectName, message, tb.ModeMarkdown, selector)
}

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

	log.Println("max similarity with previous error:", maxSimilarity)

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
	log.Println("Added new Error!")
	sessData.KnownErrorsMutex.Unlock()

	notifyAboutError(data, projectName, logData, badgerKey)
}
