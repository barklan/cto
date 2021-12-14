package checking

// TODO whole package with subpackages is dead code.

import (
	"fmt"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func LaunchChecks(b *tb.Bot, data *storage.Data, projectName string) {
	var wg sync.WaitGroup
	wg.Add(2)

	for _, url := range data.Config.P[projectName].Checks.SimpleURLChecks {
		time.Sleep(2 * time.Second) // don't want them to be started at the same time
		title := fmt.Sprintf("Healthcheck: %s", url)
		interval := 1 * time.Minute
		GoCheck(
			b,
			data,
			&wg,
			title,
			interval,
			CheckByExternalRequest,
			projectName,
			url,
			0,
		)
	}

	if data.Config.P[projectName].Checks.GitLab.FailedPipelinesMain {
		title := "Main pipelines healthcheck"
		interval := 5 * time.Minute
		branch := "main"
		GoCheck(b, data, &wg, title, interval,
			CheckFailedPipelines, projectName, branch)
	}

	wg.Wait()
	msg := "All registered checks exited."
	data.CSend(msg)
}

// // main
// wg.Add(1)
// go func() {
// 	defer func() {
// 		CrashExit(data, "All checks exited.")
// 		wg.Done()
// 	}()
// 	for projectName := range data.Config.P {
// 		time.Sleep(2 * time.Second) // we want some interval between outgoing requests
// 		checking.LaunchChecks(b, data, projectName)
// 	}
// }()
