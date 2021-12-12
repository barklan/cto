package querying

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/barklan/cto/pkg/storage"
)

type Set map[string]struct{}

func GetKnownServices(data *storage.Data, project, env string) Set {
	knownServicesRaw := data.GetVar(project, env+".knownServices")
	knownServices := Set{}
	if string(knownServicesRaw) == "" {
		return knownServices
	}
	err := json.Unmarshal(knownServicesRaw, &knownServices)
	if err != nil {
		data.PSend(project, fmt.Sprintf("Failed to unmarshal knownServices for %s.", env))
	}
	return knownServices
}

func SetKnownServices(data *storage.Data, project, env string, knownServices Set) {
	data.SetVar(project, env+".knownServices", knownServices, 120*time.Hour)
}

func GetKnownEnvs(data *storage.Data, project string) Set {
	knownEnvsRaw := data.GetVar(project, "knownEnvs")
	knownEnvs := Set{}
	if string(knownEnvsRaw) == "" {
		return knownEnvs
	}
	err := json.Unmarshal(knownEnvsRaw, &knownEnvs)
	if err != nil {
		data.PSend(project, "Failed to unmarshal knownEnvs.")
	}
	return knownEnvs
}

func SetKnownEnvs(data *storage.Data, project string, knownEnvs Set) {
	data.SetVar(project, "knownEnvs", knownEnvs, 120*time.Hour)
}
