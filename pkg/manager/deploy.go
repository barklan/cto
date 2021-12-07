package manager

import (
	"fmt"
	"log"
	"os"

	"github.com/barklan/cto/pkg/manager/alembic"
	"github.com/barklan/cto/pkg/manager/exec"
	"github.com/joho/godotenv"
)

func resetAndDeployCmds(deployCmd string) []string {
	commands := []string{}

	pgTerminateConnection := `docker exec $(docker ps -q -f name=stag_db) psql -U postgres -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'app' AND pid <> pg_backend_pid();"`
	pgTerminateConnectionMulti := []string{}
	for i := 0; i < 10; i++ {
		pgTerminateConnectionMulti = append(pgTerminateConnectionMulti, pgTerminateConnection)
	}

	commands = append(commands,
		`docker exec $(docker ps -q -f name=stag_db) psql -U postgres -c "REVOKE CONNECT ON DATABASE app FROM public;"`,
		`docker exec $(docker ps -q -f name=stag_db) psql -U postgres -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = 'app' AND pid <> pg_backend_pid();"`,
	)
	commands = append(commands, pgTerminateConnectionMulti...)
	commands = append(commands,
		// Can be used with postgres 13
		// `docker exec $(docker ps -q -f name=stag_db) psql -U postgres -c "DROP DATABASE app WITH (FORCE);"`,
		`docker exec $(docker ps -q -f name=stag_db) psql -U postgres -c "DROP DATABASE IF EXISTS app;"`,
		`docker exec $(docker ps -q -f name=stag_db) bash -c "createdb -U postgres -T template0 app"`,
		`docker exec -i $(docker ps -q -f name=stag_db) psql -U postgres app < ../db_dump_stag.sql`,
		`docker exec $(docker ps -q -f name=stag_db) psql -U postgres -c "GRANT CONNECT ON DATABASE app TO public;"`,
		deployCmd,
	)

	return commands
}

func Deploy(target, backendImage string) {
	deployCmd := fmt.Sprintf(`cd /home/ubuntu/%s \
	&& docker login -u %s -p %s registry.gitlab.com/nftgalleryx/nftgallery_backend \
	&& docker stack deploy -c docker-stack.yml --with-registry-auth %s`,
		target,
		os.Getenv("GITLAB_TOKEN_USERNAME"),
		os.Getenv("GITLAB_TOKEN_PASSWORD"),
		target,
	)

	targetTag := os.Getenv("TAG")
	targetBranch := os.Getenv("BRANCH")

	targetAlembicHead := alembic.GetAlembicHeadFromImage(backendImage, targetTag)
	targetAlembicHistory := alembic.GetAlembicHistoryFromImage(backendImage, targetTag)
	currentAlembicVersion := alembic.GetAlembicVersionFromDB(target)

	currentAlembicVersionIsInHistory := false
	for _, val := range targetAlembicHistory {
		if currentAlembicVersion == val {
			currentAlembicVersionIsInHistory = true
			break
		}
	}

	var currentTag string
	err := godotenv.Load()
	if err != nil {
		currentTag = ""
		// currentBranch := ""
	} else {
		currentTag = os.Getenv("CURRENT_TAG")
		// currentBranch := os.Getenv("CURRENT_BRANCH")
	}

	commands := []string{deployCmd}
	// TODO make migrations part of ci instead of prestart script
	needMigrate := false

	if targetTag == currentTag {
		log.Println("Fast deploy: deploying the same image.")
	} else if target == "prod" {
		log.Println("Fast deploy: deploying on prod.")
	} else if currentAlembicVersion == targetAlembicHead {
		log.Println("Fast deploy: alembic head is the same.")
	} else if currentAlembicVersionIsInHistory {
		log.Println("Fast deploy: current alembic version exists in history.")
		needMigrate = true
	} else {
		log.Println("Destructive deploy. Some data may be lost.")
		commands = resetAndDeployCmds(deployCmd)
		needMigrate = true
	}

	exec.ExecuteCmds(commands)

	WriteCurrentVersion(targetTag, targetBranch)

	if needMigrate {
		log.Println("Launching migrations! (not really)")
		Migrate(target)
	}
}
