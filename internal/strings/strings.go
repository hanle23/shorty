package strings

import (
	"fmt"
	"math"
	"strings"

	"github.com/hanle23/shorty/internal/types"
)

var (
	PADDING int16 = 5
	DOTSLEN int16 = 3
)

func FormatStringLimit(s string, limit int16) string {
	var result strings.Builder
	rL := math.Min(float64(limit), float64(len(s)))
	result.WriteString(s[:int(rL)])
	padding := limit - int16(rL) + PADDING
	if int(rL) < len(s) {
		padding -= DOTSLEN
		result.WriteString(strings.Repeat(".", int(DOTSLEN)))
	}
	result.WriteString(strings.Repeat(" ", int(padding)))
	return result.String()
}

func PrintRunnable(s types.RunnableFile, l int16) {
	var headers = []string{"Name", "Package", "Description", "Args"}
	fmt.Println("SHORTCUTS")
	for _, v := range headers {
		fmt.Print(FormatStringLimit(v, l))
	}
	fmt.Println()
	for _, v := range s.Shortcuts {
		fmt.Print(FormatStringLimit(v.Shortcut_name, l))
		fmt.Print(FormatStringLimit(v.Package_name, l))
		fmt.Print(FormatStringLimit(v.Description, l))
		argsString := ""
		for _, arg := range v.Args {
			argsString += arg
			argsString += ", "
		}
		fmt.Print(FormatStringLimit(argsString, l))
		fmt.Println()
	}

	fmt.Println()
	var scriptHeaders = []string{"Package", "Script", "Description"}
	fmt.Println("SCRIPTS")
	for _, v := range scriptHeaders {
		fmt.Print(FormatStringLimit(v, l))
	}
	fmt.Println()
	for _, v := range s.Scripts {
		fmt.Print(FormatStringLimit(v.Package_name, l))
		fmt.Print(FormatStringLimit(v.Script, l))
		fmt.Print(FormatStringLimit(v.Description, l))
		fmt.Println()
	}
}
