package main

import (
	"io/ioutil"
	"nomen_aliud/duphunter/file_info"
	"os"
	"path"
	"reflect"
	"testing"
)

func Map_[T any, U any](data []T, f func(T) U) []U {
	mapped := make([]U, len(data))
	for i, e := range data {
		mapped[i] = f(e)
	}
	return mapped
}

func TestHelloName(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "duphunter_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	os.Chdir(dir)

	os.WriteFile(path.Join(dir, "f0"), []byte(""), 0644)
	os.WriteFile(path.Join(dir, "f1"), []byte("Hello1"), 0644)
	os.WriteFile(path.Join(dir, "f2"), []byte("Hello1"), 0644)
	os.WriteFile(path.Join(dir, "f3"), []byte("Hello2"), 0644)

	files := file_info.ScanDir(".", 1)
	if len(files) != 3 {
		t.Fatalf("Found %v files: %v", len(files), files)
	}

	names := Map_(files, func(f file_info.FileInfo) string {
		return f.Path
	})
	// We don't expect f0 since it has zero length.
	want := []string{"f1", "f2", "f3"}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("Names are not the same. want: %v, got: %v", want, names)
	}

	dups := findDups(files)
	if len(dups) != 1 {
		t.Fatalf("Found %v dup groups", len(dups))
	}

	if len(dups[0]) != 2 {
		t.Fatalf("Found %v dup[0] files", len(dups[0]))
	}

	lines := getDisplayLines(dups, "$1", "$0 -- $1")
	if !reflect.DeepEqual(lines, []string{"f1", "f1 -- f2"}) {
		t.Fatalf("Incorrect display %v", lines)
	}
}
