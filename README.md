# wfs

[![PkgGoDev](https://pkg.go.dev/badge/github.com/jarxorg/wfs)](https://pkg.go.dev/github.com/jarxorg/wfs)
[![Report Card](https://goreportcard.com/badge/github.com/jarxorg/wfs)](https://goreportcard.com/report/github.com/jarxorg/wfs)

Package wfs provides writable [io/fs](https://pkg.go.dev/io/fs).FS interfaces.

```go
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

// RemoveFileFS is the interface implemented by a filesystem that provides an
// implementation of RemoveFile.
type RemoveFileFS interface {
	fs.FS
	RemoveFile(name string) error
	RemoveAll(name string) error
}
```

This is one of the solutions to an [issue](https://github.com/golang/go/issues/45757) of github.com/golango/go.

The following packages are an implementation of wfs.

- [osfs](https://pkg.go.dev/github.com/jarxorg/wfs/osfs)
- [memfs](https://pkg.go.dev/github.com/jarxorg/wfs/memfs)
- [s3fs](https://github.com/jarxorg/s3fs)
- [gcsfs](https://github.com/jarxorg/gcsfs)

## CopyFS

CopyFS walks the specified root directory on src and copies directories and files to dest filesystem.
The following code is an example.

```go
package main

import (
	"log"

	"github.com/jarxorg/s3fs"
	"github.com/jarxorg/wfs"
	"github.com/jarxorg/wfs/osfs"
)

func main() {
	src := s3fs.New("your-bucket")
	dst := osfs.DirFS("local-dir")

	// NOTE: Copy files on s3://your-bucket to local-dir.
	if err := wfs.CopyFS(dst, src, "."); err != nil {
		log.Fatal(err)
	}
}
```