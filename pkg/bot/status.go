package bot

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	tb "gopkg.in/tucnak/telebot.v2"
)

// FIXME recovery for v5
func (s *Sylon) registerStatusHandler() {
	s.B.Handle("/status", func(m *tb.Message) {
		project, _, ok := s.VerifySender(m)
		if !ok {
			return
		}
		var msg string

		// 		knownErrors := make([]types.KnownError, 0)
		// 		knownErrorsRaw := data.Get(fmt.Sprintf("knownErrors-%s", projectName))
		// 		if string(knownErrorsRaw) != "" {
		// 			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
		// 				data.CSend(fmt.Sprintf("Failed to unmarshal known errors for project: %s", projectName))
		// 				return
		// 			}
		// 		}

		// 		if len(knownErrors) == 0 {
		// 			msg += `No known issues\.`
		// 		} else {
		// 			recentIssuesHeader := fmt.Sprintf(
		// 				`Recent issues \(threshold is %.2f\):
		// `,
		// 				data.Config.Internal.Log.SimilarityThreshold,
		// 			)
		// 			recentIssuesHeader = strings.Replace(recentIssuesHeader, ".", `\.`, -1)
		// 			msg += recentIssuesHeader
		// 			for index, knownError := range knownErrors {
		// 				queryString := url.QueryEscape(knownError.OriginBadgerKey)
		// 				exactLogURL := fmt.Sprintf(
		// 					"%s/log/exact?key=%s",
		// 					data.Config.Internal.Log.ServiceHostname,
		// 					queryString,
		// 				)
		// 				badgerKeyArr := strings.Split(knownError.OriginBadgerKey, " ")
		// 				prettyOrigin := strings.Join(badgerKeyArr[1:5], " ")
		// 				prettyOrigin = strings.Replace(prettyOrigin, ".", `\.`, -1)
		// 				prettyOrigin = strings.Replace(prettyOrigin, "-", `\-`, -1)
		// 				upperFlag := strings.ToUpper(badgerKeyArr[5])

		// 				msg += fmt.Sprintf(
		// 					`*%d\. %s:* [%s](%s)
		// \- *total:* %d \| *last:* %s ago
		// `,
		// 					index+1,
		// 					upperFlag,
		// 					prettyOrigin,
		// 					exactLogURL,
		// 					// knownError.LastSeen.Format("3:04PM MST"),
		// 					knownError.Counter,
		// 					time.Since(knownError.LastSeen).Round(time.Second),
		// 				)
		// 			}
		// 		}

		// 		msg += `
		// `

		// periodicReport := logservertypes.PeriodicReport{}
		// periodicReportRaw := data.Get(fmt.Sprintf("periodicLogReport-%s", projectName))
		// if string(periodicReportRaw) == "" {
		// 	msg += `Periodic report is not ready yet\. `
		// } else {
		// 	if err := json.Unmarshal(periodicReportRaw, &periodicReport); err != nil {
		// 		log.Println("Falied to unmarshal periodic report", err)
		// 	}

		// 	msg += fmt.Sprintf(
		// 		`Last %s: recieved %d events\. `,
		// 		periodicReport.Period,
		// 		periodicReport.Recieved,
		// 	)
		// }

		msg += fmt.Sprintf(`Logs are retained for %d hours\. `, s.Config.Log.RetentionHours)

		// TODO badger keys like this should not be magic strings
		authToken, ok, err := s.Cache.Get(fmt.Sprintf("authToken-%s", project.ID))
		if err != nil {
			log.Panicln("Failed to access cache.")
		}
		if !ok {
			s.JustSend(m.Chat, "Auth token not found in cache.")
		}

		panelURL := fmt.Sprintf(
			"%s/log?token=%s",
			s.Config.Log.ServiceHostname,
			authToken,
		)
		selector := &tb.ReplyMarkup{}
		btnURL := selector.URL("Search the Logs", panelURL)
		selector.Inline(
			selector.Row(btnURL),
		)

		println(msg)

		s.JustSend(m.Chat, msg, tb.ModeMarkdownV2, selector)
		_ = s.B.Delete(m)
	})
}
