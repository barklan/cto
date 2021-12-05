package logserver

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/barklan/cto/pkg/storage"
)

func openSession(data *storage.Data) map[string]*SessionData {
	sessionDataMap := map[string]*SessionData{}

	for projectName := range data.Config.P {
		knownErrors := make([]KnownError, 0)
		sessionData := &SessionData{
			KnownErrors:      knownErrors,
			KnownErrorsMutex: sync.Mutex{},
		}
		sessionDataMap[projectName] = sessionData
	}

	if data.Config.Internal.Log.ClearOnRestart == true {
		log.Println("clearing known errors")
		for projectName := range data.Config.P {
			knownErrors := make([]KnownError, 0)
			data.SetObj(fmt.Sprintf("knownErrors-%s", projectName), knownErrors, 1*time.Hour)
		}
	} else {
		for projectName := range data.Config.P {
			knownErrors := make([]KnownError, 0)
			knownErrorsRaw := data.Get(fmt.Sprintf("knownErrors-%s", projectName))
			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
				log.Println("failed to unmarshal knownErrors", err)
			}
			sessionDataMap[projectName].KnownErrors = knownErrors
		}
	}

	go func(sessionDataMap map[string]*SessionData) {
		ticker := time.NewTicker(72 * time.Hour)
		for {
			<-ticker.C
			for projectName, sessData := range sessionDataMap {
				knownErrors := make([]KnownError, 0)
				sessData.KnownErrorsMutex.Lock()
				sessionDataMap[projectName].KnownErrors = knownErrors
				sessData.KnownErrorsMutex.Unlock()
				data.PSend(projectName, fmt.Sprintf(
					"Error records have been forcefully cleaned. Next cleaning in 3 days.",
				))
			}
		}
	}(sessionDataMap)

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

	return sessionDataMap
}
