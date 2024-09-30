package patches

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hudson-newey/2rm/src/commands"
	"hudson-newey/2rm/src/models"
	"hudson-newey/2rm/src/util"
)

const TRASH_DIR_PERMISSIONS = 0755
const HARD_DELETE_CLA = "--hard"
const SOFT_DELETE_CLA = "--soft"
const SILENT_CLA = "--silent"
const DRY_RUN_CLA = "--dry-run"
const BYPASS_PROTECTED_CLA = "--bypass-protected"

func RmPatch(arguments []string, config models.Config) {
	forceHardDelete := util.InArray(arguments, HARD_DELETE_CLA)
	forceSoftDelete := util.InArray(arguments, SOFT_DELETE_CLA)
	silent := util.InArray(arguments, SILENT_CLA)
	dryRun := util.InArray(arguments, DRY_RUN_CLA)
	bypassProtected := util.InArray(arguments, BYPASS_PROTECTED_CLA)

	actionedArgs := removeUnNeededArguments(
		removeDangerousArguments(arguments),
	)

	if shouldPassthrough(actionedArgs) {
		command := "rm " + strings.Join(actionedArgs, " ")
		commands.Execute(command)
		return
	}

	filePaths := extractFilePaths(actionedArgs)
	extractedArguments := extractArguments(actionedArgs)

	// a debug statement is useful for scripts, it provides explicit feedback
	// and prints exactly what files were backed up / moved to the trash can
	// after deletion
	debugStatement := strings.Join(filePaths, " ")
	if !silent {
		fmt.Println(debugStatement)
	}

	if dryRun {
		return
	}

	for _, path := range filePaths {
		absolutePath := relativeToAbsolute(path)
		isTmp := isTmpPath(absolutePath)

		isProtected := config.IsProtected(absolutePath)
		if isProtected && bypassProtected {
			fmt.Println("Cannot delete protected file:", absolutePath)
			fmt.Println("Use the --bypass-protected flag to force deletion")
			continue
		}

		isConfigHardDelete := config.ShouldHardDelete(absolutePath)
		isConfigSoftDelete := config.ShouldSoftDelete(absolutePath)

		if isTmp || forceHardDelete || isConfigHardDelete && !isConfigSoftDelete && !forceSoftDelete {
			hardDelete([]string{path}, extractedArguments)
		} else {
			// this function will return the default soft delete directory
			// if the user has not specified one in their config file
			softDeleteDir := config.SoftDeleteDir()
			softDelete([]string{path}, extractedArguments, softDeleteDir)
		}
	}
}

// sometimes we want to pass through the arguments to the original rm command
// e.g. when executing --help or --version
func shouldPassthrough(arguments []string) bool {
	passthroughArguments := []string{"--help", "--version"}

	for _, arg := range arguments {
		if util.InArray(passthroughArguments, arg) {
			return true
		}
	}

	return false
}

func removeUnNeededArguments(arguments []string) []string {
	unNeededArguments := []string{"-r", HARD_DELETE_CLA, SOFT_DELETE_CLA, SILENT_CLA}
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

// by default, we want to delete files to /tmp/2rm
// however, if the user has specified a different directory in their config file
// we use that instead
func softDelete(filePaths []string, arguments []string, tempDir string) {
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
	command := "rm -r " + strings.Join(arguments, " ") + " " + strings.Join(filePaths, " ")
	commands.Execute(command)
}
