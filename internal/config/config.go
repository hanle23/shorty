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
	ConfigFileName   = "config.yaml"
	RunnableFileName = "Runnable.yaml"
	DefaultFileDir   = "/.config/shorty"
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

type RunnableFile struct {
	Shortcuts map[string]Shortcut `yaml:"shortcuts"`
	Scripts   map[string]Script   `yaml:"scripts"`
}

type ConfigFile struct {
	RunnablePath string `yaml:"runnable_path"`
}

var (
	initConfigDir    = new(sync.Once)
	RunnableInstance *RunnableFile
	configInstance   *ConfigFile
	once             sync.Once
	mu               sync.RWMutex
)

func GetEmptyRunnableObject() *RunnableFile {
	newRunnableObject := &RunnableFile{
		Scripts:   make(map[string]Script),
		Shortcuts: make(map[string]Shortcut),
	}
	return newRunnableObject
}

func GetRunnablePath() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}
	if config.RunnablePath != "" {
		return config.RunnablePath, nil
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
		RunnablePath: path,
	}
	return newConfigObject, nil
}

// Grabbing config file and load it into configInstance, also return the file object
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
	configInstance = &cfg
	return &cfg, nil
}

func LoadRunnable() error {
	path, err := GetRunnablePath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read runnable: %w", err)
	}

	var runnable RunnableFile
	if err := yaml.Unmarshal(data, &runnable); err != nil {
		return fmt.Errorf("failed to parse Runnable: %w", err)
	}
	RunnableInstance = &runnable
	return nil
}

func GetRunnable() (*RunnableFile, error) {
	if RunnableInstance != nil {
		return RunnableInstance, nil
	}
	err := LoadRunnable()
	if err != nil {
		return nil, err
	}
	return RunnableInstance, nil
}

func GetScript() (*ConfigFile, error) {
	return nil, nil
}

// func GetEmptyShortcutYAML() ([]byte, error) {
//	newShortcut := GetEmptyShortcutObject()
//	bytes, err := yamlutil.ObjectToYaml(*newShortcut)
//	if err != nil {
//		return nil, err
//	}
//	return bytes, nil
//}

// Get default path for main config folder
func GetDefaultPath() (string, error) {
	homeDir, err := fs.GetHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, DefaultFileDir)
	return dir, nil
}

// TODO: Override the config shortcutDir with this new dir
func SetOverrideRunnableDir(dir string) error {
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

func SetDefaultRunnableDir() error {
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

// TODO: This function should retrieve paths from prompt and handle the creation at the same time
func InitRunnable(isNewRunnable bool) error {
	currRunnableDir, err := GetRunnablePath()
	if err != nil {
		return err
	}
	isExist := fs.IsExist(currRunnableDir)
	if !isNewRunnable && isExist {
		shouldOverride := io.OverrideConfigPrompt(currRunnableDir)
		if !shouldOverride {
			return nil
		}
	}
	defaultPath, err := GetDefaultPath()
	if err != nil {
		return err
	}
	shouldUseDefault := io.DefaultPathPrompt(defaultPath)
	if shouldUseDefault {
		fmt.Println("Initiating config to default path...")
		err := SetDefaultRunnableDir()
		if err != nil {
			return err
		}
	} else {
		newDir := io.CustomNewPathPrompt(defaultPath)
		if newDir == "" {
			return nil
		}
		fmt.Println("Setting new path...")
		err := SetOverrideRunnableDir(newDir)
		if err != nil {
			return err
		}
	}
	return nil
}
