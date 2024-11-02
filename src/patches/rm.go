package patches

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"hudson-newey/2rm/src/cli"
	"hudson-newey/2rm/src/models"
	"hudson-newey/2rm/src/util"

	"github.com/gen2brain/beeep"
)

const TRASH_DIR_PERMISSIONS = 0755

func RmPatch(arguments []string, config models.Config) {
	silent := util.InArray(arguments, cli.SILENT_CLA)
	dryRun := util.InArray(arguments, cli.DRY_RUN_CLA)
	shouldNotify := util.InArray(arguments, cli.NOTIFICATION_CLA)

	requestingHelp := util.InArray(arguments, cli.HELP_CLA)
	requestingVersion := util.InArray(arguments, cli.VERSION_CLA)

	actionedArgs := removeDangerousArguments(arguments)

	if requestingHelp {
		cli.PrintHelp()
		return
	} else if requestingVersion {
		cli.PrintVersion()
		return
	}

	filePaths := extractFilePaths(actionedArgs)

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

	deletePaths(filePaths, config, arguments)

	if shouldNotify {
		fileNames := strings.Join(filePaths, ", ")
		err := beeep.Notify("2rm", "Finished deletion request '"+fileNames+"'", "")
		if err != nil {
			panic(err)
		}
	}
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

func deletePaths(paths []string, config models.Config, arguments []string) {
	forceHardDelete := util.InArray(arguments, cli.HARD_DELETE_CLA)
	forceSoftDelete := util.InArray(arguments, cli.SOFT_DELETE_CLA)
	bypassProtected := util.InArray(arguments, cli.BYPASS_PROTECTED_CLA)
	overwrite := util.InArray(arguments, cli.OVERWRITE_CLA)

	hasInteraciveCla := util.InArray(arguments, cli.INTERACTIVE_CLA)
	hasGroupInteractiveCla := util.InArray(arguments, cli.INTERACTIVE_GROUP_CLA)
	isInteractiveGroup := hasGroupInteractiveCla && len(paths) >= config.InteractiveThreshold()
	isInteractive := hasInteraciveCla || isInteractiveGroup

	for _, path := range paths {
		if isInteractive {
			fmt.Println("Are you sure you want to delete", path, "? (y/n)")
			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "yes" {
				fmt.Println("Skipping file", path)
				continue
			}
		}

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

		// overwriting a file is not exclusive to hard/soft deletes
		// meaning that you can overwrite the contents of a file with zeros and
		// also soft delete it
		// I have made this decision because I think soft-deleting an
		// overwritten file has auditing/logging use cases
		// e.g. Who deleted this file? When was it deleted?
		// if we hard deleted the file, we would lose this information
		isConfigOverwrite := config.ShouldOverwrite(absolutePath)
		if overwrite || isConfigOverwrite {
			overwriteFile(absolutePath)
		}

		shouldHardDelete := isTmp || forceHardDelete || isConfigHardDelete && !isConfigSoftDelete && !forceSoftDelete

		deletePath(absolutePath, shouldHardDelete, config)
	}
}

func deletePath(path string, hard bool, config models.Config) {
	if hard {
		hardDelete(path)
	} else {
		softDeleteStart(path, config)
	}
}

func softDeleteStart(filePath string, config models.Config) {
	// this function will return the default soft delete directory
	// if the user has not specified one in their config file
	softDeletePath := config.SoftDeleteDir()
	softDelete(filePath, softDeletePath, "")
}

// by default, we want to delete files to /tmp/2rm
// however, if the user has specified a different directory in their config file
// we use that instead
func softDelete(filePath string, tempDir string, backupDirectory string) {
	if backupDirectory == "" {
		deletedTimestamp := time.Now().Format(time.RFC3339)
		backupDirectory = tempDir + deletedTimestamp
	}

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
			return
		}
	}

	isDirectory := util.IsDirectory(filePath)
	if isDirectory {
		// recursively delete all files in the directory
		// before deleting the directory itself
		directoryFiles := util.ListFiles(filePath)

		for _, file := range directoryFiles {
			softDelete(file, tempDir, backupDirectory)
		}

		// hard delete the directory itself
		// because we have replicated the directory structure in the trash can
		// we don't need to keep the original directory
		//
		// TODO: we should probably keep the original directory so that the
		// same file permissions and other edge cases are carried across
		hardDelete(filePath)
		return
	}

	absoluteSrcPath := relativeToAbsolute(filePath)

	backupFileName := backupFileName(filePath)
	backupLocation := backupDirectory + "/" + backupFileName

	err = util.CopyFile(absoluteSrcPath, backupLocation)
	if err != nil {
		fmt.Println("Error moving file to trash:", err)
		return
	}

	err = os.Remove(absoluteSrcPath)
	if err != nil {
		fmt.Println("Error deleting file:", err)

		// pause the program so the user can see the error message
		fmt.Scanln()
	}
}

func hardDelete(filePath string) {
	err := os.RemoveAll(filePath)
	if err != nil {
		fmt.Println("Error deleting file:", err)
	}
}

func overwriteFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(2)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		os.Exit(2)
	}

	fileSize := fileInfo.Size()
	zeroBytes := make([]byte, fileSize)

	_, err = file.WriteAt(zeroBytes, 0)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		os.Exit(2)
	}
}
