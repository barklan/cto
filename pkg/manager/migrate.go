package manager

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/barklan/cto/pkg/manager/alembic"
)

func Migrate(stackName string) {
	pgIsReady := make(chan string, 1)
	go func() {
		for i := 0; i < 40; i++ {
			out, err := exec.Command("docker", "exec", fmt.Sprintf("%s_db", stackName), "pg_isready").Output()
			if err != nil {
				pgIsReady <- string(out)
				break
			}
			time.Sleep(3 * time.Second)
		}
	}()

	select {
	case res := <-pgIsReady:
		fmt.Println(res)
		alembic.Migrate()
	case <-time.After(3 * time.Minute):
		log.Fatal("Timeout after 3 minutes waiting for postgres to be ready. " +
			"Migrations will not be applied.")
	}
}
