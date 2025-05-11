package config

import (
	"bufio"
	"fmt"
	"github.com/hanle23/shorty/internal/helper"
	"os"
	"path/filepath"
	"strings"
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

func InitFlow() error {
	currConfigDir := Dir()
	isExist := helper.IsExist(currConfigDir)
	r := bufio.NewReader(os.Stdin)
	if isExist {
		fmt.Printf("Found an existing configuration file (%s), do you want to override this? (y/n)? ", currConfigDir)
		ans, _ := r.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if ans == "n" {
			return nil
		}
	}
	defaultPath := DefaultPath()
	fmt.Printf("Do you want to override the default path? (%s) (y/n): ", defaultPath)
	ans, _ := r.ReadString('\n')
	ans = strings.TrimSpace(ans)
	if ans == "n" {
		fmt.Println("Initiating config to default path...")
		err := SetDefaultConfigDir()
		return err
	}

	fmt.Print("Please type the full path for the new config file: ")
	ans, _ = r.ReadString('\n')
	ans = strings.TrimSpace(ans)
	fmt.Println("Initiating config to overrided path...")
	err := SetOverrideConfigDir(ans)
	return err
}
