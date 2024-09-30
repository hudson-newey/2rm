package config

import (
	"hudson-newey/2rm/src/models"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseConfig(path string) models.Config {
	parsedConfig := models.Config{}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// if the config file does not exist, we want to return an empty config
		// this will act as a default config
		return parsedConfig
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(content, &parsedConfig)
	if err != nil {
		log.Fatal(err)
	}

	return parsedConfig
}
