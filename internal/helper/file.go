package helper

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
)

func GetHomeDir() string {
	home, _ := os.UserHomeDir()
	if home == "" && runtime.GOOS != "windows" {
		if u, err := user.Current(); err == nil {
			return u.HomeDir
		}
	}
	return home
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CreateDir(path string) error {
	isExist := IsExist(path)
	allowOverride := OverrideConfigPrompt(path)
	if allowOverride && isExist {
		err := os.RemoveAll(path)
		if err != nil {
			return fmt.Errorf("failed to remove existing directory: %w", err)
		}
	}
	overrideVal := UIntPrompt("%s will be created", 666)
	fileMode := os.FileMode(overrideVal)
	err := os.MkdirAll(path, fileMode)
	if err != nil {
		return err
	}
	fmt.Printf("Successfully create file %s\n", path)
	return nil
}
