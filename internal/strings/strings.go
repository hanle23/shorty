package strings

import (
	"fmt"
	"github.com/hanle23/shorty/internal/config"
	"math"
	"strings"
)

var (
	PADDING int16 = 5
)

func Split(s string, sep string) []string {
	result := make([]string, 0)
	lastIndex := 0
	for i, item := range s {
		if string(item) != sep {
			continue
		}
		curr := s[lastIndex:i]
		if curr != "" {
			result = append(result, curr)
		}
		lastIndex = i + len(sep)
	}
	if lastIndex < len(s) && s[lastIndex:] != "" {
		result = append(result, s[lastIndex:])
	}
	return result
}

func PrintStringLimit(s string, limit int16) {
	rL := math.Min(float64(limit), float64(len(s)))
	fmt.Print(s[:int(rL)])
	padding := limit - int16(rL) + PADDING
	fmt.Print(strings.Repeat(" ", int(padding)))
}

func PrintShortcut(s config.ShortcutFile, l int16) {
	var headers = []string{"Name", "Package", "Description", "Args"}
	for _, v := range headers {
		PrintStringLimit(v, l)
	}
	fmt.Println("\n")
	for _, v := range s.Shortcuts {
		PrintStringLimit(v.Shortcut_name, l)
		PrintStringLimit(v.Package_name, l)
		PrintStringLimit(v.Description, l)
	}
}
