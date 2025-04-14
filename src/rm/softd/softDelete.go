package softd

import (
	"fmt"
	"hudson-newey/2rm/src/cli"
	"hudson-newey/2rm/src/models"
	"hudson-newey/2rm/src/rm/hardd"
	"hudson-newey/2rm/src/util"
	"os"
	"strings"
	"time"
)

const TRASH_DIR_PERMISSIONS = 0755

func SoftDeleteStart(filePath string, config models.Config) {
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
		hardd.HardDelete(filePath)
		return
	}

	absoluteSrcPath := util.RelativeToAbsolute(filePath)

	backupFileName := backupFileName(filePath)
	backupLocation := backupDirectory + "/" + backupFileName

	// Soft deleting a path is semi-difficult because we want the most
	// efficiant deletion method. However, the most efficiant methods
	// typically have the most error cases.
	//
	// When soft deleting a file, we take the following steps
	//
	// 1.	Attempt to hard link the old file to the backup location.
	//	This typically fails if the backup destination is on a seperate
	//	partition or the backup location is a non-physical location.
	//	Hard link soft deletion is prefered because it does not invole
	//	copying data and is therefore an O(1) operation (where n is the
	//	size of the file).
	//
	// 2. 	(unsafe) TODO: Move the file to the backup location. This is
	//	only avaliable when running with the --unsafe flag.
	//	This operation is considered unsafe because if the system
	//	crashes half way through the move operation, neither the
	//	original or backup location will have the full picture.
	//	TODO Note: This being an unsafe operation depends upon only
	//	deleting the original files after backups have been made.
	//
	// 3.	If hard link soft deletion does not work, we just copy the data
	//	to the backup location.
	//	This is what most "soft deletion" programs do but is an O(n)
	//	operation, meaning that bigger files take longer to delete.
	//	Copying data can fail if the user does not have read access to
	//	the file.
	//
	// 4.	If none of the above operations work, we give up and do not
	//	delete the file at all.
	//	I have chosen to have a no-operation if soft deletes fail to
	//	prevent accidental data loss.
	//
	// TODO: We should cache what soft deletion operations do not throw an
	// error for each backup location so that we don't have to re-compute
	// what operations the backup location supports.
	// This could possibly be done with a /tmp/ file so that the cache
	// resets after each reboot (to prevent stale cache bugs and so that
	// 2rm can fetch system partition changes, physical disk changes, etc.

	// 1. attempt to hard link the original file to the backup location
	linkErr := os.Link(absoluteSrcPath, backupLocation)
	if linkErr == nil {
		hardd.HardDelete(absoluteSrcPath)
		return
	}

	// 2. TODO

	// 3. attempt to move the file
	moveErr := util.MovePath(absoluteSrcPath, backupLocation)
	if moveErr == nil {
		hardd.HardDelete(absoluteSrcPath)
		return
	}

	// 4. Give up. Do not delete the original file.
	cli.PrintError("Error soft deleting file:" + err.Error())
}

func backupFileName(path string) string {
	result := strings.ReplaceAll(path, ".", "-")
	result = strings.ReplaceAll(result, "/", "_")
	return result + ".bak"
}
