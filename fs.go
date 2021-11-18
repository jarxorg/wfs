package fs2

import (
	"errors"
	"io"
	"io/fs"
)

var (
	// ErrNotImplemented "not implemented"
	ErrNotImplemented = errors.New("not implemented")
)

// WriterFile is a file that provides an implementation fs.File and io.Writer.
type WriterFile interface {
	fs.File
	io.Writer
}

// WriteFileFS is the interface implemented by a filesystem that provides an
// optimized implementation of MkdirAll, CreateFile, WriteFile.
type WriteFileFS interface {
	fs.FS
	MkdirAll(dir string, mode fs.FileMode) error
	CreateFile(name string, mode fs.FileMode) (WriterFile, error)
	WriteFile(name string, p []byte, mode fs.FileMode) (n int, err error)
}

// MkdirAll creates the named directory. If the filesystem implements
// WriteFileFS calls fsys.MkdirAll otherwise returns a PathError.
func MkdirAll(fsys fs.FS, dir string, mode fs.FileMode) error {
	if fsys, ok := fsys.(WriteFileFS); ok {
		return fsys.MkdirAll(dir, mode)
	}
	return &fs.PathError{Op: "MkdirAll", Path: dir, Err: ErrNotImplemented}
}

// CreateFile creates the named file. If the filesystem implements
// WriteFileFS calls fsys.CreateFile otherwise returns a PathError.
func CreateFile(fsys fs.FS, name string, mode fs.FileMode) (WriterFile, error) {
	if fsys, ok := fsys.(WriteFileFS); ok {
		return fsys.CreateFile(name, mode)
	}
	return nil, &fs.PathError{Op: "CreateFile", Path: name, Err: ErrNotImplemented}
}

// WriteFile writes the specified bytes to the named file. If the filesystem implements
// WriteFileFS calls fsys.WriteFile otherwise returns a PathError.
func WriteFile(fsys fs.FS, name string, p []byte, mode fs.FileMode) (n int, err error) {
	if fsys, ok := fsys.(WriteFileFS); ok {
		return fsys.WriteFile(name, p, mode)
	}
	return 0, &fs.PathError{Op: "WriteFile", Path: name, Err: ErrNotImplemented}
}

// RemoveFileFS is the interface implemented by a filesystem that provides an
// implementation of RemoveFile.
type RemoveFileFS interface {
	fs.FS
	RemoveFile(name string) error
	RemoveAll(name string) error
}

// RemoveFile removes the specified named file. If the filesystem implements
// RemoveFileFS calls fsys.RemoveFile otherwise return a PathError.
func RemoveFile(fsys fs.FS, name string) error {
	if fsys, ok := fsys.(RemoveFileFS); ok {
		return fsys.RemoveFile(name)
	}
	return &fs.PathError{Op: "RemoveFile", Path: name, Err: ErrNotImplemented}
}

// RemoveAll removes path and any children it contains. If the filesystem
// implements RemoveFileFS calls fsys.RemoveAll otherwise return a PathError.
func RemoveAll(fsys fs.FS, path string) error {
	if fsys, ok := fsys.(RemoveFileFS); ok {
		return fsys.RemoveAll(path)
	}
	return &fs.PathError{Op: "RemoveAll", Path: path, Err: ErrNotImplemented}
}

// CopyFS walks the specified root directory on src and copies directories and
// files to dest filesystem.
func CopyFS(dest, src fs.FS, root string) error {
	return fs.WalkDir(src, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d == nil {
			return err
		}
		if d.IsDir() {
			return MkdirAll(dest, path, d.Type())
		}
		srcFile, err := src.Open(path)
		if err != nil {
			return err
		}
		destFile, err := CreateFile(dest, path, d.Type())
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		return err
	})
}
