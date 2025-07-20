package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/goccy/go-yaml"
	"github.com/hanle23/shorty/internal/fs"
	"github.com/hanle23/shorty/internal/io"
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
			homeDir, err := fs.GetConfigHomeDir()
			if err == nil {
				configDir = filepath.Join(homeDir, configFileDir)
			}
		}
	})
	return configDir
}

func DefaultPath() string {
	homeDir, err := fs.GetConfigHomeDir()
	if err != nil {
		return ""
	}
	dir := filepath.Join(homeDir, configFileDir)
	return dir
}

// TODO: Currently not persistant, need to either make a .env file if this is being set
func SetOverrideConfigDir(dir string) error {
	isExist := fs.IsExist(dir)
	if !isExist {
		err := fs.CreateDir(dir, false)
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
	isExist := fs.IsExist(dir)
	if !isExist {
		err := fs.CreateDir(dir, false)
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
		isExist := fs.IsExist(currConfigDir)
		if isExist {
			shouldOverride := io.OverrideConfigPrompt(currConfigDir)
			if !shouldOverride {
				return nil
			}

		}
	}

	defaultPath := DefaultPath()
	shouldUseDefault := io.DefaultPathPrompt(defaultPath)
	if shouldUseDefault {
		fmt.Println("Initiating config to default path...")
		err := SetDefaultConfigDir()
		if err != nil {
			return err
		}
	} else {
		newDir := io.CustomNewPathPrompt(defaultPath)
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
		//TODO: Load config into instance
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
