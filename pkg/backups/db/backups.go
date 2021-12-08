package db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/sshclient"
	"github.com/barklan/cto/pkg/storage"
	tb "gopkg.in/tucnak/telebot.v2"
)

func writeBufferToFile(buffer *bytes.Buffer, fullFilename string) {
	file, err := os.Create(fullFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// FIXME
}

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

func SendBackupToTelegram(data *storage.Data, projectName, filename string) {
	file := &tb.Document{
		File:     tb.FromDisk(filename),
		Caption:  fmt.Sprintf("Compatible backend rev: %s", "TODO"), // FIXME
		FileName: filename,
	}

	data.PSend(projectName, file)
}

func BackupDatabase(data *storage.Data, projectName string) {
	driver := data.Config.P[projectName].Backups.DB.Driver

	if !data.Config.P[projectName].Backups.DB.Enable {
		data.CSend(fmt.Sprintf("Skipping db backups for project %s", projectName))
	}

	if driver != "postgres" {
		data.PSend(projectName, "Only postgres database backups supported.")
	}

	containerName := data.Config.P[projectName].Backups.DB.ContainerName
	databaseName := data.Config.P[projectName].Backups.DB.Database
	var cmd string

	switch driver {
	case "postgres":
		cmd = backupPostgresCmdString(containerName, databaseName)
	default:
		return
	}

	buffer, err := sshclient.ConnectAndExecute(data, projectName, cmd)
	if err != nil {
		data.CSend("Error when making backup.")
		data.CSend(err.Error())
		return
	}

	baseFileName := constructBackupFilename(containerName, databaseName)
	fullFilename := data.MediaPath + "/" + baseFileName
	err = ioutil.WriteFile(fullFilename, buffer.Bytes(), 0777)
	if err != nil {
		data.CSend(fmt.Sprintf("Failed to save backup for project %s", projectName))
		log.Println(err)
	}

	SendBackupToTelegram(data, projectName, fullFilename)
}

func PeriodicDBBackups(data *storage.Data, projectName string) {
	defer data.CSend("PeriodicDBBackups exited for project %s", projectName)

	interval := data.Config.P[projectName].Backups.DB.IntervalMinutes
	if interval > 0 {
		ticker := time.NewTicker(time.Duration(interval) * time.Minute)

		for range ticker.C {
			BackupDatabase(data, projectName)
		}
	}
}

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