package config

import (
	"github.com/hanle23/shorty/internal/helper"
	"os"
	"path/filepath"
	"sync"
)

const (
	EnvOverrideConfigDir = "SHORTY_CONFIG"
	ConfigFileName       = "config.yml"
	configFileDir        = "/.config/shorty"
)

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
		err := helper.CreateDir(dir)
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
		err := helper.CreateDir(dir)
		if err != nil {
			return err
		}
	}
	err := os.Setenv(EnvOverrideConfigDir, "")
	return err
}
