package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	logservertypes "github.com/barklan/cto/pkg/logserver/types"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func getSLAinfo(data *storage.Data, projectName string) string {
	totalRuningTimeKey := fmt.Sprintf("%s-totalRunningTime", projectName)
	var totalRunningTime time.Duration
	totalRunningTimeRaw := data.Get(totalRuningTimeKey)
	if string(totalRunningTimeRaw) == "" {
		totalRunningTime = 0
	} else {
		err := json.Unmarshal(totalRunningTimeRaw, &totalRunningTime)
		if err != nil {
			data.CSend("Failed to unmarshal totalRunningTime")
		}
	}

	downTimeKey := fmt.Sprintf("%s-downTime", projectName)
	var totalDownTime time.Duration
	totalDownTimeRaw := data.Get(downTimeKey)
	if string(totalDownTimeRaw) == "" {
		totalDownTime = 0
	} else {
		err := json.Unmarshal(totalDownTimeRaw, &totalDownTime)
		if err != nil {
			data.CSend("Failed to unmarshal totalDownTime")
		}
	}

	invertedSLA := float64(totalDownTime) / float64(totalRunningTime) * 100
	sla := 100.0 - invertedSLA
	return fmt.Sprintf(
		`Uptime: %.10f%% \(total: %s, down: %s\)\. `,
		sla,
		totalRunningTime.Round(time.Second),
		totalDownTime.Round(time.Second),
	)
}

func registerStatusHandler(b *tb.Bot, data *storage.Data) {
	b.Handle("/status", func(m *tb.Message) {
		projectName, ok := VerifySender(data, m)
		if !ok {
			return
		}

		knownErrors := make([]logservertypes.KnownError, 0)
		knownErrorsRaw := data.Get(fmt.Sprintf("knownErrors-%s", projectName))
		if string(knownErrorsRaw) != "" {
			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
				data.CSend(fmt.Sprintf("Failed to unmarshal known errors for project: %s", projectName))
				return
			}
		}

		var msg string
		if len(knownErrors) == 0 {
			msg += `No known issues\.`
		} else {
			recentIssuesHeader := fmt.Sprintf(
				`Recent issues \(threshold is %.2f\):
`,
				data.Config.Internal.Log.SimilarityThreshold,
			)
			recentIssuesHeader = strings.Replace(recentIssuesHeader, ".", `\.`, -1)
			msg += recentIssuesHeader
			for index, knownError := range knownErrors {
				queryString := url.QueryEscape(knownError.OriginBadgerKey)
				exactLogURL := fmt.Sprintf(
					"%s/log/exact?key=%s",
					data.Config.Internal.Log.ServiceHostname,
					queryString,
				)
				badgerKeyArr := strings.Split(knownError.OriginBadgerKey, " ")
				prettyOrigin := strings.Join(badgerKeyArr[1:5], " ")
				prettyOrigin = strings.Replace(prettyOrigin, ".", `\.`, -1)
				prettyOrigin = strings.Replace(prettyOrigin, "-", `\-`, -1)
				upperFlag := strings.ToUpper(badgerKeyArr[5])

				msg += fmt.Sprintf(
					`*%d\. %s:* [%s](%s)
\- *last:* %s ago \| *total:* %d
`,
					index+1,
					upperFlag,
					prettyOrigin,
					exactLogURL,
					time.Since(knownError.LastSeen).Round(time.Second),
					// knownError.LastSeen.Format("3:04PM MST"),
					knownError.Counter,
				)
			}
		}

		msg += `
`

		periodicReport := logservertypes.PeriodicReport{}
		periodicReportRaw := data.Get(fmt.Sprintf("periodicLogReport-%s", projectName))
		if string(periodicReportRaw) == "" {
			msg += `Periodic report is not ready yet\. `
		} else {
			if err := json.Unmarshal(periodicReportRaw, &periodicReport); err != nil {
				log.Println("Falied to unmarshal periodic report", err)
			}

			msg += fmt.Sprintf(
				`Last %s: recieved %d events\.`,
				periodicReport.Period,
				periodicReport.Recieved,
			)
		}

		msg += fmt.Sprintf(`Logs are retained for %d hours\. `, data.Config.Internal.Log.RetentionHours)

		msg += getSLAinfo(data, projectName)

		msg += fmt.Sprintf(`%s\.`, projectName)

		// TODO badger keys like this should not be magic strings
		authToken := data.GetStr(fmt.Sprintf("authToken-%s", projectName))
		if authToken == "" {
			data.CSend(fmt.Sprintf("There is no auth token in badger for project %s.", projectName))
		}
		panelURL := fmt.Sprintf(
			"%s/log?token=%s",
			data.Config.Internal.Log.ServiceHostname,
			authToken,
		)
		selector := &tb.ReplyMarkup{}
		btnURL := selector.URL("Search the Logs", panelURL)
		selector.Inline(
			selector.Row(btnURL),
		)

		println(msg)

		data.PSend(projectName, msg, tb.ModeMarkdownV2, selector)
		_ = b.Delete(m)
	})
}
