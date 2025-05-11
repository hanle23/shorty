package config

import (
	"testing"
)

func setupConfigDir(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()
	return tmpDir
}
