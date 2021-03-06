package bot

import (
	"fmt"

	"github.com/barklan/cto/pkg/security"
	"github.com/barklan/cto/pkg/vars"
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

		authToken, err := security.CreateJWT(s.Config, vars.Guest, project.ID)
		if err != nil {
			s.JustSend(m.Chat, "Failed to create JWT.")
			return
		}

		panelURL := fmt.Sprintf(
			"%s/guest?token=%s&name=%s&project=%s",
			s.Config.Log.ServiceHostname,
			authToken,
			"guest",
			project.ID,
		)
		selector := &tb.ReplyMarkup{}
		btnURL := selector.URL("Panel", panelURL)
		selector.Inline(
			selector.Row(btnURL),
		)

		println(msg)

		s.JustSend(m.Chat, msg, tb.ModeMarkdownV2, selector)
		_ = s.B.Delete(m)
	})
}
