package logserver

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/logserver/types"
	"github.com/barklan/cto/pkg/storage"
)

// TODO recovery v5
// func registerResetErrorsHandler(data *storage.Data, sessionDataMap map[string]*SessionData) {
// 	data.B.Handle("/reset", func(m *tb.Message) {
// 		projectName, ok := bot.VerifySender(data, m)
// 		if !ok {
// 			return
// 		}

// 		sessData := sessionDataMap[projectName]
// 		knownErrors := make([]types.KnownError, 0)

// 		sessData.KnownErrorsMutex.Lock()
// 		sessData.KnownErrors = knownErrors
// 		data.SetObj(
// 			fmt.Sprintf("knownErrors-%s", projectName),
// 			sessData.KnownErrors,
// 			12*time.Hour,
// 		)
// 		sessData.KnownErrorsMutex.Unlock()

// 		data.PSend(projectName, "Error records have been reset.")
// 	})
// }

func openSession(data *storage.Data) map[string]*SessionData {
	sessionDataMap := map[string]*SessionData{}

	projects := make([]string, 0)
	if err := data.R.Select(&projects, "select id from project"); err != nil {
		log.Println("no projects found when opening logserver session")
	}

	for _, projectName := range projects {
		knownErrors := make([]types.KnownError, 0)
		sessionData := &SessionData{
			KnownErrors:      knownErrors,
			KnownErrorsMutex: sync.Mutex{},
		}
		sessionDataMap[projectName] = sessionData
	}

	if data.Config.Internal.Log.ClearOnRestart {
		log.Println("clearing known errors")
		for _, projectName := range projects {
			knownErrors := make([]types.KnownError, 0)
			data.SetObj(fmt.Sprintf("knownErrors-%s", projectName), knownErrors, 1*time.Hour)
		}
	} else {
		for _, projectName := range projects {
			knownErrors := make([]types.KnownError, 0)
			knownErrorsRaw := data.Get(fmt.Sprintf("knownErrors-%s", projectName))
			if string(knownErrorsRaw) != "" {
				if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
					log.Println("failed to unmarshal knownErrors", err)
				}
			}
			sessionDataMap[projectName].KnownErrors = knownErrors
		}
	}

	go func(sessDataMap map[string]*SessionData) {
		ticker := time.NewTicker(1 * time.Minute)
		for {
			<-ticker.C
			for projectName, sessData := range sessDataMap {
				sessData.KnownErrorsMutex.Lock()
				data.SetObj(
					fmt.Sprintf("knownErrors-%s", projectName),
					sessData.KnownErrors,
					5*time.Hour,
				) // ttl in case goroutine fails
				sessData.KnownErrorsMutex.Unlock()
			}
		}
	}(sessionDataMap)

	go func(sessDataMap map[string]*SessionData) {
		ticker := time.NewTicker(5 * time.Minute)
		for {
			<-ticker.C
			for _, sessData := range sessDataMap {
				sessData.KnownErrorsMutex.Lock()
				for i, knownError := range sessData.KnownErrors {
					if knownError.LastSeen.Before(time.Now().Add(time.Duration(-12) * time.Hour)) {
						sessData.KnownErrors = remove(sessData.KnownErrors, i)
						log.Println("deteted old error", knownError.OriginBadgerKey)
						break
					}
				}
				sessData.KnownErrorsMutex.Unlock()
			}
		}
	}(sessionDataMap)

	// TODO recovery v5
	// registerResetErrorsHandler(data, sessionDataMap)

	return sessionDataMap
}
