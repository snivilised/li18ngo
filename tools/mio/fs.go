package mio

import "os"

const (
	perm = 0o666
)

var _ Transport = (*NativeReaderWriterFS)(nil)

// NativeReaderWriterFS
type NativeReaderWriterFS struct {
	Path string
}

// Address returns the file system location to read and write from
func (w *NativeReaderWriterFS) Address() string {
	return w.Path
}

// Read reads data from the file at the specified path
func (w *NativeReaderWriterFS) Read() ([]byte, error) {
	return os.ReadFile(w.Path)
}

// Write writes data to the file at the specified path
func (w *NativeReaderWriterFS) Write(data []byte) error {
	return os.WriteFile(w.Path, data, perm)
}
