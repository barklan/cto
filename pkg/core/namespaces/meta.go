package namespaces

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/barklan/cto/pkg/storage"
	"github.com/barklan/cto/pkg/storage/vars"
	"go.uber.org/zap"
)

type Set map[string]struct{}

func GetKnownServices(data *storage.Data, project, env string) Set {
	knownServicesRaw := data.GetVar(project, env+vars.KnownServices)
	knownServices := Set{}
	if string(knownServicesRaw) == "" {
		return knownServices
	}

	err := json.Unmarshal(knownServicesRaw, &knownServices)
	if err != nil {
		data.InternalAlert(fmt.Sprintf("Failed to unmarshal knownServices for %s.", env))
	}
	return knownServices
}

// TODO json marshalling twice here
func SetKnownServices(data *storage.Data, project, env string, knownServices Set) {
	data.SetVar(project, env+vars.KnownServices, knownServices, 120*time.Hour)
	raw, err := json.Marshal(knownServices)
	if err != nil {
		data.Log.Error("failed to marshal knownServices", zap.Error(err))
		return
	}
	if err := data.Cache.SetVar(project, env+vars.KnownServices, raw, 120*time.Hour); err != nil {
		data.Log.Error("error caching knownServices", zap.String("project", project), zap.Error(err))
		return
	}
}

func GetKnownEnvs(data *storage.Data, project string) Set {
	knownEnvsRaw := data.GetVar(project, vars.KnownEnvs)
	knownEnvs := Set{}
	if string(knownEnvsRaw) == "" {
		return knownEnvs
	}
	err := json.Unmarshal(knownEnvsRaw, &knownEnvs)
	if err != nil {
		data.InternalAlert(fmt.Sprintf("Failed to unmarshal knownEnvs for project %q", project))
	}
	return knownEnvs
}

// TODO json marshalling twice here
func SetKnownEnvs(data *storage.Data, project string, knownEnvs Set) {
	data.SetVar(project, vars.KnownEnvs, knownEnvs, 120*time.Hour)
	raw, err := json.Marshal(knownEnvs)
	if err != nil {
		data.Log.Error("failed to marshal knownEnvs", zap.Error(err))
		return
	}
	if err := data.Cache.SetVar(project, vars.KnownEnvs, raw, 120*time.Hour); err != nil {
		data.Log.Error("error caching knownEnvs", zap.String("project", project), zap.Error(err))
		return
	}
}

func SetLastRefresh(data *storage.Data, project string) {
	data.SetVar(project, vars.MetaLastRefesh, time.Now(), 48*time.Hour)
}

func GetLastRefresh(data *storage.Data, project string) time.Duration {
	b := data.GetVar(project, vars.MetaLastRefesh)
	var last time.Time
	if err := json.Unmarshal(b, &last); err != nil {
		data.Log.Error("failed to unmarshal last meta time", zap.Error(err))
		return 24 * time.Hour
	}
	return time.Since(last)
}

func Clear(data *storage.Data, project string) {
	envs := GetKnownEnvs(data, project)
	for k := range envs {
		SetKnownServices(data, project, k, Set{})
	}
	SetKnownEnvs(data, project, Set{})
}
