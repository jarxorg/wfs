package osfs_test

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/jarxorg/fs2"
	"github.com/jarxorg/fs2/osfs"
)

func ExampleDirFS() {
	tmpDir, err := ioutil.TempDir("", "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	name := "example.txt"
	content := []byte(`Hello`)

	fsys := osfs.DirFS(tmpDir)
	_, err = fs2.WriteFile(fsys, name, content, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	wrote, err := ioutil.ReadFile(tmpDir + "/" + name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(wrote))

	// Output: Hello
}
