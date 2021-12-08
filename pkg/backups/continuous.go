package backups

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/storage"
)

func ContinuousDBBackupsAllProjects(data *storage.Data) {
	defer data.CSend("ContinuousDBBackupsAllProjects exited for all projects.")
	wg := new(sync.WaitGroup)
	wg.Add(len(data.Config.P))

	for projectName := range data.Config.P {
		go func(pName string) {
			defer wg.Done()
			ContinuousDBBackups(data, pName)
		}(projectName)
	}

	wg.Wait()
}

func ContinuousDBBackups(data *storage.Data, projectName string) {
	defer data.CSend(fmt.Sprintf("ContinuousDBBackups exited for project %s", projectName))

	if !data.Config.P[projectName].Backups.DB.Continuous {
		data.CSend(fmt.Sprintf("Skipping ContinuousDBBackups for project %s", projectName))
		return
	}

	if data.Config.P[projectName].Backups.DB.Driver != "postgres" {
		data.PSend(projectName, "ContinuousDBBackups only for postgres.")
		return
	}

	ticker := time.NewTicker(30 * time.Second)

	// FIXME base backup should be performed on handle (check if continuity is enabled and wal files exist)
	// FIXME there should be handle (or periodic goroutine that cleans both basebackup and all wal files)
	for range ticker.C {
		err := PerformContinuity(data, projectName)
		if err != nil {
			data.PSend(projectName, "Contious backup failed. Will retry later.")
			data.PSend(projectName, err.Error())
			log.Println(err)
		}
	}
}
