package models

import (
	"hudson-newey/2rm/src/util"
	"path/filepath"
)

type Config struct {
	// the backup location for soft-deletes
	Backups string

	// any file paths that match these patterns will be overwritten with
	// zeros when deleted
	Overwrite []string

	// any file paths that match these patterns will be hard-deleted
	Hard []string

	// any file paths that match these patterns will be soft-deleted
	// soft-deletes take precedence over hard-deletes
	// meaning that if a file matches both a hard and soft delete pattern
	// the file will be soft-deleted
	Soft []string

	// any file paths that match these patterns will be protected from deletion
	// protected files cannot be deleted without the --bypass-protected flag
	Protected []string
}

func (config Config) ShouldHardDelete(path string) bool {
	for _, hardDeletePath := range config.Hard {
		matched := matchesPattern(hardDeletePath, path)
		if matched {
			return true
		}
	}

	return false
}

func (config Config) ShouldSoftDelete(path string) bool {
	for _, softDeletePath := range config.Soft {
		matched := matchesPattern(softDeletePath, path)
		if matched {
			return true
		}
	}

	return false
}

func (config Config) ShouldOverwrite(path string) bool {
	for _, overwritePath := range config.Overwrite {
		matched := matchesPattern(overwritePath, path)
		if matched {
			return true
		}
	}

	return false
}

func (config Config) IsProtected(path string) bool {
	return util.InArray(config.Protected, path)
}

// if the user has not specified a backup directory, we will use a default
// directory of /tmp/2rm
// I have chosen this directory because it will be automatically cleaned up
// by the system after a reboot
func (config Config) SoftDeleteDir() string {
	if config.Backups != "" {
		return config.Backups
	}

	return "/tmp/2rm/"
}

func matchesPattern(pattern string, path string) bool {
	matched, _ := filepath.Match(pattern, path)
	return matched
}
