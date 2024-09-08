package patches

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hudson-newey/2rm/src/commands"
	"hudson-newey/2rm/src/util"
)

const TRASH_DIR_PERMISSIONS = 0755

func RmPatch(arguments []string) {
	actionedArgs := removeUnNeededArguments(
		removeDangerousArguments(arguments),
	)

	filePaths := extractFilePaths(actionedArgs)
	extractedArguments := extractArguments(actionedArgs)

	// a debug statement is useful for scripts, it provides explicit feedback
	// and prints exactly what files were backed up / moved to the trash can
	// after deletion
	debugStatement := strings.Join(filePaths, " ")
	fmt.Println(debugStatement)

	for _, path := range filePaths {
		absolutePath := relativeToAbsolute(path)
		isTmp := isTmpPath(absolutePath)

		// TODO: this should probably support batch deletions
		if isTmp {
			hardDelete([]string{path}, extractedArguments)
		} else {
			softDelete([]string{path}, extractedArguments)
		}
	}
}

func removeUnNeededArguments(arguments []string) []string {
	unNeededArguments := []string{"-r"}
	returnedArguments := []string{}

	for _, arg := range arguments {
		if !util.InArray(unNeededArguments, arg) {
			returnedArguments = append(returnedArguments, arg)
		}
	}

	return returnedArguments
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

func extractFilePaths(input []string) []string {
	filePaths := []string{}

	for _, str := range input {
		if !strings.HasPrefix(str, "-") {
			filePaths = append(filePaths, str)
		}
	}

	return filePaths
}

func extractArguments(input []string) []string {
	arguments := []string{}

	for _, str := range input {
		if strings.HasPrefix(str, "-") {
			arguments = append(arguments, str)
		}
	}

	return arguments
}

func isTmpPath(absolutePath string) bool {
	return strings.HasPrefix(absolutePath, "/tmp")
}

func relativeToAbsolute(path string) string {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(3)
	}

	return absolutePath
}

func backupFileName(path string) string {
	result := strings.ReplaceAll(path, ".", "-")
	result = strings.ReplaceAll(result, "/", "_")
	return result + ".bak"
}

func softDelete(filePaths []string, arguments []string) {
	tempDir := "/tmp/2rm/"
	deletedTimestamp := time.Now().Format(time.RFC3339)
	backupDirectory := tempDir + deletedTimestamp

	err := os.MkdirAll(backupDirectory, TRASH_DIR_PERMISSIONS)
	if err != nil {
		fmt.Println("Failed to create trash can entry in", backupDirectory)
		fmt.Println("Continue? (y/n)")

		var response string
		fmt.Scanln(&response)

		// unless the user explicitly states that they want to continue
		// without a backup, we want to exit
		// everything other than a "y"/"yes" response will not delete the file
		if response != "y" && response != "yes" {
			fmt.Println("Exiting without removing file(s).")
			os.Exit(1)
		}
	}

	commandArguments := strings.Join(arguments, " ")

	for _, path := range filePaths {
		absoluteSrcPath := relativeToAbsolute(path)

		backupFileName := backupFileName(path)
		backupLocation := backupDirectory + "/" + backupFileName

		moveCommand := "mv " + commandArguments + " " + absoluteSrcPath + " " + backupLocation
		commands.Execute(moveCommand)
	}
}

func hardDelete(filePaths []string, arguments []string) {
	command := "rm " + strings.Join(arguments, " ") + " " + strings.Join(filePaths, " ")
	commands.Execute(command)
}
