package models

import "path/filepath"

type Config struct {
	Hard []string
}

func (config Config) ShouldHardDelete(path string) bool {
	for _, hardDeletePath := range config.Hard {
		matched, _ := filepath.Match(hardDeletePath, path)
		if matched {
			return true
		}
	}

	return false
}
