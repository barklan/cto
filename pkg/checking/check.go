package checking

import (
	"fmt"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func LaunchChecks(b *tb.Bot, data *storage.Data, projectName string) {
	var wg sync.WaitGroup

	for _, url := range data.Config.P[projectName].Checks.SimpleURLChecks {
		time.Sleep(2 * time.Second) // don't want them to be started at the same time
		title := fmt.Sprintf("Healthcheck: %s", url)
		interval := 1 * time.Minute
		GoCheck(b, data, &wg, title, interval,
			CheckByExternalRequest, projectName, url, 0)
	}

	// if data.Config.P[projectName].Checks.GitLab.FailedPipelinesMain == true {
	// 	title := "Main pipelines healthcheck"
	// 	interval := 5 * time.Minute
	// 	branch := "main"
	// 	GoCheck(b, data, &wg, title, interval,
	// 		CheckFailedPipelines, branch)
	// }

	// if data.Config.P[projectName].Checks.GitLab.MRApprovals == true {
	// 	title := "Check approved mrs (slow)"
	// 	interval := 8 * time.Hour
	// 	GoCheck(b, data, &wg, title, interval, CheckEachMrInApprovedMREventsSlow)

	// 	data.SetObj("fastMRChecksMuted", "false", -1)
	// 	interval = 45 * time.Minute
	// 	title = "Check approved mrs (fast)"
	// 	GoCheck(b, data, &wg, title, interval, CheckEachMrInApprovedMREventsFast)
	// }

	// TODO do we really need this handle? Aren't "less is more"?
	// b.Handle("/checks", func(m *tb.Message) {
	// 	var str string
	// 	for _, check := range checklist {
	// 		str = str + fmt.Sprintf("\n - %s", check)
	// 	}
	// 	msg := fmt.Sprintf("Registered checks: %v.", str)
	// 	data.CSend(msg)
	// })

	wg.Wait()
	msg := "All registered checks exited."
	data.CSend(msg)
}
