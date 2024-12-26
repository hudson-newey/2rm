package util

import (
	"io"
	"os"
)

func CopyFile(src string, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func ListFiles(directory string) []string {
	files, err := os.ReadDir(directory)
	if err != nil {
		return []string{}
	}

	var fileNames []string
	for _, file := range files {
		relativeName := directory + "/" + file.Name()
		fileNames = append(fileNames, relativeName)
	}

	return fileNames
}

func IsDirectoryEmpty(directory string) bool {
	files, err := os.ReadDir(directory)
	if err != nil {
		return true
	}

	return len(files) == 0
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
