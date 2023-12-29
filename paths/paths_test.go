package paths_test

import (
	"testing"

	"github.com/kukymbr/core2go/paths"
	"github.com/stretchr/testify/assert"
)

func TestCalled(t *testing.T) {
	called := paths.Called()

	assert.NotEmpty(t, called)
}

func TestExecutable(t *testing.T) {
	executable := paths.Executable()

	assert.NotEmpty(t, executable)
}

func TestExecutableDir(t *testing.T) {
	executable := paths.Executable()
	dir := paths.ExecutableDir()

	assert.Contains(t, executable, dir)
}

func TestAbs(t *testing.T) {
	dir := paths.ExecutableDir()
	abs, err := paths.Abs(dir)

	assert.NoError(t, err)
	assert.NotEmpty(t, abs)
}

func TestEvalSymlinks(t *testing.T) {
	dir := paths.ExecutableDir()
	abs, err := paths.EvalSymlinks(dir)

	assert.NoError(t, err)
	assert.NotEmpty(t, abs)
}

func TestExt(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", ""},
		{`/test/test`, ""},
		{`/test.pkg/test`, ""},
		{`/test/test.csv`, "csv"},
		{`/test/test.csv.test`, "test"},
	}

	for i, test := range tests {
		ext := paths.Ext(test.Input)

		assert.Equal(t, test.Expected, ext, i)
	}
}

func TestRemoveExt(t *testing.T) {
	tests := []struct {
		Input    string
		Expected string
	}{
		{"", ""},
		{`/test/test`, "/test/test"},
		{`/test.pkg/test`, "/test.pkg/test"},
		{`/test/test.csv`, "/test/test"},
		{`/test/test.csv.test`, "/test/test.csv"},
	}

	for i, test := range tests {
		ext := paths.RemoveExt(test.Input)

		assert.Equal(t, test.Expected, ext, i)
	}
}
