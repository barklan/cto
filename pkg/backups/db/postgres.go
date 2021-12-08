package db

import (
	"fmt"
)

func backupPostgresCmdString(containerName, databaseName string) string {
	cmd := fmt.Sprintf(
		"docker exec $(docker ps -q -f name=%s) pg_dump -U postgres %s",
		containerName,
		databaseName,
	)
	return cmd
}
