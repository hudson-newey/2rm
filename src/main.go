package main

import (
	"os"

	"hudson-newey/2rm/src/config"
	"hudson-newey/2rm/src/models"
	"hudson-newey/2rm/src/patches"
)

func main() {
	originalArguments := os.Args[1:]

	homeLocation := os.Getenv("HOME")
	defaultConfigDirectory := homeLocation + "/.local/share/2rm/config.yml"
	parsedConfig := models.Config{}

	_, err := os.Stat(defaultConfigDirectory)
	if err == nil {
		parsedConfig = config.ParseConfig(defaultConfigDirectory)
	}

	patches.RmPatch(originalArguments, parsedConfig)
}
