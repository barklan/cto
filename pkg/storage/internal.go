package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var IChatState = "IChatState"

type InternalConfig struct {
	TG struct {
		BotToken   string `yaml:"bot_token"`
		BossChatID int64  `yaml:"boss_chat_id"`
	} `yaml:"tg"`
	JWTExpHours int `yaml:"jwt_exp_hours"`
	Log         struct {
		ClearOnRestart      bool    `yaml:"clear_on_restart"`
		ServiceHostname     string  `yaml:"service_hostname"`
		RetentionHours      int     `yaml:"retention_hours"`
		SimilarityThreshold float64 `yaml:"similarity_threshold"`
	} `yaml:"log"`
}

func ReadInternalConfig(path string) (InternalConfig, error) {
	var config_path string
	if path == "" {
		configEnvironment, ok := os.LookupEnv("CONFIG_ENV")
		if !ok {
			log.Panic("Config environment variable CONFIG_ENV must be specified.")
		}
		if configEnvironment == "dev" {
			config_path = "environment/local.yml"
		} else if configEnvironment == "devdocker" {
			config_path = "/app/config/local.yml"
		} else if configEnvironment == "prod" {
			config_path = "/app/config/internal.yml"
		} else {
			log.Panic("CONFIG_ENV must be either dev or prod.")
		}
	} else {
		config_path = path
	}

	var config InternalConfig
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		return InternalConfig{}, fmt.Errorf("Config must exist.")
	}

	content, err := ioutil.ReadFile(config_path)
	if err != nil {
		log.Panic("Failed to read config file.")
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		return config, fmt.Errorf("Failed to parse config file. %v", err)
	}

	log.Printf("ReadConfig: %#v\n", fmt.Sprint(config))

	return config, nil
}
