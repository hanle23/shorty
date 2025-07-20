package strings_test

import (
	"github.com/hanle23/shorty/internal/strings"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		input string
		sep   string
		want  []string
	}{
		{input: "a,b,c", sep: ",", want: []string{"a", "b", "c"}},
		{input: "something, is, better, than, nothing", sep: ",", want: []string{"something", " is", " better", " than", " nothing"}},
		{input: "/path/to/file", sep: "/", want: []string{"path", "to", "file"}},
		{input: "/path/to/file/", sep: "/", want: []string{"path", "to", "file"}},
	}
	for _, test := range tests {
		got := strings.Split(test.input, test.sep)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Split(%q, %q) = %q, want %q", test.input, test.sep, got, test.want)
		}
	}
}
