package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/barklan/cto/pkg/logserver"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

// TODO fix this for v3
func registerStatusHandler(b *tb.Bot, data *storage.Data) {
	b.Handle("/status", func(m *tb.Message) {
		projectName, ok := VerifySender(data, m)
		if !ok {
			return
		}

		knownErrors := make([]logserver.KnownError, 0)
		knownErrorsRaw := data.Get(fmt.Sprintf("knownErrors-%s", projectName))
		if string(knownErrorsRaw) != "" {
			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
				data.CSend(fmt.Sprintf("Failed to unmarshal known errors for project: %s", projectName))
				return
			}
		}

		var msg string
		if len(knownErrors) == 0 {
			msg += "No known issues."
		} else {
			msg += fmt.Sprintf(
				"Recent issues (similarity threshold is set to %.2f):\n",
				data.Config.Internal.Log.SimilarityThreshold,
			)
			for index, knownError := range knownErrors {
				queryString := url.QueryEscape(knownError.OriginBadgerKey)
				exactLogURL := fmt.Sprintf(
					"%s/log/exact?key=%s",
					data.Config.Internal.Log.ServiceHostname,
					queryString,
				)
				badgerKeyArr := strings.Split(knownError.OriginBadgerKey, " ")
				prettyOrigin := strings.Join(badgerKeyArr[1:5], " ") // Omit project and everything after service name
				upperFlag := strings.ToUpper(badgerKeyArr[5])

				msg += fmt.Sprintf(
					" *%d. %s origin:* %s ([link to log](%s))\n    *last seen:* %s ago (%s)\n    *total count:* %d\n",
					index+1,
					upperFlag,
					prettyOrigin,
					exactLogURL,
					time.Since(knownError.LastSeen).Round(time.Second),
					knownError.LastSeen.Format("3:04PM MST"),
					knownError.Counter,
				)
			}
		}

		msg += "\n"

		periodicReport := logserver.PeriodicReport{}
		periodicReportRaw := data.Get(fmt.Sprintf("periodicLogReport-%s", projectName))
		if string(periodicReportRaw) == "" {
			msg += "Periodic report is not ready yet."
		} else {
			if err := json.Unmarshal(periodicReportRaw, &periodicReport); err != nil {
				log.Println("Falied to unmarshal periodic report", err)
			}

			msg += fmt.Sprintf(
				"Last %s: recieved %d events.",
				periodicReport.Period,
				periodicReport.Recieved,
			)
		}

		msg += "\n"

		msg += fmt.Sprintf("Logs are retained for %d hours. ", data.Config.Internal.Log.RetentionHours)

		// TODO delete this later
		msg += fmt.Sprintf("%s.", projectName)

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

		data.PSend(projectName, msg, tb.ModeMarkdown, selector)
		b.Delete(m)
	})
}
