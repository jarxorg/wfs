package fs2

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestMkdirAll(t *testing.T) {
	got := ""
	fsys := &FSDelegator{
		MkdirAllFunc: func(dir string, _ fs.FileMode) error {
			got = dir
			return nil
		},
	}

	want := "path/to/dir"
	err := MkdirAll(fsys, want, fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Error MkdirAll called with %s; want %s", got, want)
	}
}

func TestMkdirAll_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	dir := "path/to/dir"
	wantErr := &fs.PathError{Op: "MkdirAll", Path: dir, Err: ErrNotImplemented}

	err := MkdirAll(fsys, dir, fs.ModePerm)
	if err == nil {
		t.Errorf("Error MkdirAll returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error MkdirAll returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error MkdirAll returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestCreateFile(t *testing.T) {
	want := &FileDelegator{}
	called := false
	fsys := &FSDelegator{
		CreateFileFunc: func(_ string, _ fs.FileMode) (WriterFile, error) {
			called = true
			return want, nil
		},
	}

	got, err := CreateFile(fsys, "test.txt", fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error CreateFile is not called")
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error CreateFile returns %v; want %v", got, want)
	}
}

func TestCreateFile_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	name := "test.txt"
	wantErr := &fs.PathError{Op: "CreateFile", Path: name, Err: ErrNotImplemented}

	var err error
	_, err = CreateFile(fsys, name, fs.ModePerm)
	if err == nil {
		t.Errorf("Error CreateFile returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error CreateFile returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error CreateFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestWriteFile(t *testing.T) {
	want := 1
	called := false
	fsys := &FSDelegator{
		WriteFileFunc: func(_ string, _ []byte, _ fs.FileMode) (int, error) {
			called = true
			return want, nil
		},
	}

	got, err := WriteFile(fsys, "", []byte{}, fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error WriteFile is not called")
	}
	if got != want {
		t.Errorf("Error WriteFile returns %d; want %d", got, want)
	}
}

func TestWriteFile_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	name := "test.txt"
	wantErr := &fs.PathError{Op: "WriteFile", Path: name, Err: ErrNotImplemented}

	var err error
	_, err = WriteFile(fsys, name, []byte{}, fs.ModePerm)
	if err == nil {
		t.Errorf("Error WriteFile returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error WriteFile returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error WriteFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestRemoveFile(t *testing.T) {
	called := false
	fsys := &FSDelegator{
		RemoveFileFunc: func(name string) error {
			called = true
			return nil
		},
	}

	err := RemoveFile(fsys, "")
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error RemoveFile is not called")
	}
}

func TestRemoveAll(t *testing.T) {
	called := false
	fsys := &FSDelegator{
		RemoveAllFunc: func(name string) error {
			called = true
			return nil
		},
	}

	err := RemoveAll(fsys, "")
	if err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Error("Error RemoveAll is not called")
	}
}

func TestRemoveFile_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	name := "test.txt"
	wantErr := &fs.PathError{Op: "RemoveFile", Path: name, Err: ErrNotImplemented}

	err := RemoveFile(fsys, name)
	if err == nil {
		t.Errorf("Error RemoveFile returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error RemoveFile returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error RemoveFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestRemoveAll_ErrNotImplemented(t *testing.T) {
	fsys := &OpenFSDelegator{}

	path := "path/to/dir"
	wantErr := &fs.PathError{Op: "RemoveAll", Path: path, Err: ErrNotImplemented}

	err := RemoveAll(fsys, path)
	if err == nil {
		t.Errorf("Error RemoveAll returns no error")
	}
	gotErr, ok := err.(*fs.PathError)
	if !ok {
		t.Errorf("Error RemoveAll returns unknown error %v", err)
	}
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error RemoveAll returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestCopyFS(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	root := "dir0"
	want := map[string][]byte{}
	src := os.DirFS("osfs/testdata")
	err = fs.WalkDir(src, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		p, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}
		want[path] = p
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	got := map[string][]byte{}
	dest := DelegateFS(os.DirFS(tmpDir))
	dest.CreateFileFunc = func(name string, mode fs.FileMode) (WriterFile, error) {
		return &FileDelegator{
			WriteFunc: func(p []byte) (int, error) {
				got[name] = p
				return len(p), nil
			},
		}, nil
	}

	err = CopyFS(dest, src, root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error CopyFS %v; want %v", got, want)
	}
}

func TestCopyFS_StatError(t *testing.T) {
	wantErr := errors.New("test")

	src := DelegateFS(os.DirFS("osfs/testdata"))
	src.StatFunc = func(name string) (fs.FileInfo, error) {
		return nil, wantErr
	}

	gotErr := CopyFS(&FSDelegator{}, src, ".")
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error CopyFS returns unknown error %+v; want %v", gotErr, wantErr)
	}
}

func TestCopyFS_OpenError(t *testing.T) {
	wantErr := errors.New("test")

	src := DelegateFS(os.DirFS("osfs/testdata"))
	src.OpenFunc = func(name string) (fs.File, error) {
		return nil, wantErr
	}

	gotErr := CopyFS(&FSDelegator{}, src, ".")
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error CopyFS returns unknown error %+v; want %v", gotErr, wantErr)
	}
}

func TestCopyFS_CreateFileError(t *testing.T) {
	wantErr := errors.New("test")

	src := os.DirFS("osfs/testdata")
	dest := &FSDelegator{
		CreateFileFunc: func(_ string, _ fs.FileMode) (WriterFile, error) {
			return nil, wantErr
		},
	}

	gotErr := CopyFS(dest, src, ".")
	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error CopyFS returns unknown error %+v; want %v", gotErr, wantErr)
	}
}
