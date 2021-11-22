package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
)

type InternalConfig struct {
	TG struct {
		BotToken       string `yaml:"bot_token"`
		Boss           string `yaml:"boss"`
		ClearOnRestart bool   `yaml:"clear_on_restart"`
	} `yaml:"tg"`
	DataPath      string `yaml:"data_path"`
	JWTExpHours   int    `yaml:"jwt_exp_hours"`
	MagicJWTToken string `yaml:"magic_jwt_token"`
	Log           struct {
		ClearOnRestart      bool    `yaml:"clear_on_restart"`
		ServiceHostname     string  `yaml:"service_hostname"`
		RetentionHours      int     `yaml:"retention_hours"`
		SimilarityThreshold float64 `yaml:"similarity_threshold"`
	} `yaml:"log"`
}

func ReadInternalConfig(path string) (InternalConfig, error) {
	var config_path string
	if path == "" {
		if _, ok := os.LookupEnv("CTO_LOCAL_ENV"); ok {
			config_path = "environment/local.yml"
		} else {
			config_path = "/app/config/internal.yml"
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

	fmt.Printf("ReadConfig: %# v", pretty.Formatter(config))

	return config, nil
}
