package config

import "log"

type Config struct {
	Internal            InternalConfig
	P                   map[string]ProjectConfig
	EnvToProjectName    map[string]string
	ChatIDToProjectName map[int64]string
}

func ReadConfig() Config {
	internal, err := ReadInternalConfig("")
	if err != nil {
		log.Fatal(err)
	}
	projects := ReadAllProjectsConfigs()

	envToProjectName := map[string]string{}
	chatIDToProjectName := map[int64]string{}

	for projectName, projectConfig := range projects {
		for _, env := range projectConfig.Envs {
			envToProjectName[env] = projectName
		}
		chatIDToProjectName[projectConfig.TG.ChatID] = projectName
	}

	config := Config{
		Internal:            internal,
		P:                   projects,
		EnvToProjectName:    envToProjectName,
		ChatIDToProjectName: chatIDToProjectName,
	}
	return config
}
