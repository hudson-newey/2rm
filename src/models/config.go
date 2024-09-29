package models

import "path/filepath"

type Config struct {
	Hard []string
	Soft []string
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

func matchesPattern(pattern string, path string) bool {
	matched, _ := filepath.Match(pattern, path)
	return matched
}
