package storage

import (
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Internal InternalConfig
}

func ReadConfig(data *Data) *Config {
	internal, err := ReadInternalConfig("")
	if err != nil {
		log.Fatal(err)
	}

	config := &Config{
		Internal: internal,
	}
	return config
}
