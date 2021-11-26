package osfs

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/jarxorg/fs2"
	"github.com/jarxorg/fs2/fstest2"
)

func TestFS(t *testing.T) {
	fsys := New("testdata")
	if err := fstest.TestFS(fsys, "dir0", "dir0/file01.txt"); err != nil {
		t.Errorf("Error testing/fstest: %+v", err)
	}
}

func TestWriteFileFS(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := New(filepath.Dir(tmpDir))
	if err := fstest2.TestWriteFileFS(fsys, filepath.Base(tmpDir)); err != nil {
		t.Errorf(`Error fs2/fstest2: %+v`, err)
	}
}

func TestMkdirAll(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := NewOSFS(tmpDir)
	err = fsys.MkdirAll("dir", fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = fsys.MkdirAll("../invalid", fs.ModePerm)
	if err == nil {
		t.Fatal(err)
	}
}

func TestCreateFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	got, err := fs2.CreateFile(fsys, "test.txt", fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer got.Close()
}

func TestCreateFile_MkdirAllError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	orgMkdirAllFunc := osMkdirAllFunc
	defer func() { osMkdirAllFunc = orgMkdirAllFunc }()

	wantErr := errors.New("test")
	osMkdirAllFunc = func(dir string, perm os.FileMode) error {
		return wantErr
	}

	fsys := DirFS(tmpDir)
	var gotErr error
	_, gotErr = fs2.CreateFile(fsys, "name.txt", fs.ModePerm)

	if !reflect.DeepEqual(gotErr, wantErr) {
		t.Errorf("Error CreateFile returns unknown error %v; want %v", gotErr, wantErr)
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	name := "test.txt"
	want := []byte(`test`)

	fsys := DirFS(tmpDir)
	n, err := fs2.WriteFile(fsys, name, want, fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(want) {
		t.Errorf("Error len %d; want %d", n, len(want))
	}

	got, err := ioutil.ReadFile(tmpDir + "/" + name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error content %s; want %s", got, want)
	}
}

func TestWriteFile_InvalidError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	_, err = fs2.WriteFile(fsys, "../invalid.txt", []byte{}, fs.ModePerm)
	if err == nil {
		t.Fatal("Error WriteFile returns no error")
	}
}

func TestContainsDenyWin(t *testing.T) {
	testCases := []struct {
		name string
		want bool
	}{
		{
			name: `allow.txt`,
			want: false,
		}, {
			name: `path/to/allow.txt`,
			want: false,
		}, {
			name: `deny:txt`,
			want: true,
		}, {
			name: `C:/deny.txt`,
			want: true,
		}, {
			name: `path\to\deny.txt`,
			want: true,
		},
	}
	for i, testCase := range testCases {
		got := containsDenyWin(testCase.name)
		if got != testCase.want {
			t.Errorf("Error[%d] containsDenyWin(%s) %v; want %v",
				i, testCase.name, got, testCase.want)
		}
	}
}

func TestSub_WriteFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	dir := "sub"
	name := "test.txt"
	want := []byte(`test`)

	fsys, err := fs.Sub(DirFS(tmpDir), dir)
	if err != nil {
		t.Fatal(err)
	}
	n, err := fs2.WriteFile(fsys, name, want, fs.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(want) {
		t.Errorf("Error len %d; want %d", n, len(want))
	}

	got, err := ioutil.ReadFile(tmpDir + "/" + dir + "/" + name)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Error content %s; want %s", got, want)
	}
}

func TestRemoveFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	name := "test.txt"

	if err = ioutil.WriteFile(tmpDir+"/"+name, []byte{}, fs.ModePerm); err != nil {
		t.Fatal(err)
	}

	err = fs2.RemoveFile(fsys, name)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveFile_InvalidError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	err = fs2.RemoveFile(fsys, "../invalid-dir")
	if err == nil {
		t.Fatal("Error RemoveFile returns no error")
	}
}

func TestRemoveAll(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	path := "dir"
	name := "test.txt"

	if err = os.Mkdir(tmpDir+"/"+path, fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err = ioutil.WriteFile(tmpDir+"/"+path+"/"+name, []byte{}, fs.ModePerm); err != nil {
		t.Fatal(err)
	}

	err = fs2.RemoveAll(fsys, path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRemoveAll_InvalidError(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	fsys := DirFS(tmpDir)
	err = fs2.RemoveAll(fsys, "../invalid-dir")
	if err == nil {
		t.Fatal("Error RemoveAll returns no error")
	}
}
