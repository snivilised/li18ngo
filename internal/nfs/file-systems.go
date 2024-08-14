package nfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type readDirFS struct {
	fsys fs.FS
}

// NewReadDirFS creates a ReadDirFS file system
func NewReadDirFS(path string) fs.ReadDirFS {
	return &readDirFS{
		fsys: os.DirFS(path),
	}
}

// Open opens the named file.
func (n *readDirFS) Open(path string) (fs.File, error) {
	return n.fsys.Open(path)
}

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
func (n *readDirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(n.fsys, name)
}

type queryStatusFS struct {
	fsys fs.FS
}

func NewQueryStatusFS(fsys fs.FS) fs.StatFS {
	return &queryStatusFS{
		fsys: fsys,
	}
}

// Open opens the named file.
func (q *queryStatusFS) Open(name string) (fs.File, error) {
	return q.fsys.Open(name)
}

// Stat returns a FileInfo describing the named file.
// If there is an error, it will be of type *PathError.
func (q *queryStatusFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}

type nativeDirFS struct {
	statFS fs.StatFS
}

// NewNativeDirFS creates an instance of MkDirAllFS from a path
func NewNativeDirFS(path string) MkDirAllFS {
	return &nativeDirFS{
		statFS: NewQueryStatusFS(NewReadDirFS(path)),
	}
}

// FromNativeDirFS creates an instance of MkDirAllFS from a fs.FS
func FromNativeDirFS(fsys fs.FS) MkDirAllFS {
	return &nativeDirFS{
		statFS: NewQueryStatusFS(fsys),
	}
}

// FileExists return true if item at path exists as a file
func (f *nativeDirFS) FileExists(path string) bool {
	fi, err := f.statFS.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	return true
}

// DirectoryExists return true if item at path exists as a directory
func (f *nativeDirFS) DirectoryExists(path string) bool {
	if strings.HasPrefix(path, string(filepath.Separator)) {
		path = path[1:]
	}

	fileInfo, err := f.statFS.Stat(path)
	if err != nil {
		return false
	}

	if !fileInfo.IsDir() {
		return false
	}

	return true
}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
func (f *nativeDirFS) MkDirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}
