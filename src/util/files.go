package util

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
)

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
	_, err := os.Lstat(path)
	return err == nil
}

// I copy the file here instead of moving the file because sometimes linux
// will start throwing errors when moving a file across different partitions
// fix in: https://github.com/hudson-newey/2rm/issues/17
func MovePath(src string, dst string) error {
	if isLink(src) {
		return moveLink(src, dst)
	}

	return moveFile(src, dst)
}

func RelativeToAbsolute(path string) string {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		os.Exit(2)
	}

	return absolutePath
}

// returns a boolean representing if the path is a symbolic or hard link
func isLink(path string) bool {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return !fileInfo.Mode().IsRegular()
}

func moveFile(src string, dst string) error {
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

func moveLink(src string, dst string) error {
	linkTarget, err := os.Readlink(src)
	if err != nil {
		return err
	}

	return os.Symlink(linkTarget, dst)
}

func GetFileDevice(path string) (uint64, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	sys := fileInfo.Sys()
	if sys == nil {
		return 0, nil // Not supported on this OS
	}

	sysStat, ok := sys.(*syscall.Stat_t)
	if !ok {
		return 0, nil // Not a syscall.Stat_t
	}

	return sysStat.Dev, nil
}
