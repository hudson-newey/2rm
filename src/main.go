package main

import (
	"fmt"
	"os"
	"strings"

	"hudson-newey/2rm/src/commands"
)

func extractFileNames(args []string) []string {
	fileNames := []string{}

	for _, str := range args {
		if !strings.HasPrefix(str, "-") {
			fileNames = append(fileNames, str)
		}
	}

	return fileNames
}

func main() {
	originalArguments := os.Args[1:];

	// I remove and modify some arguments from the rm command
	args := strings.Join(originalArguments, " ")
	actionedArgs := strings.ReplaceAll(args, " --no-preserve-root", "")

	fileNames := extractFileNames(originalArguments)
	debugStatement := strings.Join(fileNames, " ")
	fmt.Println(debugStatement)

	gitCommand := "rm " + actionedArgs
	commands.Execute(gitCommand)
}
