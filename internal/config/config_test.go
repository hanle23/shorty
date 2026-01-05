package config_test

import (
	"testing"

	"github.com/hanle23/shorty/internal/config"
)

func TestLoadConfig(t *testing.T) {
	config, err := config.LoadConfig()
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}
	if config == nil {
		t.Error("Config should not be nil")
	}
}
