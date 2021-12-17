package logserver

import (
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/barklan/cto/pkg/logserver/types"
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

func openOrEnterSession(
	data *storage.Data,
	sessDataMap map[string]*SessionData,
	projectName string,
) {
	if v, ok := sessDataMap[projectName]; ok {
		atomic.AddUint64(v.Using, 1)
		log.Printf("%q project entered session", projectName)
	} else {
		knownErrorsMutex := new(sync.Mutex)
		knownErrors := make([]types.KnownError, 0)
		knownErrorsRaw := data.GetVar(projectName, vars.KnownErrors)
		if string(knownErrorsRaw) != "" {
			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
				log.Println("failed to unmarshal KnownErrors", err)
			}
		}
		var using uint64 = 1
		sessDataMap[projectName] = &SessionData{
			KnownErrorsMutex: knownErrorsMutex,
			KnownErrors:      knownErrors,
			Using:            &using,
		}
		log.Printf("%q project opened session", projectName)
	}
}

func closeOrLeaveSession(
	data *storage.Data,
	sessDataMap map[string]*SessionData,
	projectName string,
) {
	sessData := sessDataMap[projectName]

	sessData.KnownErrorsMutex.Lock()
	for i, knownError := range sessData.KnownErrors {
		if knownError.LastSeen.Before(time.Now().Add(time.Duration(-12) * time.Hour)) {
			sessData.KnownErrors = remove(sessData.KnownErrors, i)
			log.Println("deteted old error", knownError.OriginBadgerKey)
			break
		}
	}
	data.SetVar(projectName, vars.KnownErrors, sessData.KnownErrors, 48*time.Hour)
	sessData.KnownErrorsMutex.Unlock()

	if *sessDataMap[projectName].Using == uint64(1) {
		delete(sessDataMap, projectName)
		log.Printf("%q project closed session", projectName)
	} else {
		atomic.AddUint64(sessDataMap[projectName].Using, ^uint64(0))
		log.Printf("%q project left session", projectName)
	}
}
