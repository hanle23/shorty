package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"

	"github.com/hanle23/shorty/internal/fs"
	"github.com/hanle23/shorty/internal/io"
	"github.com/hanle23/shorty/internal/types"
	"github.com/hanle23/shorty/internal/yamlutil"
)

const (
	ConfigFileName   = "config.yaml"
	RunnableFileName = "runnables.yaml"
	DefaultFileDir   = "/.config/shorty"
)

var (
	RunnableInstance *types.RunnableFile
	configInstance   *types.ConfigFile
)

func GetEmptyRunnableObject() *types.RunnableFile {
	newRunnableObject := &types.RunnableFile{
		Scripts:   make(map[string]types.Script),
		Shortcuts: make(map[string]types.Shortcut),
	}
	return newRunnableObject
}

func GetEmptyConfigObject(runnablePath string) (*types.ConfigFile, error) {
	newConfigObject := &types.ConfigFile{
		RunnablePath: runnablePath,
	}
	return newConfigObject, nil
}

// Get the path to the runnable file, return empty string if not found from config
func GetRunnablePath() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}
	if config.RunnablePath != "" {
		return config.RunnablePath, nil
	}
	return "", nil
}

// Grabbing config file and load it into configInstance, also return the file object
func LoadConfig() (*types.ConfigFile, error) {
	if configInstance != nil {
		return configInstance, nil
	}
	path, err := GetDefaultFolderPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get default folder path: %w", err)
	}
	path = filepath.Join(path, ConfigFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg types.ConfigFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	configInstance = &cfg
	return configInstance, nil
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

	var runnable types.RunnableFile
	if err := yaml.Unmarshal(data, &runnable); err != nil {
		return fmt.Errorf("failed to parse Runnable: %w", err)
	}
	RunnableInstance = &runnable
	return nil
}

func GetRunnable() (*types.RunnableFile, error) {
	if RunnableInstance != nil {
		return RunnableInstance, nil
	}
	err := LoadRunnable()
	if err != nil {
		return nil, fmt.Errorf("failed to load runnable: %w", err)
	}
	return RunnableInstance, nil
}

// Get default path for main config folder
func GetDefaultFolderPath() (string, error) {
	homeDir, err := fs.GetHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, DefaultFileDir)
	return dir, nil
}

func AddShortcut(newShortcut *types.Shortcut) error {
	currRunnable, err := GetRunnable()
	if err != nil {
		return err
	}
	currShortcut := currRunnable.Shortcuts
	_, exist := currShortcut[newShortcut.Shortcut_name]
	if exist {
		proceed := io.YesNoPrompt("This shortcut is already exist, do you wish to overwrite it? (Y/n)")
		if !proceed {
			return nil
		}
	}
	newRunnable := *currRunnable
	newRunnable.Shortcuts[newShortcut.Shortcut_name] = *newShortcut
	RunnableInstance = &newRunnable
	runnablePath, err := GetRunnablePath()
	if err != nil {
		return fmt.Errorf("failed to get runnable path: %w", err)
	}
	runnableYaml, err := yamlutil.ObjectToYaml(newRunnable)
	if err != nil {
		return fmt.Errorf("failed to convert runnable to yaml: %w", err)
	}
	err = SaveYamlFile(runnablePath, runnableYaml)
	if err != nil {
		return fmt.Errorf("failed to save runnable: %w", err)
	}
	return nil
}

func AddScript(newScript *types.Script) error {
	currRunnable, err := GetRunnable()
	if err != nil {
		return err
	}
	currScript := currRunnable.Scripts
	_, exist := currScript[newScript.Package_name]
	if exist {
		proceed := io.YesNoPrompt("This script is already exist, do you wish to overwrite it? (Y/n)")
		if !proceed {
			return nil
		}
	}
	newRunnable := *currRunnable
	newRunnable.Scripts[newScript.Package_name] = *newScript
	RunnableInstance = &newRunnable
	runnablePath, err := GetRunnablePath()
	if err != nil {
		return fmt.Errorf("failed to get runnable path: %w", err)
	}
	runnableYaml, err := yamlutil.ObjectToYaml(newRunnable)
	if err != nil {
		return fmt.Errorf("failed to convert runnable to yaml: %w", err)
	}
	err = SaveYamlFile(runnablePath, runnableYaml)
	if err != nil {
		return fmt.Errorf("failed to save runnable: %w", err)
	}
	return nil
}

func SetRunnableDir(dir string) error {
	defaultFolderPath, err := GetDefaultFolderPath()
	if err != nil {
		return fmt.Errorf("failed to get default folder path: %w", err)
	}
	if !fs.IsExist(dir) {
		fmt.Println("Directory does not exist, creating...", dir)
		err := fs.CreateDir(dir, false)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
		fmt.Println("Directory created successfully...")
	}
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	config.RunnablePath = filepath.Join(dir, RunnableFileName)
	configYaml, err := yamlutil.ObjectToYaml(*config)
	if err != nil {
		return fmt.Errorf("failed to convert config to yaml: %w", err)
	}
	configFullPath := filepath.Join(defaultFolderPath, ConfigFileName)
	err = SaveYamlFile(configFullPath, configYaml)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}
	return nil
}

func InitFlow() error {
	defaultFolderPath, err := GetDefaultFolderPath()
	if err != nil {
		return fmt.Errorf("failed to get default folder path: %w", err)
	}
	configPath := filepath.Join(defaultFolderPath, ConfigFileName)

	if fs.IsExist(configPath) {
		shouldReset := io.YesNoPrompt("Found existing configuration. Reset config and runnables to original state? (y/n)")
		if !shouldReset {
			fmt.Println("Init cancelled.")
			return nil
		}
	}

	if err := initConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}
	fmt.Printf("Config initialized: %s\n", configPath)

	runnablePath := filepath.Join(defaultFolderPath, RunnableFileName)
	if err := initRunnable(); err != nil {
		return fmt.Errorf("failed to initialize runnables: %w", err)
	}
	fmt.Printf("Runnables initialized: %s\n", runnablePath)

	fmt.Println("Initialization complete.")
	return nil
}

func initConfig() error {
	defaultFolderPath, err := GetDefaultFolderPath()
	if err != nil {
		return fmt.Errorf("failed to get default folder path: %w", err)
	}
	if err := os.MkdirAll(defaultFolderPath, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	defaultRunnablePath := filepath.Join(defaultFolderPath, RunnableFileName)
	emptyConfig, err := GetEmptyConfigObject(defaultRunnablePath)
	if err != nil {
		return fmt.Errorf("failed to create default config: %w", err)
	}
	configYaml, err := yamlutil.ObjectToYaml(*emptyConfig)
	if err != nil {
		return fmt.Errorf("failed to serialize config: %w", err)
	}
	configFullPath := filepath.Join(defaultFolderPath, ConfigFileName)
	if err := SaveYamlFile(configFullPath, configYaml); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	configInstance = nil
	return nil
}

func initRunnable() error {
	runnablePath, err := GetRunnablePath()
	if err != nil {
		return fmt.Errorf("failed to get runnable path: %w", err)
	}
	if runnablePath == "" {
		defaultFolderPath, err := GetDefaultFolderPath()
		if err != nil {
			return fmt.Errorf("failed to get default folder path: %w", err)
		}
		runnablePath = filepath.Join(defaultFolderPath, RunnableFileName)
	}
	emptyRunnable := GetEmptyRunnableObject()
	runnableYaml, err := yamlutil.ObjectToYaml(*emptyRunnable)
	if err != nil {
		return fmt.Errorf("failed to serialize runnables: %w", err)
	}
	if err := SaveYamlFile(runnablePath, runnableYaml); err != nil {
		return fmt.Errorf("failed to write runnables file: %w", err)
	}
	RunnableInstance = nil
	return nil
}

func SaveYamlFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}
