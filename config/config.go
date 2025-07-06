package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/hanle23/shorty/internal/helper"
)

const (
	EnvOverrideConfigDir = "SHORTY_CONFIG"
	ConfigFileName       = "config.yml"
	configFileDir        = "/.config/shorty"
)

// All shortcut name needs to be unique
// Package Name is the name of the program that will be replaced by shortcut name
// Args is all the arguments that the user wants to include after package name
// Description is the string description of the shortcut during list command
type Shortcut struct {
	Shortcut_name string `yaml:"shortcut_name"`
	Package_name  string `yaml:"package_name"`
	Args          string `yaml:"args"`
	Description   string `yaml:"description,omitempty"`
}

// All package name needs to be unique
// Script is the actual script that will be triggered, including all the args
// Description is the string description of the shortcut during list command
type Script struct {
	Package_name string `yaml:"package_name"`
	Script       string `yaml:"script"`
	Description  string `yaml:"description,omitempty"`
}

type Config struct {
	Shortcuts map[string]Shortcut `yaml:"shortcuts"`
	Scripts   map[string]Script   `yaml:"scripts"`
}

var (
	initConfigDir = new(sync.Once)
	configDir     string
	instance      *Config
	once          sync.Once
	mu            sync.RWMutex
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

// TODO: Currently not persistant, need to either make a .env file if this is being set
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

// TODO: Currently not persistant, need to either make a .env file if this is being set
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

func GetConfig() *Config {
	once.Do(func() {
		path := Dir()
		// Load config into instance
		fmt.Println(path)

		instance = &Config{}
	})
	return instance
}

func GetEmptyConfigYAML() ([]byte, error) {
	newConfig := CreateEmptyConfigObject()
	bytes, err := ObjectToYaml(*newConfig)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func ObjectToYaml(object Config) ([]byte, error) {
	bytes, err := yaml.Marshal(object)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func CreateEmptyConfigObject() *Config {
	newObjectConfig := &Config{
		Scripts:   make(map[string]Script),
		Shortcuts: make(map[string]Shortcut),
	}
	return newObjectConfig
}
