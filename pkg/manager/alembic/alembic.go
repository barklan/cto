package alembic

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/barklan/cto/pkg/manager/docker"
	exec "github.com/barklan/cto/pkg/manager/exec"
)

func checkStringForMultipleHeads(input string) {
	r, err := regexp.Compile(`head`)
	if err != nil {
		log.Panic(err)
	}
	heads := r.FindAllString(input, -1)
	if len(heads) > 1 {
		log.Fatal("Multiple heads detected! " +
			"Deploy would not proceed. " +
			"Run 'alembic merge heads' to create a merge migration.")
	}
}

func GetAlembicHeadFromImage(image, imageTag string) string {
	output := docker.Run(image, imageTag, []string{"alembic", "heads"})

	checkStringForMultipleHeads(output)

	r, err := regexp.Compile(`^(\w+)`)
	if err != nil {
		log.Panic(err)
	}
	head := r.FindString(output)
	fmt.Println(head)

	return head
}

func GetAlembicHistoryFromImage(image, imageTag string) []string {
	output := docker.Run(image, imageTag, []string{"alembic", "history"})

	r, err := regexp.Compile(`Revision ID: (\w+)\b`)
	if err != nil {
		log.Panic(err)
	}
	historyRaw := r.FindAllStringSubmatch(output, -1)
	var history []string
	for _, val := range historyRaw {
		history = append(history, val[1])
	}
	return history
}

func GetAlembicVersionFromDB(stackName string) string {
	command := fmt.Sprintf(`docker exec $(docker ps -q -f name=%s_db) psql -U postgres -d app -t -c "SELECT version_num FROM alembic_version;"`,
		stackName)
	version, err := exec.ExecuteCmd(command)
	if err != nil {
		log.Panic(err)
	}
	version = strings.Trim(version, "\n ")
	return version
}
