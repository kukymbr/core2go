//go:build windows

package paths_test

import (
	"testing"

	"github.com/kukymbr/core2go/paths"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T) {
	tests := []struct {
		Parts    []string
		Expected string
	}{
		{[]string{""}, ""},
		{[]string{`//fs\\test`, "/test1"}, `\\fs\test\test1`},
		{[]string{`C://test//`, `test1/test2`, `//`}, `C:\\test\test1\test2`},
	}

	for i, test := range tests {
		path := paths.Join(test.Parts[0], test.Parts[1:]...)

		assert.Equal(t, test.Expected, path, i)
	}
}

func TestClean(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", "."},
		{`//fs\\test`, `\\fs\test`},
		{`C://test//`, `C:\\test`},
	}

	for i, test := range tests {
		path := paths.Clean(test.Input)

		assert.Equal(t, test.Expected, path, i)
	}
}

func TestDir(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", "."},
		{`//fs\\test`, `\\fs`},
		{`C://test//../test`, `C:\\test`},
	}

	for i, test := range tests {
		dir := paths.Dir(test.Input)

		assert.Equal(t, test.Expected, dir, i)
	}
}

func TestFixSeparators(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", ""},
		{`//fs\\test`, `\\fs\test`},
		{`C://test//`, `C:\\test`},
	}

	for i, test := range tests {
		path := paths.FixSeparators(test.Input)

		assert.Equal(t, test.Expected, path, i)
	}
}

func TestBase(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", "."},
		{`//fs\\test`, `test`},
		{`C://test//`, `test`},
	}

	for i, test := range tests {
		dir := paths.Base(test.Input)

		assert.Equal(t, test.Expected, dir, i)
	}
}

func TestVolumeName(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", ""},
		{`//fs\\test`, `\\fs`},
		{`c://test//`, `C:`},
	}

	for i, test := range tests {
		dir := paths.VolumeName(test.Input)

		assert.Equal(t, test.Expected, dir, i)
	}
}
