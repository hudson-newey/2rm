package config_test

import (
	"hudson-newey/2rm/src/config"
	"hudson-newey/2rm/src/models"
	"path/filepath"
	"reflect"
	"testing"
)

func assertConfig(t *testing.T, configPath string, expectedConfig models.Config) {
	testConfigDir := "../../tests/assets/configs/"
	absolutePath, err := filepath.Abs(testConfigDir + configPath)
	if err != nil {
		t.Errorf("Failed to get absolute path")
	}

	realizedConfig := config.ParseConfig(absolutePath)

	if !reflect.DeepEqual(expectedConfig, realizedConfig) {
		t.Errorf("Expected %v but got %v", expectedConfig, realizedConfig)
	}
}

func TestParsingConfig(t *testing.T) {
	expectedConfig := models.Config{
		Backups: "/tmp/2rm/",
		Overwrite: []string{
			".ssh/*",
		},
		Hard: []string{
			"node_modules/",
			"target/",
			".angular/",
			".next/",
			"*.partial",
		},
		Soft: []string{
			"*.bak",
		},
		Protected: []string{
			".ssh/",
		},
		Interactive: 10,
	}

	assertConfig(t, "valid.yml", expectedConfig)
}

// this test asserts that we can parse a partial config
// we do not have to check every combination of partial configs
// because that would result in a massive explosion of tests
// that would not provide much value
func TestOnlyHardConfig(t *testing.T) {
	expectedConfig := models.Config{
		Hard: []string{
			"node_modules/",
			"target/",
			".angular/",
			".next/",
		},
	}

	assertConfig(t, "only_hard.yml", expectedConfig)
}

func TestOnlyBackups(t *testing.T) {
	expectedConfig := models.Config{
		Backups: "/tmp/2rm/",
	}
	assertConfig(t, "only_backups.yml", expectedConfig)
}

func TestParsingEmptyConfig(t *testing.T) {
	expectedConfig := models.Config{}
	assertConfig(t, "empty.yml", expectedConfig)
}

func TestParsingConfigWithMissingFile(t *testing.T) {
	expectedConfig := models.Config{}
	assertConfig(t, "missing.yml", expectedConfig)
}
