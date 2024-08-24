package ifs

import (
	"io/fs"
	"os"

	"github.com/snivilised/li18ngo/nfs"
)

// NewStatFS creates a new fs.StatFS from a path
func NewStatFS(path string) fs.StatFS {
	return StatFSFromFS(os.DirFS(path))
}

type readDirFS struct {
	fsys fs.FS
}

// NewReadDirFS creates a ReadDirFS file system, which
// contains the ReadDir method that reads entries of a path
// from the native file system.
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

// StatFSFromFS creates a file system upon which Stat can be invoked
func StatFSFromFS(fsys fs.FS) fs.StatFS {
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
func NewNativeDirFS(path string) nfs.MkDirAllFS {
	return &nativeDirFS{
		statFS: StatFSFromFS(NewReadDirFS(path)),
	}
}

// DirFSFromFS creates a native instance of MkDirAllFS from a fs.FS
func DirFSFromFS(fsys fs.FS) nfs.MkDirAllFS {
	return &nativeDirFS{
		statFS: StatFSFromFS(fsys),
	}
}

// FileExists return true if item at path exists as a file
func (f *nativeDirFS) FileExists(path string) bool {
	info, err := f.statFS.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

// DirectoryExists return true if item at path exists as a directory
func (f *nativeDirFS) DirectoryExists(path string) bool {
	info, err := f.statFS.Stat(path)
	if err != nil {
		return false
	}

	if !info.IsDir() {
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
