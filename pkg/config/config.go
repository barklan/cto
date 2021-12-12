package config

import "log"

type Config struct {
	Internal            InternalConfig
	P                   map[string]ProjectConfig
	ChatIDToProjectName map[int64]string
}

func ReadConfig() Config {
	internal, err := ReadInternalConfig("")
	if err != nil {
		log.Fatal(err)
	}
	projects := ReadAllProjectsConfigs()

	chatIDToProjectName := map[int64]string{}

	for projectName, projectConfig := range projects {
		chatIDToProjectName[projectConfig.TG.ChatID] = projectName
	}

	config := Config{
		Internal:            internal,
		P:                   projects,
		ChatIDToProjectName: chatIDToProjectName,
	}
	return config
}
