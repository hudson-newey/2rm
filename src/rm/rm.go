package rm

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

func Execute(args models.CliOptions, config models.Config) {
	sanitizedArgs := removeDangerousArguments(args.RawArguments)

	if args.RequestingHelp {
		cli.PrintHelp()
		return
	} else if args.RequestingVersion {
		cli.PrintVersion()
		return
	}

	filePaths := extractFilePaths(sanitizedArgs)
	deletePaths(filePaths, config, args)

	if args.ShouldNotify {
		fileNames := strings.Join(filePaths, ", ")
		err := beeep.Notify("2rm", "Completed deletion '"+fileNames+"'", "")
		if err != nil {
			panic(err)
		}
	}
}

func removeDangerousArguments(arguments []string) []string {
	// I have excluded the root slash as a forbidden argument just in case
	// you make a typo like rm ./myDirectory /
	// when you were just trying to delete myDirectory
	// If you really have to delete your root directory consider using the GNU
	// rm command
	forbiddenArguments := []string{"/"}
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

	for _, path := range input {
		if !strings.HasPrefix(path, "-") {
			if util.IsDirectory(path) && !strings.HasSuffix(path, "/") {
				filePaths = append(filePaths, path+"/")
			} else {
				filePaths = append(filePaths, path)
			}
		} else {
			isSupportedCla := false
			for _, argument := range cli.SupportedCliArguments {
				cliFormat := "--" + argument
				shortCliFormat := "-" + argument
				if cliFormat == path || shortCliFormat == path {
					isSupportedCla = true
				}
			}

			if isSupportedCla {
				continue
			}

			errorMessage := "unrecognized option '" + path + "'\nTry '2rm --help' for more information."
			if util.PathExists(path) {
				errorMessage = "unrecognized option '" + path + "'\n" +
					"Try '2rm ./" + path + "' to remove the file '" + path + "'.\n" +
					"Try '2rm --help' for more information."
			}

			// I (personally) don't like this behavior of only erroring one file at a
			// time. However, to maintain compatibility with GNU rm I have replicated
			// this behavior of only showing the first error every run.
			cli.PrintError(errorMessage)
			os.Exit(1)
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
		cli.PrintErrorValue(err)
		os.Exit(2)
	}

	return absolutePath
}

func backupFileName(path string) string {
	result := strings.ReplaceAll(path, ".", "-")
	result = strings.ReplaceAll(result, "/", "_")
	return result + ".bak"
}

func deletePaths(paths []string, config models.Config, args models.CliOptions) {
	isInteractiveGroup := args.IsGroupInteractive && len(paths) >= config.InteractiveThreshold()
	isInteractive := args.IsInteractive || isInteractiveGroup

	removedFiles := []string{}
	for _, path := range paths {
		if !util.PathExists(path) {
			cli.PrintError("Cannot remove '" + path + "': No such file or directory")
			continue
		}

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
		if isProtected && args.BypassProtected {
			// TODO: maybe these two print lines should be combined
			cli.PrintError("Cannot delete protected file:" + absolutePath)
			cli.PrintError("Use the --bypass-protected flag to force deletion")
			continue
		}

		if args.OnlyEmptyDirs {
			isDirEmpty := util.IsDirectoryEmpty(absolutePath)
			if !isDirEmpty {
				cli.PrintError("cannot remove '" + path + "': Directory not empty")
				continue
			}
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
		if args.Overwrite || isConfigOverwrite && !args.DryRun {
			overwriteFile(absolutePath)
		}

		shouldHardDelete := isTmp || args.HardDelete || isConfigHardDelete && !isConfigSoftDelete && !args.SoftDelete

		deletePath(absolutePath, shouldHardDelete, args.DryRun, config)

		if args.Verbose && !args.Silent {
			fmt.Printf("removed '%s'\n", path)
		}

		removedFiles = append(removedFiles, path)
	}

	removedFilesLen := len(removedFiles)
	for i, path := range removedFiles {
		fmt.Print(path)

		// TODO: we might want to make this the shells IFS value
		if i != removedFilesLen {
			fmt.Print(" ")
		}
	}

	// do an empty print line last so that when the shell returns to a
	// prompt the prompt will be on its own line
	if removedFilesLen > 0 {
		fmt.Println()
	}
}

func deletePath(path string, hard bool, dryRun bool, config models.Config) {
	// I break out during dry runs at the last possible point so that dry
	// runs are as close to real runs as possible
	// all the same debug information, deletion prompts and errors will be
	// shown during dry runs
	if dryRun {
		return
	}

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
		fmt.Println("Failed to create soft delete can entry in", backupDirectory)
		fmt.Println("Continue? (y/n)")

		var response string
		fmt.Scanln(&response)

		// unless the user explicitly states that they want to continue
		// without a backup, we want to exit
		// everything other than a "y"/"yes" response will not delete the file
		if response != "y" && response != "yes" {
			cli.PrintError("Exiting without removing file(s).")
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
		// because we have replicated the directory structure in the soft delete
		// directory we don't need to keep the original directory
		//
		// TODO: we should probably keep the original directory so that the
		// same file permissions and other edge cases are carried across
		hardDelete(filePath)
		return
	}

	absoluteSrcPath := relativeToAbsolute(filePath)

	backupFileName := backupFileName(filePath)
	backupLocation := backupDirectory + "/" + backupFileName

	err = util.MovePath(absoluteSrcPath, backupLocation)
	if err != nil {
		cli.PrintError("Error soft deleting file:" + err.Error())
		return
	}

	hardDelete(absoluteSrcPath)
}

func hardDelete(filePath string) {
	err := os.RemoveAll(filePath)
	if err != nil {
		cli.PrintError("Error deleting file:" + err.Error())
	}
}

func overwriteFile(filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		cli.PrintError("Error opening file:" + err.Error())
		os.Exit(2)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		cli.PrintError("Error reading file info:" + err.Error())
		os.Exit(2)
	}

	fileSize := fileInfo.Size()
	zeroBytes := make([]byte, fileSize)

	_, err = file.WriteAt(zeroBytes, 0)
	if err != nil {
		cli.PrintError("Error writing to file:" + err.Error())
		os.Exit(2)
	}
}
