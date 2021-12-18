package bot

// FIXME recovery for v5
// func (s *Sylon) registerStatusHandler() {
// 	s.B.Handle("/status", func(m *tb.Message) {
// 		project, chat, ok := s.VerifySender(m)
// 		if !ok {
// 			return
// 		}

// 		knownErrors := make([]types.KnownError, 0)
// 		knownErrorsRaw := data.Get(fmt.Sprintf("knownErrors-%s", projectName))
// 		if string(knownErrorsRaw) != "" {
// 			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
// 				data.CSend(fmt.Sprintf("Failed to unmarshal known errors for project: %s", projectName))
// 				return
// 			}
// 		}

// 		var msg string
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

// 		periodicReport := logservertypes.PeriodicReport{}
// 		periodicReportRaw := data.Get(fmt.Sprintf("periodicLogReport-%s", projectName))
// 		if string(periodicReportRaw) == "" {
// 			msg += `Periodic report is not ready yet\. `
// 		} else {
// 			if err := json.Unmarshal(periodicReportRaw, &periodicReport); err != nil {
// 				log.Println("Falied to unmarshal periodic report", err)
// 			}

// 			msg += fmt.Sprintf(
// 				`Last %s: recieved %d events\. `,
// 				periodicReport.Period,
// 				periodicReport.Recieved,
// 			)
// 		}

// 		msg += fmt.Sprintf(`Logs are retained for %d hours\. `, data.Config.Internal.Log.RetentionHours)

// 		msg += fmt.Sprintf(`%s\.`, projectName)

// 		// TODO badger keys like this should not be magic strings
// 		authToken := data.GetStr(fmt.Sprintf("authToken-%s", projectName))
// 		if authToken == "" {
// 			data.CSend(fmt.Sprintf("There is no auth token in badger for project %s.", projectName))
// 		}
// 		panelURL := fmt.Sprintf(
// 			"%s/log?token=%s",
// 			data.Config.Internal.Log.ServiceHostname,
// 			authToken,
// 		)
// 		selector := &tb.ReplyMarkup{}
// 		btnURL := selector.URL("Search the Logs", panelURL)
// 		selector.Inline(
// 			selector.Row(btnURL),
// 		)

// 		println(msg)

// 		data.PSend(projectName, msg, tb.ModeMarkdownV2, selector)
// 		_ = b.Delete(m)
// 	})
// }
