package main

import (
	"fmt"
	"os"

	"hudson-newey/2rm/src/config"
	"hudson-newey/2rm/src/models"
	"hudson-newey/2rm/src/rm"
)

func main() {
	originalArguments := os.Args[1:]

	if (len(originalArguments) == 0) {
		fmt.Println("2rm: missing operand\nTry '2rm --help' for more information")
		os.Exit(1)
	}

	homeLocation := os.Getenv("HOME")
	defaultConfigDirectory := homeLocation + "/.local/share/2rm/config.yml"
	parsedConfig := models.Config{}

	_, err := os.Stat(defaultConfigDirectory)
	if err == nil {
		parsedConfig = config.ParseConfig(defaultConfigDirectory)
	}

	rm.Execute(originalArguments, parsedConfig)
}
