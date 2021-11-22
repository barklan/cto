package querying

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/barklan/cto/pkg/storage"
)

type Set map[string]struct{}

func GetKnownServices(data *storage.Data, env string) Set {
	knownServicesRaw := data.Get(fmt.Sprintf("knownServices-%s", env))
	knownServices := Set{}
	if string(knownServicesRaw) == "" {
		return knownServices
	}
	err := json.Unmarshal(knownServicesRaw, &knownServices)
	if err != nil {
		data.CSend(fmt.Sprintf("Failed to unmarshal knownServices for %s", env))
	}
	return knownServices
}

func SetKnownServices(data *storage.Data, env string, knownServices Set) {
	data.SetObj(fmt.Sprintf("knownServices-%s", env), knownServices, 10*time.Hour)
}
