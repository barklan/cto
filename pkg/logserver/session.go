package logserver

import (
	"log"

	"github.com/barklan/cto/pkg/storage"
	"github.com/barklan/cto/pkg/storage/vars"
)

func openSession(data *storage.Data) map[string]*SessionData {
	sessionDataMap := map[string]*SessionData{}

	projects := make([]string, 0)
	if err := data.R.Select(&projects, "select id from project"); err != nil {
		log.Println("no projects found when opening logserver session")
	}

	if data.Config.Internal.Log.ClearOnRestart {
		log.Println("clearing known errors")
		for _, projectName := range projects {
			data.DeleteVar(projectName, vars.KnownErrors)
		}
	}

	return sessionDataMap
}
