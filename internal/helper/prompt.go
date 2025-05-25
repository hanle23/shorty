package helper

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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
	if err != nil {
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
	isExist := IsExist(path)
	r := bufio.NewReader(os.Stdin)
	if !isExist {
		return false
	}
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
	fmt.Printf("Do you want to usethe default path? (%s) (y/n): ", defaultPath)
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
