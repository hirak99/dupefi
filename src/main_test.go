package main

import (
	"io/ioutil"
	"nomen_aliud/duphunter/file_info"
	"os"
	"path"
	"testing"
)

func addFile(name string, content []byte) {
}

func TestHelloName(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "duphunter_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	os.WriteFile(path.Join(dir, "f0"), []byte(""), 0644)
	os.WriteFile(path.Join(dir, "f1"), []byte("Hello1"), 0644)
	os.WriteFile(path.Join(dir, "f2"), []byte("Hello1"), 0644)
	os.WriteFile(path.Join(dir, "f3"), []byte("Hello2"), 0644)
	files := file_info.ScanDir(dir, 1)
	if len(files) != 3 {
		t.Fatalf("Found %v files: %v", len(files), files)
	}

	dups := findDups(files)
	if len(dups) != 1 {
		t.Fatalf("Found %v dup groups", len(dups))
	}

	if len(dups[0]) != 2 {
		t.Fatalf("Found %v dup[0] files", len(dups[0]))
	}
}
