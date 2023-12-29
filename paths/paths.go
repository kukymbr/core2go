package paths

import (
	"os"
	"path/filepath"
	"strings"
)

// Join is a wrapper for the path.Join function;
// differences in the behavior:
//   - strict first element
//   - fixes separator to os.PathSeparator inside the parts.
func Join(path string, parts ...string) string {
	elem := make([]string, 0, len(parts)+1)
	elem = append(elem, path)

	if len(parts) > 0 {
		elem = append(elem, parts...)
	}

	for i, e := range elem {
		elem[i] = FixSeparators(e)
	}

	return filepath.Join(elem...)
}

// Clean is a wrapper for the filepath.Clean function;
// differences in the behavior:
//   - fixes separator to os.PathSeparator before processing.
func Clean(path string) string {
	path = FixSeparators(path)

	return filepath.Clean(path)
}

// Dir is a wrapper for the filepath.Dir function;
// differences in the behavior:
//   - calls Clean before processing.
func Dir(path string) string {
	path = Clean(path)

	return filepath.Dir(path)
}

// FixSeparators replaces all path separators to the OS-correct.
func FixSeparators(path string) string {
	if path == "" {
		return ""
	}

	sepToReplace := '/'
	if os.PathSeparator == sepToReplace {
		sepToReplace = '\\'
	}

	return strings.ReplaceAll(path, string(sepToReplace), string(os.PathSeparator))
}

// Base is a wrapper for the filepath.Base function;
// differences in the behavior:
//   - fixes separator to os.PathSeparator before processing.
func Base(path string) string {
	path = FixSeparators(path)

	return filepath.Base(path)
}

// VolumeName is a wrapper for a filepath.VolumeName;
// differences in the behavior:
//   - fixes separator to os.PathSeparator before processing.
func VolumeName(path string) string {
	path = FixSeparators(path)

	return filepath.VolumeName(path)
}

// Ext is a wrapper for the filepath.Ext function;
// differences in the behavior:
//   - removes the dot from the result string
//   - fixes separator to os.PathSeparator before processing.
func Ext(path string) string {
	path = FixSeparators(path)

	ext := filepath.Ext(path)

	return strings.TrimLeft(ext, ".")
}

// RemoveExt returns a path without an extension.
func RemoveExt(path string) string {
	ext := Ext(path)
	if ext == "" {
		return path
	}

	return strings.TrimSuffix(path, "."+ext)
}

// Called returns the path or the name of the called executable from the args.
func Called() string {
	called := os.Args[0]
	called = FixSeparators(called)

	return called
}

// Executable returns a path to the running executable.
// It's a wrapper for an os.Executable with the differences:
//   - symlinks are followed
//   - path is cleaned.
//
// If os.Executable fails, a path will be taken from the args.
func Executable() string {
	path, err := os.Executable()
	if err != nil {
		called := Called()

		path, err = filepath.Abs(called)
		if err != nil {
			path = called
		}
	}

	path = FixSeparators(path)

	absPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return Clean(path)
	}

	return absPath
}

// ExecutableDir returns a running executable dir path.
func ExecutableDir() string {
	exe := Executable()

	return Dir(exe)
}

// Abs is a wrapper for the filepath.Abs function;
// differences in the behavior:
//   - fixes separator to os.PathSeparator before processing.
func Abs(path string) (string, error) {
	path = FixSeparators(path)

	return filepath.Abs(path)
}

// EvalSymlinks is a wrapper for the filepath.EvalSymlinks function;
// differences in the behavior:
//   - fixes separator to os.PathSeparator before processing.
func EvalSymlinks(path string) (string, error) {
	path = FixSeparators(path)

	return filepath.EvalSymlinks(path)
}
