package nfs

import (
	"os"
)

// 📚 package nfs: file system definitions

type (
	// ExistsInFS provides the facility to check the existence
	// of a path in the underlying file system.
	ExistsInFS interface {
		// FileExists does file exist at the path specified
		FileExists(path string) bool

		// DirectoryExists does directory exist at the path specified
		DirectoryExists(path string) bool
	}

	// MkDirAllFS is a file system with a MkDirAll method.
	MkDirAllFS interface {
		ExistsInFS
		MkDirAll(path string, perm os.FileMode) error
	}
)
