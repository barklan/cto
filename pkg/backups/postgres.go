package backups

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/barklan/cto/pkg/sshclient"
	"github.com/barklan/cto/pkg/storage"
)

func backupPostgresCmdString(containerName, databaseName string) string {
	cmd := fmt.Sprintf(
		"docker exec $(docker ps -q -f name=%s) pg_dump -U postgres %s | gzip",
		containerName,
		databaseName,
	)
	return cmd
}

func performDump(data *storage.Data, projectName, containerName, databaseName string) {
	cmd := backupPostgresCmdString(containerName, databaseName)

	sshData := sshclient.SSHConnectionData{
		Hostname: data.Config.P[projectName].Backups.DB.SSHHostname,
		Username: data.Config.P[projectName].Backups.DB.SSHHostname,
		Keypath:  data.Config.P[projectName].Backups.DB.SSHHostname,
	}

	buffer, err := sshclient.ConnectAndExecute(data, projectName, sshData, cmd)
	if err != nil {
		data.PSend(projectName, "Error when making backup.")
		data.PSend(projectName, err.Error())
		return
	}

	baseFileName := constructBackupFilename(containerName, databaseName)
	data.CreateMediaDirIfNotExists(projectName + "/sqldumps")
	fullFilename := data.MediaPath + "/" + projectName + "/sqldumps/" + baseFileName
	err = ioutil.WriteFile(fullFilename, buffer.Bytes(), 0777)
	if err != nil {
		data.PSend(projectName, fmt.Sprintf("Failed to save backup for project %s", projectName))
		log.Println(err)
	}

	sendBackupToTelegram(data, projectName, fullFilename)
}

func baseBackupPostgresCmdString(containerName string) string {
	cmd := fmt.Sprintf(
		"docker exec $(docker ps -q -f name=%s) pg_basebackup -U postgres -D /pgbackups -Ft -z",
		containerName,
	)
	return cmd
}

func PerformBaseBackup(data *storage.Data, projectName string) {
	// FIXME
}

func PerformContinuity(data *storage.Data, projectName string) error {
	// Remote command lists only files (without directories).
	sshData := sshclient.SSHConnectionData{
		Hostname: data.Config.P[projectName].Backups.DB.SSHHostname,
		Username: data.Config.P[projectName].Backups.DB.SSHHostname,
		Keypath:  data.Config.P[projectName].Backups.DB.SSHHostname,
	}

	remoteFolder := data.Config.P[projectName].Backups.DB.ContinuousPath
	output, err := sshclient.ConnectAndExecute(
		data,
		projectName,
		sshData,
		fmt.Sprintf("ls -p %s | grep -v /", remoteFolder),
	)
	if err != nil {
		return err
	}

	remoteFiles := strings.Fields(output.String())
	log.Printf("Remote files: %s", remoteFiles)

	localFiles := make([]string, 0)

	data.CreateMediaDirIfNotExists(projectName + "/pgbackups")
	localFolder := data.MediaPath + "/" + projectName + "/pgbackups"
	files, err := ioutil.ReadDir(localFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			localFiles = append(localFiles, file.Name())
		}
	}

	log.Printf("Local files: %s", localFiles)

	filesToCopy := difference(remoteFiles, localFiles)
	log.Printf("Files to copy: %s", filesToCopy)

	if len(filesToCopy) == 0 {
		return nil
	}

	containerName := data.Config.P[projectName].Backups.DB.ContainerName

	for _, file := range filesToCopy {
		remoteFilePath := remoteFolder + "/" + file
		localFilePath := localFolder + "/" + file

		log.Printf("Transferring %s", file)

		err = sshclient.SFTP(data, projectName, sshData, remoteFilePath, localFilePath)
		if err != nil {
			_, e := os.Stat(localFilePath)
			if errors.Is(e, os.ErrNotExist) {
				return err
			}

			e = os.Remove(localFilePath)
			if e != nil {
				data.PSend(projectName, "Continous backup is corrupted! Please intervene. Panic in 10 seconds.")
				time.Sleep(10 * time.Second)
				log.Panic(e)
			}

			return err
		}

		log.Printf("Removing remote  %s", file)

		_, err = sshclient.ConnectAndExecute(
			data,
			projectName,
			sshData,
			fmt.Sprintf("rm %s", remoteFilePath),
		)
		if err != nil {
			data.PSend(
				projectName,
				fmt.Sprintf("%s. Failed to remove copied WAL files from server.", containerName),
			)
		}
	}

	data.PSend(
		projectName,
		fmt.Sprintf(
			"%s. Successfully backed up %d new WAL files.",
			containerName,
			len(filesToCopy),
		),
	)

	return nil
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
