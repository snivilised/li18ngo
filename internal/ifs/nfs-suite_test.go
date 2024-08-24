package ifs_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nfs Suite")
}

func TrimRoot(root string) string {
	// omit leading '/', because test-fs stupidly doesn't like it,
	// so we have to jump through hoops
	if strings.HasPrefix(root, string(filepath.Separator)) {
		return root[1:]
	}

	pattern := `^[a-zA-Z]:[\\/]*`
	re, _ := regexp.Compile(pattern)

	return re.ReplaceAllString(root, "")
}

const (
	perm = 0o766
)

type (
	ensureTE struct {
		given     string
		should    string
		relative  string
		expected  string
		directory bool
	}
	mkDirAllMapFS struct {
		mapFS fstest.MapFS
	}
)

func (f *mkDirAllMapFS) FileExists(path string) bool {
	fi, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	return true
}

func (f *mkDirAllMapFS) DirectoryExists(path string) bool {
	if strings.HasPrefix(path, string(filepath.Separator)) {
		path = path[1:]
	}

	fileInfo, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if !fileInfo.IsDir() {
		return false
	}

	return true
}

func (f *mkDirAllMapFS) MkDirAll(path string, perm os.FileMode) error {
	var current string
	segments := filepath.SplitList(path)

	for _, part := range segments {
		if current == "" {
			current = part
		} else {
			current += string(filepath.Separator) + part
		}

		if exists := f.DirectoryExists(current); !exists {
			f.mapFS[current] = &fstest.MapFile{
				Mode: fs.ModeDir | perm,
			}
		}
	}

	return nil
}

// AbsFunc signature of function used to obtain the absolute representation of
// a path.
type AbsFunc func(path string) (string, error)

// Abs function invoker, allows a function to be used in place where
// an instance of an interface would be expected.
func (f AbsFunc) Abs(path string) (string, error) {
	return f(path)
}

// HomeUserFunc signature of function used to obtain the user's home directory.
type HomeUserFunc func() (string, error)

// Home function invoker, allows a function to be used in place where
// an instance of an interface would be expected.
func (f HomeUserFunc) Home() (string, error) {
	return f()
}

// ResolveMocks, used to override the internal functions used
// to resolve the home path (os.UserHomeDir) and the abs path
// (filepath.Abs). In normal usage, these do not need to be provided,
// just used for testing purposes.
type ResolveMocks struct {
	HomeFunc HomeUserFunc
	AbsFunc  AbsFunc
}
