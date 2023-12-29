//go:build unix

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
		{[]string{`/home\\user`, "/bin"}, `/home/user/bin`},
		{[]string{`/usr/`, `local/../local`, `//bin`}, `/usr/local/bin`},
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
		{`/home\\user`, `/home/user`},
		{`/usr/local/../local//bin`, `/usr/local/bin`},
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
		{`/home\\user`, `/home`},
		{`/usr/local/../local//bin`, `/usr/local`},
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
		{`/home\user`, `/home/user`},
		{`/usr\local/../local\/bin`, `/usr/local/../local//bin`},
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
		{`\dev`, `dev`},
		{`/home\user`, `user`},
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
		{`/usr/local`, ""},
	}

	for i, test := range tests {
		dir := paths.VolumeName(test.Input)

		assert.Equal(t, test.Expected, dir, i)
	}
}
