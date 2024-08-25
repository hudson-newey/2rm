package patches

import (
	"fmt"
	"strings"

	"hudson-newey/2rm/src/commands"
	"hudson-newey/2rm/src/util"
)

func RmPatch(arguments []string) {
	actionedArgs := removeDangerousArguments(arguments)

	fileNames := extractFileNames(actionedArgs)
	debugStatement := strings.Join(fileNames, " ")
	fmt.Println(debugStatement)

	command := "rm " + strings.Join(actionedArgs, " ")
	commands.Execute(command)
}

func removeDangerousArguments(arguments []string) []string {
	// I have excluded the root slash as a forbidden argument just incase
	// you make a typo like rm ./myDirectory /
	// when you were just trying to delete myDirectory
	forbiddenArguments := []string{"/", "--no-preserve-root"}
	returnedArguments := []string{}

	for _, arg := range arguments {
		isForbidden := util.InArray(forbiddenArguments, arg)
		if !isForbidden {
			returnedArguments = append(returnedArguments, arg)
		}
	}

	return returnedArguments
}

func extractFileNames(args []string) []string {
	fileNames := []string{}

	for _, str := range args {
		if !strings.HasPrefix(str, "-") {
			fileNames = append(fileNames, str)
		}
	}

	return fileNames
}
