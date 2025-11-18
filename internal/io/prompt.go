package io

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/hanle23/shorty/internal/interfaces"
	"github.com/pkg/term"
)

const (
	KeyUp     = byte(65)
	KeyDown   = byte(66)
	KeyEscape = byte(27)
	KeyEnter  = byte(13)
	KeyJ      = byte(106)
	KeyK      = byte(107)
)

var NavigationKeys = map[byte]bool{
	KeyUp:   true,
	KeyDown: true,
}

func ApprovalPrompt(prompt string) bool {
	fmt.Printf("%s, do you want to proceed? (y/N): ", prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return false
	}
	input = strings.TrimSpace(input)
	if input == "y" {
		fmt.Println("Hello,", input)
		return true
	}
	return false
}

func UIntPrompt(prompt string, defaultValue uint32) uint32 {
	fmt.Printf("%s, type the override value or press Enter to select the default (%d), q to exit: ", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if input == "" && err != nil {
		fmt.Println("Error reading input, set setting to default value: ", err)
		return defaultValue
	}
	input = strings.TrimSpace(input)
	if input == "q" {
		os.Exit(1)
	}
	value, err := strconv.Atoi(input)
	if err == nil {
		fmt.Println("Select number is invalid, set setting to default value: ", err)
		return defaultValue
	}
	return uint32(value)
}

// OverridePrompt is a function that prompts the user to override a file or directory
// If the user does not want to override, it will return false
// If the user wants to override, it will return true
func OverrideConfigPrompt(path string) bool {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("Found an existing file or directory (%s), do you want to override this? (y/n)? ", path)
	ans, _ := r.ReadString('\n')
	ans = strings.TrimSpace(ans)
	return ans == "y"
}

// DefaultPathPrompt is a function that prompts the user to use the default path
// If the user does not want to use the default path, it will return false
// If the user wants to use the default path, it will return true
func DefaultPathPrompt(defaultPath string) bool {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("Do you want to use the default path? (%s) (y/n): ", defaultPath)
	ans, _ := r.ReadString('\n')
	ans = strings.TrimSpace(ans)
	return ans == "y"
}

func CustomNewPathPrompt(path string) string {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Please type the full path for the new config file: ")
	newDir, _ := r.ReadString('\n')
	newDir = strings.TrimSpace(newDir)
	return newDir
}

func AddNewShortcutPrompt() *interfaces.Shortcut {
	newShortcut := &interfaces.Shortcut{
		Package_name:  "",
		Shortcut_name: "",
		Args:          []string{},
		Description:   "",
	}
	// TODO: Adding prompt, inputs and assign values into newShortcut
	return newShortcut

}

func AddNewScriptPrompt() *interfaces.Script {
	newScript := &interfaces.Script{
		Package_name: "",
		Script:       "",
		Description:  "",
	}
	// TODO: Adding prompt, inputs and assign values into newScript
	return newScript

}

func RenderModeSelector(index int, redraw bool) {
	modeTypes := [2]string{"Shortcut", "Script"}
	if redraw {
		fmt.Printf("\033[3A\033[0J")
	}
	fmt.Print("Select the type you want to interact with:\n")
	for idx, val := range modeTypes {
		if idx == index {
			fmt.Printf("> %v\n", val)
		} else {
			fmt.Printf("%d %v\n", idx+1, val)
		}
	}
}

// This prompt will return 0 if shortcut or 1 if script
func ModeSelectorPrompt() int {
	index := 0
	RenderModeSelector(index, false)
	for {
		keyCode := getInput()
		switch keyCode {
		case KeyEscape:
			return -1
		case KeyEnter:
			return index
		case KeyUp:
			index = (index - 1) % 2
			if index < 0 {
				index = 1
			}
			RenderModeSelector(index, true)
		case KeyK:
			index = (index - 1) % 2
			if index < 0 {
				index = 1
			}
			RenderModeSelector(index, true)
		case KeyDown:
			index = (index + 1) % 2
			RenderModeSelector(index, true)
		case KeyJ:
			index = (index + 1) % 2
			RenderModeSelector(index, true)
		}
	}
}

func getInput() byte {
	t, _ := term.Open("/dev/tty")
	defer t.Close()

	err := term.RawMode(t)
	if err != nil {
		log.Fatal(err)
	}

	var read int
	readBytes := make([]byte, 3)
	read, err = t.Read(readBytes)
	if err != nil {
		return 0
	}

	defer t.Restore()

	if read == 3 {
		if _, ok := NavigationKeys[readBytes[2]]; ok {
			return readBytes[2]
		}
	} else {
		return readBytes[0]
	}
	return 0
}
