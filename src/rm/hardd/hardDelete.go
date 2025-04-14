package hardd

import (
	"hudson-newey/2rm/src/cli"
	"os"
)

func HardDelete(filePath string) {
	err := os.RemoveAll(filePath)
	if err != nil {
		cli.PrintError("Error deleting file:" + err.Error())
	}
}
