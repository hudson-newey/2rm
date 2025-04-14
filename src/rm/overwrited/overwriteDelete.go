package overwrited

import (
	"hudson-newey/2rm/src/cli"
	"os"
)

func OverwriteFile(filePath string) {
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
