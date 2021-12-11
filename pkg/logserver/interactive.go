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

func handleErrorRecordInteractive(
	data *storage.Data,
	logData map[string]string,
	record RawLogRecord,
	badgerKey string,
	sessData *SessionData,
	projectName string,
) error {
	sessData.KnownErrorsMutex.Lock()

	recordStr := fmt.Sprint(record)

	var maxSimilarity float64
	var maxSimilarityIndex int
	for i, knownError := range sessData.KnownErrors {
		similarity := strutil.Similarity(knownError.LogStr, recordStr, metrics.NewHamming())
		if similarity > maxSimilarity {
			maxSimilarity = similarity
			maxSimilarityIndex = i
		}
	}
	if maxSimilarity > data.Config.Internal.Log.SimilarityThreshold {
		sessData.KnownErrors[maxSimilarityIndex].Counter++
		sessData.KnownErrors[maxSimilarityIndex].LastSeen = time.Now()
		sessData.KnownErrorsMutex.Unlock()
		return nil
	}

	// similar errors not found at this point

	// add error
	newError := types.KnownError{
		OriginBadgerKey: badgerKey,
		LogStr:          recordStr,
		Counter:         1,
		LastSeen:        time.Now(),
	}
	sessData.KnownErrors = append(sessData.KnownErrors, newError)
	log.Println("Added new Error!")
	sessData.KnownErrorsMutex.Unlock()

	// notify
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
	timeStr := logData["fluentd_time"]
	var extraTimeStr string
	fluentdTime, err := time.Parse("2006-01-02 15:04:05 -0700", timeStr)
	if err == nil {
		now := time.Now()
		elapsed := now.Sub(fluentdTime).Round(time.Second)
		extraTimeStr = fmt.Sprintf(" (%s ago)", elapsed)
	}

	upperFlag := strings.ToUpper(logData["flag"])
	message := fmt.Sprintf(
		"*%s* in %s (%s) at %s%s.",
		upperFlag,
		logData["service_name"],
		logData["fluentd_hostname"],
		timeStr,
		extraTimeStr,
	)

	data.PSend(projectName, message, tb.ModeMarkdown, selector)

	return nil
}
