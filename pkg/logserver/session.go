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
) *SessionData {
	if v, ok := sessDataMap[projectName]; ok {
		atomic.AddUint64(v.Using, 1)
		data.Log.Info("project entered session", zap.String("project", projectName))
	} else {
		Mutex := new(sync.Mutex)
		knownErrors := make(map[string]types.KnownError)
		knownErrorsRaw := data.GetVar(projectName, vars.KnownErrors)
		if string(knownErrorsRaw) != "" {
			data.Log.Info("known issues found", zap.String("project", projectName))
			if err := json.Unmarshal(knownErrorsRaw, &knownErrors); err != nil {
				data.Log.Error("failed to unmarshal KnownErrors", zap.Error(err))
			}
		}
		var using uint64 = 1
		sessDataMap[projectName] = &SessionData{
			Mutex:       Mutex,
			KnownErrors: knownErrors,
			Using:       &using,
		}
		data.Log.Info("project opened session", zap.String("project", projectName))
	}
	return sessDataMap[projectName]
}

func closeOrLeaveSession(
	data *storage.Data,
	sessDataMap map[string]*SessionData,
	projectName string,
) {
	sessData, ok := sessDataMap[projectName]
	if !ok {
		data.Log.Warn("project session not found when trying to close", zap.String("pid", projectName))
		return
	}

	sessData.Mutex.Lock()
	for i, knownError := range sessData.KnownErrors {
		if knownError.LastSeen.Before(time.Now().Add(time.Duration(-12) * time.Hour)) {
			delete(sessData.KnownErrors, i)
			data.Log.Info("deteted old error", zap.String("key", knownError.OriginBadgerKey))
			break
		}
	}
	data.SetVar(projectName, vars.KnownErrors, sessData.KnownErrors, 24*time.Hour)
	knownErrorsJson, err := json.Marshal(sessData.KnownErrors)
	if err != nil {
		data.Log.Error("failed to marshal knownErrors", zap.String("project", projectName))
		if e := data.Cache.SetVar(projectName, vars.KnownErrors, knownErrorsJson, 48*time.Hour); e != nil {
			data.Log.Error("failed to set knownErrors to cache", zap.String("project", projectName), zap.Error(e))
		}
	}
	sessData.Mutex.Unlock()

	// FIXME nil pointer deref here
	if *sessDataMap[projectName].Using == uint64(1) {
		delete(sessDataMap, projectName)
		data.Log.Info("project closed session", zap.String("project", projectName))
	} else {
		atomic.AddUint64(sessDataMap[projectName].Using, ^uint64(0))
		data.Log.Info("project left session", zap.String("project", projectName))
	}
}
