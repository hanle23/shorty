package config

import (
	"fmt"
	"log"
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
	RunnableFileName = "Runnable.yaml"
	DefaultFileDir   = "/.config/shorty"
)

var (
	// initConfigDir    = new(sync.Once)
	RunnableInstance *types.RunnableFile
	configInstance   *types.ConfigFile
	// once             sync.Once
	// mu               sync.RWMutex
)

func GetEmptyRunnableObject() *types.RunnableFile {
	newRunnableObject := &types.RunnableFile{
		Scripts:   make(map[string]types.Script),
		Shortcuts: make(map[string]types.Shortcut),
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

func GetEmptyConfigObject() (*types.ConfigFile, error) {
	defaultPath, err := GetDefaultPath()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(defaultPath, ConfigFileName)
	newConfigObject := &types.ConfigFile{
		RunnablePath: path,
	}
	return newConfigObject, nil
}

// Grabbing config file and load it into configInstance, also return the file object
func LoadConfig() (*types.ConfigFile, error) {
	path, err := GetDefaultPath()
	if err != nil {
		return nil, err
	}
	path = filepath.Join(path, ConfigFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg types.ConfigFile
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
		return nil, err
	}
	return RunnableInstance, nil
}

func GetScript() (*types.ConfigFile, error) {
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
	err = SaveRunnableInstance(&newRunnable)
	if err != nil {
		return err
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
	err = SaveRunnableInstance(&newRunnable)
	if err != nil {
		return err
	}
	return nil
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

func SaveRunnableInstance(currRunnable *types.RunnableFile) error {
	runnablePath, err := GetRunnablePath()
	if err != nil {
		return err
	}

	yamlRunnable, err := yamlutil.ObjectToYaml(*currRunnable)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(runnablePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(yamlRunnable); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}
