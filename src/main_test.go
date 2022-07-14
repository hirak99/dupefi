package main

import (
	"io/ioutil"
	"nomen_aliud/duphunter/file_info"
	"os"
	"path"
	"regexp"
	"testing"
)

func TestEmptyDir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "duphunter_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	os.Chdir(dir)

	files := ChanToSlice(file_info.ScanDir(".", 1, nil))
	AssertEqual(t, len(files), 0)
	dups := findDups(files)
	AssertEqual(t, len(dups), 0)
}

func TestDuphunting(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "duphunter_test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	os.Chdir(dir)

	os.WriteFile(path.Join(dir, "f0"), []byte(""), 0644)
	os.WriteFile(path.Join(dir, "f1"), []byte("Hello1"), 0644)
	os.WriteFile(path.Join(dir, "f2"), []byte("Hello Dup"), 0644)
	os.WriteFile(path.Join(dir, "f3"), []byte("Hello2"), 0644)
	os.Mkdir(path.Join(dir, "subd1"), 0744)
	os.WriteFile(path.Join(dir, "subd1", "f4.txt"), []byte("Hello Dup"), 0644)
	os.WriteFile(path.Join(dir, "subd1", "f5.txt"), []byte("Hello"), 0644)

	// Test regexp.
	AssertSliceEqual(t,
		Map(
			ChanToSlice(file_info.ScanDir(".", 1, regexp.MustCompile(`\.txt$`))),
			func(f file_info.FileInfo) string { return f.Path }),
		[]string{"subd1/f4.txt", "subd1/f5.txt"})

	files := ChanToSlice(file_info.ScanDir(".", 1, nil))
	// We don't expect f0 since it has zero length.
	AssertSliceEqual(t,
		Map(files, func(f file_info.FileInfo) string { return f.Path }),
		[]string{"f1", "f2", "f3", "subd1/f4.txt", "subd1/f5.txt"})

	dups := findDups(files)
	AssertEqual(t, len(dups), 1)
	AssertEqual(t, len(dups[0]), 2)

	AssertSliceEqual(t,
		ChanToSlice(getDisplayLines(dups, "$1", "$0 -- $1")),
		[]string{"f2", "f2 -- subd1/f4.txt"})
	AssertSliceEqual(t,
		ChanToSlice(getDisplayLines(dups, "", "$0 -- $1")),
		[]string{"f2 -- subd1/f4.txt"})
}
