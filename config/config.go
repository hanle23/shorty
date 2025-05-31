package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/hanle23/shorty/internal/helper"
)

const (
	EnvOverrideConfigDir = "SHORTY_CONFIG"
	ConfigFileName       = "config.yml"
	configFileDir        = "/.config/shorty"
)

type Shortcut struct {
	PackageName string   `yaml:"package_name"`
	Args        []string `yaml:"args"`
	Description string   `yaml:"description"`
}

type Script struct {
	Content     string `yaml:"content"`
	Description string `yaml:"description,omitempty"`
}

type Config struct {
	Shortcuts map[string]Shortcut `yaml:"shortcuts"`
	Scripts   map[string]Script   `yaml:"scripts"`
}

var (
	initConfigDir = new(sync.Once)
	configDir     string
)

func Dir() string {
	initConfigDir.Do(func() {
		configDir = os.Getenv(EnvOverrideConfigDir)
		if configDir == "" {
			configDir = filepath.Join(helper.GetHomeDir(), configFileDir)
		}
	})
	return configDir
}

func DefaultPath() string {
	dir := filepath.Join(helper.GetHomeDir(), configFileDir)
	return dir
}

func SetOverrideConfigDir(dir string) error {
	isExist := helper.IsExist(dir)
	if !isExist {
		err := helper.CreateDir(dir, false)
		if err != nil {
			return err
		}
	}
	err := os.Setenv(EnvOverrideConfigDir, dir)
	return err

}

func SetDefaultConfigDir() error {
	dir := DefaultPath()
	isExist := helper.IsExist(dir)
	if !isExist {
		err := helper.CreateDir(dir, false)
		if err != nil {
			return err
		}
	}
	err := os.Setenv(EnvOverrideConfigDir, "")
	return err
}

func InitFlow(isNewConfig bool) error {
	currConfigDir := Dir()
	if !isNewConfig {
		shouldOverride := helper.OverrideConfigPrompt(currConfigDir)
		if !shouldOverride {
			return nil
		}
	}

	defaultPath := DefaultPath()
	shouldUseDefault := helper.DefaultPathPrompt(defaultPath)
	if shouldUseDefault {
		fmt.Println("Initiating config to default path...")
		err := SetDefaultConfigDir()
		if err != nil {
			return err
		}
	} else {
		newDir := helper.CustomNewPathPrompt(defaultPath)
		if newDir == "" {
			return nil
		}
		fmt.Println("Initiating config to overriding path...")
		err := SetOverrideConfigDir(newDir)
		if err != nil {
			return err
		}
	}
	return nil
}
