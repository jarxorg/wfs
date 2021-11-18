package memfs_test

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/jarxorg/fs2"
	"github.com/jarxorg/fs2/memfs"
)

func ExampleNew() {
	name := "path/to/example.txt"
	content := []byte(`Hello`)

	fsys := memfs.New()
	var err error
	_, err = fs2.WriteFile(fsys, name, content, fs.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	wrote, err := fs.ReadFile(fsys, name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(wrote))

	// Output: Hello
}
