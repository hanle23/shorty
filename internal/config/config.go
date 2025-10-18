package config

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"os"
	"path/filepath"
	"sync"

	"github.com/hanle23/shorty/internal/fs"
	"github.com/hanle23/shorty/internal/io"
)

const (
	ConfigFileName    = "config.yaml"
	ShortcutsFileName = "shortcuts.yaml"
	DefaultFileDir    = "/.config/shorty"
)

// All shortcut name needs to be unique
// Package Name is the name of the program that will be replaced by shortcut name
// Args is all the arguments that the user wants to include after package name
// Description is the string description of the shortcut during list command
type Shortcut struct {
	Shortcut_name string   `yaml:"shortcut_name"`
	Package_name  string   `yaml:"package_name"`
	Args          []string `yaml:"args"`
	Description   string   `yaml:"description,omitempty"`
}

// All package name needs to be unique
// Script is the actual script that will be triggered, including all the args
// Description is the string description of the shortcut during list command
type Script struct {
	Package_name string `yaml:"package_name"`
	Script       string `yaml:"script"`
	Description  string `yaml:"description,omitempty"`
}

type ShortcutFile struct {
	Shortcuts map[string]Shortcut `yaml:"shortcuts"`
	Scripts   map[string]Script   `yaml:"scripts"`
}

type ConfigFile struct {
	ShortcutPath string `yaml:"shortcut_path"`
}

var (
	initConfigDir    = new(sync.Once)
	configDir        string
	shortcutInstance *ShortcutFile
	configInstance   *ShortcutFile
	once             sync.Once
	mu               sync.RWMutex
)

func GetEmptyShortcutObject() *ShortcutFile {
	newShortcutObject := &ShortcutFile{
		Scripts:   make(map[string]Script),
		Shortcuts: make(map[string]Shortcut),
	}
	return newShortcutObject
}

func GetShortcutPath() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	if config.ShortcutPath != "" {
		return config.ShortcutPath, nil
	}
	return "", nil
}

func GetEmptyConfigObject() (*ConfigFile, error) {
	defaultPath, err := GetDefaultPath()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(defaultPath, ConfigFileName)
	newConfigObject := &ConfigFile{
		ShortcutPath: path,
	}
	return newConfigObject, nil
}

func LoadConfig() (*ConfigFile, error) {
	path, err := GetDefaultPath()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, ConfigFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg ConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func LoadShortcut() error {
	path, err := GetShortcutPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read shortcut: %w", err)
	}

	var shortcut ShortcutFile
	if err := yaml.Unmarshal(data, &shortcut); err != nil {
		return fmt.Errorf("failed to parse shortcut: %w", err)
	}
	shortcutInstance = &shortcut
	return nil
}

func GetShortcut() (*ShortcutFile, error) {
	if shortcutInstance != nil {
		return shortcutInstance, nil
	}
	err := LoadShortcut()
	if err != nil {
		return nil, err
	}
	return shortcutInstance, nil
}

// func GetEmptyShortcutYAML() ([]byte, error) {
//	newShortcut := GetEmptyShortcutObject()
//	bytes, err := yamlutil.ObjectToYaml(*newShortcut)
//	if err != nil {
//		return nil, err
//	}
//	return bytes, nil
//}

// Get default path for shortcut
func GetDefaultPath() (string, error) {
	homeDir, err := fs.GetHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, DefaultFileDir)
	return dir, nil
}

// Get current shortcut path from config if was overrided
func GetShortcutDir() (string, error) {
	initConfigDir.Do(func() {
		// TODO: Need to load config file here and get from config instead of env
		// TODO: Need to also rethink on how to write this function
	})
	return configDir, nil
}

// TODO: Currently not persistant, need to set it into config file
func SetOverrideShortcutDir(dir string) error {
	isExist := fs.IsExist(dir)
	if !isExist {
		err := fs.CreateDir(dir, false)
		if err != nil {
			return err
		}
	}
	// TODO: Need to load or create config object here and append, then write it into config file
	return nil
}

// TODO: Currently not persistant, need to either make a .env file if this is being set
func SetDefaultShortcutDir() error {
	dir, err := GetDefaultPath()
	if err != nil {
		return err
	}
	isExist := fs.IsExist(dir)
	if !isExist {
		err := fs.CreateDir(dir, false)
		if err != nil {
			return err
		}
	}
	// TODO: Need to load or create config object here and append, then write it into config file
	return nil
}

func InitShortcut(isNewShortcut bool) error {
	currShortcutDir, err := GetShortcutDir()
	if err != nil {
		return err
	}
	if !isNewShortcut {
		isExist := fs.IsExist(currShortcutDir)
		if isExist {
			shouldOverride := io.OverrideConfigPrompt(currShortcutDir)
			if !shouldOverride {
				return nil
			}
		}
	}

	defaultPath, err := GetDefaultPath()
	if err != nil {
		return err
	}
	shouldUseDefault := io.DefaultPathPrompt(defaultPath)
	if shouldUseDefault {
		fmt.Println("Initiating config to default path...")
		err := SetDefaultShortcutDir()
		if err != nil {
			return err
		}
	} else {
		newDir := io.CustomNewPathPrompt(defaultPath)
		if newDir == "" {
			return nil
		}
		fmt.Println("Setting new path...")
		err := SetOverrideShortcutDir(newDir)
		if err != nil {
			return err
		}
	}
	return nil
}
