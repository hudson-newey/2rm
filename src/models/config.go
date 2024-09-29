package models

import "path/filepath"

type Config struct {
	Backups string
	Hard    []string
	Soft    []string
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
