package models_test

import (
	"hudson-newey/2rm/src/config"
	"hudson-newey/2rm/src/models"
	"path/filepath"
	"testing"
)

func loadConfig(path string) models.Config {
	testConfigDir := "../../tests/assets/configs/"
	absolutePath, _ := filepath.Abs(testConfigDir + path)
	return config.ParseConfig(absolutePath)
}

func assertConfig(
	t *testing.T,
	configPath string,
	configFunction func(models.Config, string) bool,
	testedPath string,
	expectedResult bool,
) {
	expectedConfig := expectedResult

	testedConfig := loadConfig(configPath)
	realizedConfig := configFunction(testedConfig, testedPath)

	if expectedConfig != realizedConfig {
		t.Fatalf("Expected %v but got %v", expectedConfig, realizedConfig)
	}
}

func runTests(t *testing.T, tests []Test) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assertConfig(t, test.configPath, test.configFunction, test.testedPath, test.expectedResult)
		})
	}
}

type Test struct {
	name           string
	configPath     string
	configFunction func(models.Config, string) bool
	testedPath     string
	expectedResult bool
}

func TestConfig(t *testing.T) {
	tests := []Test{
		// hard deletes
		{
			name:           "HardDelete",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "node_modules/",
			expectedResult: true,
		},
		{
			name:           "NotHardDelete",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "src/",
			expectedResult: false,
		},
		{
			name:           "HardDeleteEmpty",
			configPath:     "only_backups.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "node_modules/",
			expectedResult: false,
		},
		{
			name:           "HardMatchesAbsolutePath",
			configPath:     "abs_path.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "/tmp/2rm/",
			expectedResult: true,
		},

		// soft deletes
		{
			name:           "SoftDelete",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldSoftDelete,
			testedPath:     "file.bak",
			expectedResult: true,
		},
		{
			name:           "NotSoftDelete",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldSoftDelete,
			testedPath:     "file.txt",
			expectedResult: false,
		},
		{
			name:           "SoftDeleteEmpty",
			configPath:     "only_backups.yml",
			configFunction: models.Config.ShouldSoftDelete,
			testedPath:     "file.bak",
			expectedResult: false,
		},
		{
			name:           "SoftMatchesAbsolutePath",
			configPath:     "abs_path.yml",
			configFunction: models.Config.ShouldSoftDelete,
			testedPath:     "/home/john-doe/.local/share/2rm/config.yml",
			expectedResult: true,
		},

		// protected files
		{
			name:           "IsProtected",
			configPath:     "valid.yml",
			configFunction: models.Config.IsProtected,
			testedPath:     ".ssh/",
			expectedResult: true,
		},
		{
			name:           "NotProtected",
			configPath:     "valid.yml",
			configFunction: models.Config.IsProtected,
			testedPath:     "src/",
			expectedResult: false,
		},
		{
			name:           "ProtectedEmpty",
			configPath:     "only_backups.yml",
			configFunction: models.Config.IsProtected,
			testedPath:     ".ssh/",
			expectedResult: false,
		},
		{
			name:           "ProtectedMatchesAbsolutePath",
			configPath:     "abs_path.yml",
			configFunction: models.Config.IsProtected,
			testedPath:     "/home/john-doe/.ssh/id_rsa",
			expectedResult: true,
		},

		// overwrite
		{
			name:           "Overwrite",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldOverwrite,
			testedPath:     ".ssh/test.pem",
			expectedResult: true,
		},
		{
			name:           "DontOverwrite",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldOverwrite,
			testedPath:     "non-existent.txt",
			expectedResult: false,
		},
		{
			name:           "OverwriteEmpty",
			configPath:     "only_backups.yml",
			configFunction: models.Config.ShouldOverwrite,
			testedPath:     ".ssh/test.pem",
			expectedResult: false,
		},
		{
			name:           "OverwriteAbsolutePath",
			configPath:     "abs_path.yml",
			configFunction: models.Config.ShouldOverwrite,
			testedPath:     "/home/john-doe/.ssh/key.pem",
			expectedResult: true,
		},
	}

	runTests(t, tests)
}
