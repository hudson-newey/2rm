package main

import (
	"os"

	"hudson-newey/2rm/src/cli"
	"hudson-newey/2rm/src/config"
	"hudson-newey/2rm/src/models"
	"hudson-newey/2rm/src/rm"
)

func main() {
	originalArguments := os.Args[1:]

	if len(originalArguments) == 0 {
		cli.PrintError("2rm: missing operand\nTry '2rm --help' for more information")
		os.Exit(1)
	}

	homeLocation := os.Getenv("HOME")

	systemConfigPath := "/etc/2rm/config.yml"
	userConfigPath := homeLocation + "/.local/share/2rm/config.yml"
	parsedConfig := models.Config{}

	// TODO: we might want to extend user configs from system configs once the
	// "extend" config property is implemented
	// see: https://github.com/hudson-newey/2rm/issues/9
	_, err := os.Stat(userConfigPath)
	if err == nil {
		parsedConfig = config.ParseConfig(userConfigPath)
	} else {
		// attempt to parse system config
		_, err := os.Stat(systemConfigPath)
		if err == nil {
			parsedConfig = config.ParseConfig(systemConfigPath)
		}
	}

	argumentModel := cli.ParseCliFlags(originalArguments)

	rm.Execute(argumentModel, parsedConfig)
}
