package backups

import (
	"fmt"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func constructBackupFilename(targetContainer, targetDatabase string) string {
	backupFilenameStamp := time.Now().Format("2006_01_02__15_04_05")
	backupFilename := fmt.Sprintf(
		"%s_%s_%s.sql.gz",
		targetContainer,
		targetDatabase,
		backupFilenameStamp,
	)
	return backupFilename
}

func sendBackupToTelegram(data *storage.Data, projectName, filename string) {
	file := &tb.Document{
		File:     tb.FromDisk(filename),
		Caption:  fmt.Sprintf("Compatible backend rev: %s", "TODO"), // FIXME
		FileName: filename,
	}

	data.PSend(projectName, file)
}

func BackupDatabase(data *storage.Data, projectName string) {
	driver := data.Config.P[projectName].Backups.DB.Driver

	if driver != "postgres" {
		data.PSend(projectName, "Only postgres database backups supported.")
	}

	containerName := data.Config.P[projectName].Backups.DB.ContainerName
	databaseName := data.Config.P[projectName].Backups.DB.Database

	switch driver {
	case "postgres":
		performDump(data, projectName, containerName, databaseName)
	default:
		return
	}
}

func PeriodicDBBackups(data *storage.Data, projectName string) {
	defer data.CSend(fmt.Sprintf("PeriodicDBBackups exited for project %s", projectName))

	if !data.Config.P[projectName].Backups.DB.Enable {
		data.CSend(fmt.Sprintf("Skipping db backups for project %s", projectName))
		return
	}

	interval := data.Config.P[projectName].Backups.DB.IntervalMinutes
	if interval <= 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Minute)

	for range ticker.C {
		BackupDatabase(data, projectName)
	}
}

// TODO should use scheduling
func PeriodicDBBackupsAllProjects(data *storage.Data) {
	defer data.CSend("PeriodicDBBackupsAllProjects exited for all projects.")
	wg := new(sync.WaitGroup)
	wg.Add(len(data.Config.P))

	for projectName := range data.Config.P {
		go func(pName string) {
			defer wg.Done()
			PeriodicDBBackups(data, pName)
		}(projectName)
	}

	wg.Wait()
}
