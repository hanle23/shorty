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
