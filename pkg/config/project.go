package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/kr/pretty"
	"gopkg.in/yaml.v3"
)

type ProjectConfig struct {
	ProjectID string   `yaml:"project_id"`
	Envs      []string `yaml:"envs"`
	TG        struct {
		ChatID int64 `yaml:"chat_id"`
	} `yaml:"tg"`
	Checks struct {
		GitLab struct {
			ProjectID           int    `yaml:"project_id"`
			APIToken            string `yaml:"api_token"`
			FailedPipelinesMain bool   `yaml:"failed_pipelines_main"`
			MRApprovals         bool   `yaml:"mr_approvals"`
		} `yaml:"gitlab"`
		Traefik struct {
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		}
		SimpleURLChecks []string `yaml:"simple_url_checks"`
		SLA             string   `yaml:"sla"`
	} `yaml:"checks"`
	Backups struct {
		DB struct {
			Enable          bool   `yaml:"enable"`
			Continuous      bool   `yaml:"continuous"`
			ContinuousPath  string `yaml:"continuous_path"`
			IntervalMinutes int    `yaml:"interval_minutes"`
			SSHUser         string `yaml:"ssh_user"`
			SSHHostname     string `yaml:"ssh_hostname"`
			SSHKeyFilename  string `yaml:"ssh_key_filename"`
			Driver          string `yaml:"driver"`
			ContainerName   string `yaml:"container_name"`
			Database        string `yaml:"database"`
		} `yaml:"db"`
	} `yaml:"backups"`
	Log struct {
		RetentionHours      int     `yaml:"retention_hours"`
		SimilarityThreshold float64 `yaml:"similarity_threshold"`
	} `yaml:"log"`
}

func ReadProjectConfig(path string) (ProjectConfig, error) {
	var config_path string
	if path == "" {
		if v, ok := os.LookupEnv("CTO_CONFIG_PATH"); ok {
			config_path = v
		} else {
			config_path = "/app/config/cto.yml"
		}
	} else {
		config_path = path
	}

	var config ProjectConfig
	if _, err := os.Stat(config_path); os.IsNotExist(err) {
		return ProjectConfig{}, fmt.Errorf("Config must exist.")
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

func ReadAllProjectsConfigs() map[string]ProjectConfig {
	config_path := "/app/config"
	if _, ok := os.LookupEnv("CTO_LOCAL_ENV"); ok {
		config_path = "environment"
	}
	items, _ := ioutil.ReadDir(config_path)

	configs := map[string]ProjectConfig{}
	for _, item := range items {
		filename := item.Name()
		isLocal := strings.Contains(filename, "_local")

		if _, ok := os.LookupEnv("CTO_LOCAL_ENV"); ok {
			if filename != "internal.yml" && filename != "local.yml" && isLocal {
				projectConfig, err := ReadProjectConfig(config_path + "/" + filename)
				if err != nil {
					log.Printf("failed to read project config %s, %v:", filename, err)
					continue
				}

				configs[projectConfig.ProjectID] = projectConfig
			}
		} else {
			if filename != "internal.yml" && filename != "local.yml" && !isLocal {
				projectConfig, err := ReadProjectConfig(config_path + "/" + filename)
				if err != nil {
					log.Printf("failed to read project config %s, %v:", filename, err)
					continue
				}

				configs[projectConfig.ProjectID] = projectConfig
			}
		}

	}
	log.Println("configs:", configs)
	return configs
}
