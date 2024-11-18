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
			name:           "HardDeleteWithPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "./node_modules/",
			expectedResult: true,
		},
		{
			name:           "HardDeleteWithLongPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "./Documents/client/node_modules/",
			expectedResult: true,
		},
		{
			// while we specified node_modules/ as a hard delete inside the
			// config, test.txt does not match the hard delete config pattern
			// so we expect the nested file to use the default soft-delete
			name:           "HardDeleteWithSuffix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "node_modules/test.txt",
			expectedResult: false,
		},
		{
			name:           "HardDeleteWithLongSuffix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "node_modules/sub_package/dist/test.txt",
			expectedResult: false,
		},
		{
			name:           "HardDeleteWithSuffixAndPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldHardDelete,
			testedPath:     "/home/john_doe/Documents/client/node_modules/test.txt",
			expectedResult: false,
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
			name:           "SoftDeleteWithPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldSoftDelete,
			testedPath:     "./file.bak",
			expectedResult: true,
		},
		{
			name:           "SoftDeleteWithLongPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldSoftDelete,
			testedPath:     "./.local/share/2rm/file.bak",
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
			name:           "ProtectedWithPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.IsProtected,
			testedPath:     "./.ssh/",
			expectedResult: true,
		},
		{
			name:           "ProtectedWithLongPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.IsProtected,
			testedPath:     "./john-doe/crypto/.ssh/",
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
			name:           "OverwriteWithPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldOverwrite,
			testedPath:     "./.ssh/test.pem",
			expectedResult: true,
		},
		{
			name:           "OverwriteWithLongPrefix",
			configPath:     "valid.yml",
			configFunction: models.Config.ShouldOverwrite,
			testedPath:     "./john-doe/crypto/.ssh/test.pem",
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
