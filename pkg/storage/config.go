package storage

import (
	"encoding/json"
	"log"
	"sync"
)

type Config struct {
	Internal InternalConfig
	P        map[string]int64
	PMutex   sync.Mutex
}

func ReadConfig(data *Data) *Config {
	internal, err := ReadInternalConfig("")
	if err != nil {
		log.Fatal(err)
	}
	projects := ReadAllProjectsConfigs(data)

	config := &Config{
		Internal: internal,
		P:        projects,
		PMutex:   sync.Mutex{},
	}
	return config
}

func ReadAllProjectsConfigs(data *Data) map[string]int64 {
	projects := map[string]int64{}
	projectsRaw := Get(data.DB, "!projects")
	if len(projectsRaw) == 0 {
		return projects
	}
	err := json.Unmarshal(projectsRaw, &projects)
	if err != nil {
		log.Panicln("failed to unmarshal projects")
	}
	return projects
}

func AddProject(data *Data, project string, chatID int64) {
	data.Config.PMutex.Lock()
	data.Config.P[project] = chatID
	byteObj, err := json.Marshal(data.Config.P)
	if err != nil {
		data.Config.PMutex.Unlock()
		log.Panic(err)
	}
	Set(data.DB, "!projects", byteObj)
	data.Config.PMutex.Unlock()
}

func DeleteProject(data *Data, project string) {
	data.Config.PMutex.Lock()
	delete(data.Config.P, project)
	byteObj, err := json.Marshal(data.Config.P)
	if err != nil {
		data.Config.PMutex.Unlock()
		log.Panic(err)
	}
	Set(data.DB, "!projects", byteObj)
	data.Config.PMutex.Unlock()
}
