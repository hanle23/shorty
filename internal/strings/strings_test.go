package strings_test

import (
	"testing"

	"github.com/hanle23/shorty/internal/strings"
)

func TestPrintStringLimit(t *testing.T) {
	tests := []struct {
		input string
		limit int16
		want  string
	}{
		{input: "Hello, world!", limit: 10, want: "Hello, wor...  "},
		{input: "Hello, world!", limit: 20, want: "Hello, world!            "},
		{input: "Testing description through CLI call", limit: 15, want: "Testing descrip...  "},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got := strings.FormatStringLimit(test.input, test.limit)
			if got != test.want {
				t.Errorf("FormatStringLimit(%q, %d) = %q, want %q", test.input, test.limit, got, test.want)
			}
		})
	}
}
