package models

import (
	"strings"
)

type Config struct {
	Hard []string
}

func (config Config) ShouldHardDelete(path string) bool {
	for _, hardDeletePath := range config.Hard {
		// if the config hard delte path is an absolute path, we only want to
		// hard delete the path if the full path matches
		isAbsolutePath := strings.HasPrefix(hardDeletePath, "/")
		if isAbsolutePath && path == hardDeletePath {
			return true
		}

		lastConfigPathLocation := lastPathLocation(hardDeletePath)
		lastPathLocation := lastPathLocation(path)
		if lastConfigPathLocation == lastPathLocation {
			return true
		}
	}

	return false
}

func lastPathLocation(path string) string {
	splitPath := strings.Split(path, "/")
	lastDirectory := splitPath[len(splitPath)-1]
	return lastDirectory
}
