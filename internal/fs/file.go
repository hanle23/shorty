package fs

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"syscall"

	"github.com/hanle23/shorty/internal/context"
	"github.com/hanle23/shorty/internal/io"
)

//const (
//	fileModeVal = uint32(0755)
//)

func GetHomeDir() (string, error) {
	c := context.GetContext()
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		if c.Debug {
			fmt.Println("Getting env from XDG_CONFIG_HOME")
		}
		return xdg, nil
	}
	home, _ := os.UserHomeDir()
	if home != "" && (runtime.GOOS == "windows" || runtime.GOOS == "darwin") {
		if c.Debug {
			fmt.Println("Getting env from os.UserHomeDir()")
		}
		return home, nil
	}
	u, err := user.Current()
	if err != nil {
		return u.HomeDir, err
	}
	if c.Debug {
		fmt.Println("Getting env from user.Current()")
	}
	return "", err
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// TODO: Maybe display the previous override file mode and ask if the user wants to reuse this setting
func CreateDir(path string, bypassPrompt bool) error {
	isExist := IsExist(path)
	fileModeVal := uint32(0755)
	if !bypassPrompt {
		allowOverride := io.OverrideConfigPrompt(path)
		if allowOverride && isExist {
			err := os.RemoveAll(path)
			if err != nil {
				return fmt.Errorf("failed to remove existing directory: %w", err)
			}
		}
		prompt := fmt.Sprintf("%s will be created", path)
		fileModeVal = io.UIntPrompt(prompt, fileModeVal)
	}
	fileMode := os.FileMode(fileModeVal)

	// Save current umask
	oldMask := syscall.Umask(0)
	// Restore umask when function returns
	defer syscall.Umask(oldMask)

	err := os.MkdirAll(path, fileMode)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat directory: %w", err)
	}
	if info.Mode().Perm() != fileMode.Perm() {
		return fmt.Errorf("directory permissions mismatch. got: %v, expected: %v", info.Mode().Perm(), fileMode)
	}
	fmt.Printf("Successfully create file %s\n", path)
	return nil
}
