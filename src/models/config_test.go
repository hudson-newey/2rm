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

func TestShouldHardDelete(t *testing.T) {
	testedConfig := loadConfig("valid.yml")

	expected := true
	realized := testedConfig.ShouldHardDelete("node_modules/")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestShouldNotHardDelete(t *testing.T) {
	testedConfig := loadConfig("valid.yml")

	expected := false
	realized := testedConfig.ShouldHardDelete("src/")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestShouldHardDeleteEmpty(t *testing.T) {
	testedConfig := loadConfig("only_backups.yml")

	expected := false
	realized := testedConfig.ShouldHardDelete("node_modules/")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestShouldSoftDelete(t *testing.T) {
	testedConfig := loadConfig("valid.yml")

	expected := true
	realized := testedConfig.ShouldSoftDelete("file.bak")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestShouldNotSoftDelete(t *testing.T) {
	testedConfig := loadConfig("valid.yml")

	expected := false
	realized := testedConfig.ShouldSoftDelete("file.txt")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestShouldSoftDeleteEmpty(t *testing.T) {
	testedConfig := loadConfig("only_backups.yml")

	expected := false
	realized := testedConfig.ShouldSoftDelete("file.bak")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestHardMatchesAbsolutePath(t *testing.T) {
	testedConfig := loadConfig("abs_path.yml")

	expected := true
	realized := testedConfig.ShouldHardDelete("/tmp/2rm/")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestSoftMatchesAbsolutePath(t *testing.T) {
	testedConfig := loadConfig("abs_path.yml")

	expected := true
	realized := testedConfig.ShouldSoftDelete("/home/john-doe/.local/share/2rm/config.yml")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestIsProtected(t *testing.T) {
	testedConfig := loadConfig("valid.yml")

	expected := true
	realized := testedConfig.IsProtected(".ssh/")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}

func TestNotProtected(t *testing.T) {
	testedConfig := loadConfig("valid.yml")

	expected := false
	realized := testedConfig.IsProtected("src/")

	if expected != realized {
		t.Fatalf("Expected %v but got %v", expected, realized)
	}
}
