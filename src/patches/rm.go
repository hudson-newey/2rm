package patches

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hudson-newey/2rm/src/util"
)

const TRASH_DIR_PERMISSIONS = 0755

func RmPatch(arguments []string) {
	actionedArgs := removeDangerousArguments(arguments)

	filePaths := extractFilePaths(actionedArgs)

	// TODO: this should probably allow hard deletion of /tmp files
	createBackups(filePaths)

	// a debug statement is useful for scripts, it provides explicit feedback
	// and prints exactly what files were backed up / moved to the trash can
	// after deletion
	debugStatement := strings.Join(filePaths, " ")
	fmt.Println(debugStatement)

	// command := "rm " + strings.Join(actionedArgs, " ")
	// commands.Execute(command)
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

func extractFilePaths(args []string) []string {
	filePaths := []string{}

	for _, str := range args {
		if !strings.HasPrefix(str, "-") {
			filePaths = append(filePaths, str)
		}
	}

	return filePaths
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

func createBackups(filePaths []string) {
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

	for _, path := range filePaths {
		absoluteSrcPath := relativeToAbsolute(path)

		backupFileName := backupFileName(path)
		backupLocation := backupDirectory + "/" + backupFileName

		sourceFile, err := os.Open(absoluteSrcPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		defer sourceFile.Close()

		destFile, err := os.Create(backupLocation)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		defer destFile.Close()

		_, err = io.Copy(sourceFile, destFile)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			fmt.Println("Failed to copy file for", path)
			fmt.Println("Continue? (y/n)")

			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "yes" {
				fmt.Println("Exiting without removing file(s).")
				os.Exit(2)
			}
		}
	}
}
