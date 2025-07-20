package rm

import (
	"fmt"
	"strings"

	"hudson-newey/2rm/src/cli"
	"hudson-newey/2rm/src/models"
	"hudson-newey/2rm/src/rm/hardd"
	"hudson-newey/2rm/src/rm/overwrited"
	"hudson-newey/2rm/src/rm/softd"
	"hudson-newey/2rm/src/util"

	"github.com/gen2brain/beeep"
)

// When using the --interactive=once command line flag, we prompt the user on
// the first file that they delete, but do not prompt them on any further
// deletions.
// This variable exists to keep track of if we have sent that initial prompt.
// TODO: Module level state is gross. I should refactor this out.
var doneInitialInteractive bool = false

func ProcessDeletion(args models.CliOptions, config models.Config) {
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

func deletePaths(paths []string, config models.Config, args models.CliOptions) {
	isInteractiveGroup := args.IsGroupInteractive && len(paths) >= config.InteractiveThreshold()
	isInteractive := args.IsInteractive || isInteractiveGroup

	removedFiles := []string{}
	initialDevice, err := util.GetFileDevice(".")
	if err != nil {
		cli.PrintError(fmt.Sprintf("Error getting device for current directory: %v", err))
		return
	}

	for _, path := range paths {
		if !util.PathExists(path) {
			cli.PrintError("Cannot remove '" + path + "': No such file or directory")
			continue
		}

		if args.OneFileSystem {
			fileDevice, err := util.GetFileDevice(path)
			if err != nil {
				cli.PrintError(fmt.Sprintf("Error getting device for file %s: %v", path, err))
				continue
			}
			if fileDevice != initialDevice {
				cli.PrintError(fmt.Sprintf("Cannot remove '%s': It is on a different file system.", path))
				continue
			}
		}

		if isInteractive || (args.IsOnceInteractive && !doneInitialInteractive) {
			fmt.Println("Are you sure you want to delete", path, "? (y/n)")
			var response string
			fmt.Scanln(&response)

			// TODO: why is this being set even when the -I and -i flags are used?
			// this is smell that it should be somewhere else
			doneInitialInteractive = true

			if response != "y" && response != "yes" {
				fmt.Println("Skipping file", path)
				continue
			}
		}

		absolutePath := util.RelativeToAbsolute(path)
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
			overwrited.OverwriteFile(absolutePath)
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
		hardd.HardDelete(path)
	} else {
		softd.SoftDeleteStart(path, config)
	}
}
