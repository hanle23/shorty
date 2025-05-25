package helper

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	home := GetHomeDir()
	if home == "" {
		t.Error("GetHomeDir returned empty string")
	}

	// Test that the directory exists
	if !IsExist(home) {
		t.Errorf("Home directory %s does not exist", home)
	}
}

func TestIsExist(t *testing.T) {
	// Test with existing directory
	tempDir := t.TempDir()
	if !IsExist(tempDir) {
		t.Errorf("IsExist returned false for existing directory: %s", tempDir)
	}

	// Test with non-existing directory
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	if IsExist(nonExistentDir) {
		t.Errorf("IsExist returned true for non-existing directory: %s", nonExistentDir)
	}
}

func TestCreateDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	testDir := filepath.Join(tempDir, "testdir")

	// Test creating a new directory
	err := CreateDir(testDir)
	if err != nil {
		t.Errorf("CreateDir() failed to create new directory: %v", err)
	}

	// Verify directory was created
	if !IsExist(testDir) {
		t.Errorf("CreateDir() did not create directory: %s", testDir)
	}

	// Verify directory permissions
	info, err := os.Stat(testDir)
	if err != nil {
		t.Errorf("Failed to stat created directory: %v", err)
		return
	}
	if !info.IsDir() {
		t.Errorf("Created path is not a directory: %s", testDir)
	}

	// Test creating nested directory
	nestedDir := filepath.Join(testDir, "nested", "dir")
	err = CreateDir(nestedDir)
	if err != nil {
		t.Errorf("CreateDir() failed to create nested directory: %v", err)
	}

	// Verify nested directory was created
	if !IsExist(nestedDir) {
		t.Errorf("CreateDir() did not create nested directory: %s", nestedDir)
	}
}

func TestCreateDirWithInvalidPath(t *testing.T) {
	// Test creating directory with invalid path
	invalidPath := filepath.Join(t.TempDir(), "invalid", "path", "with", "special", "chars", "\\/*?")
	err := CreateDir(invalidPath)
	if err == nil {
		t.Error("CreateDir() should fail with invalid path")
	}
}

func TestCreateDirWithExistingDir(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Test creating the same directory again
	err := CreateDir(tempDir)
	if err != nil {
		t.Errorf("CreateDir() failed to handle existing directory: %v", err)
	}

	// Verify the directory still exists
	if !IsExist(tempDir) {
		t.Error("CreateDir() removed existing directory")
	}
}
