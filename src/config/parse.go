package config

import (
	"hudson-newey/2rm/src/models"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseConfig(path string) models.Config {
	parsedConfig := models.Config{}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal([]byte(content), &parsedConfig)
	if err != nil {
		log.Fatal(err)
	}

	return parsedConfig
}
