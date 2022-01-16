package logserver

import (
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/barklan/cto/pkg/logserver/types"
	"github.com/barklan/cto/pkg/storage"
	"github.com/barklan/cto/pkg/storage/vars"
)

// TODO this should not select all projects
func openSession(data *storage.Data) map[string]*SessionData {
	sessionDataMap := map[string]*SessionData{}

	projects := make([]string, 0)
	if err := data.R.Select(&projects, "select id from project"); err != nil {
		data.Log.Warn("no projects found when opening logserver session")
	}

	if data.Config.Internal.Log.ClearOnRestart {
		data.Log.Warn("clearing known errors for all projects")
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
		data.Log.Info("project entered session", zap.String("project", projectName))
	} else {
		knownErrorsMutex := new(sync.Mutex)
		knownErrors := make([]types.KnownError, 0)
		knownErrorsRaw := data.GetVar(projectName, vars.KnownErrors)
		if string(knownErrorsRaw) != "" {
			data.Log.Info("known issues found", zap.String("project", projectName))
			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
				data.Log.Error("failed to unmarshal KnownErrors", zap.Error(err))
			}
		}
		var using uint64 = 1
		sessDataMap[projectName] = &SessionData{
			KnownErrorsMutex: knownErrorsMutex,
			KnownErrors:      knownErrors,
			Using:            &using,
		}
		data.Log.Info("project opened session", zap.String("project", projectName))
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
			data.Log.Info("deteted old error", zap.String("key", knownError.OriginBadgerKey))
			break
		}
	}
	data.SetVar(projectName, vars.KnownErrors, sessData.KnownErrors, 48*time.Hour)
	sessData.KnownErrorsMutex.Unlock()

	if *sessDataMap[projectName].Using == uint64(1) {
		delete(sessDataMap, projectName)
		data.Log.Info("project closed session", zap.String("project", projectName))
	} else {
		atomic.AddUint64(sessDataMap[projectName].Using, ^uint64(0))
		data.Log.Info("project left session", zap.String("project", projectName))
	}
}
