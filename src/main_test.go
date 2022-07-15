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

// Create a temp directory, and set up some common files.
// Returns a cleanup function that must be called as defer.
func setupCommonFiles() func() {
	dir, err := ioutil.TempDir(os.TempDir(), "duphunter_test")
	if err != nil {
		panic(err)
	}

	os.Chdir(dir)

	os.WriteFile(path.Join(dir, "f0"), []byte(""), 0644)
	os.WriteFile(path.Join(dir, "f1"), []byte("Hello1"), 0644)
	os.WriteFile(path.Join(dir, "f2"), []byte("Hello Dup"), 0644)
	os.WriteFile(path.Join(dir, "f3"), []byte("Hello2"), 0644)
	os.Mkdir(path.Join(dir, "subd1"), 0744)
	os.WriteFile(path.Join(dir, "subd1", "f4.txt"), []byte("Hello Dup"), 0644)
	os.WriteFile(path.Join(dir, "subd1", "f5.txt"), []byte("Hello"), 0644)

	return func() {
		os.RemoveAll(dir)
	}
}

func TestRegexpScan(t *testing.T) {
	defer setupCommonFiles()()

	// Test regexp.
	AssertSliceEqual(t,
		Map(
			ChanToSlice(file_info.ScanDir(".", 1, regexp.MustCompile(`\.txt$`))),
			func(f file_info.FileInfo) string { return f.Path }),
		[]string{"subd1/f4.txt", "subd1/f5.txt"})
}

func testDuphunting(t *testing.T) {
	defer setupCommonFiles()()

	files := ChanToSlice(file_info.ScanDir(".", 1, nil))
	// We don't expect f0 since it has zero length.
	AssertSliceEqual(t,
		Map(files, func(f file_info.FileInfo) string { return f.Path }),
		[]string{"f1", "f2", "f3", "subd1/f4.txt", "subd1/f5.txt"})

	dups := postProcessDups(findDups(files))
	AssertEqual(t, len(dups), 1)
	AssertEqual(t, len(dups[0]), 2)

	opts.BaseTemplate = "$1"
	opts.OutTemplate = "$0 -- $1"
	AssertSliceEqual(t,
		ChanToSlice(getDisplayLines(dups)),
		[]string{"f2", "f2 -- subd1/f4.txt"})
	opts.BaseTemplate = ""
	opts.OutTemplate = "$0 -- $1"
	AssertSliceEqual(t,
		ChanToSlice(getDisplayLines(dups)),
		[]string{"f2 -- subd1/f4.txt"})
}

func TestWithChecksum(t *testing.T) {
	opts.Checksum = true
	testDuphunting(t)
}

func TestWithoutChecksum(t *testing.T) {
	opts.Checksum = false
	testDuphunting(t)
}

func TestPostProcessAllInodesSame(t *testing.T) {
	group := []file_info.FileInfo{
		file_info.FakeFile("f1", 100, 2001),
		file_info.FakeFile("f2", 100, 2001),
		file_info.FakeFile("f3", 100, 2001),
	}
	func() {
		// Removing duplicates.
		opts.InodeAsDup = false
		result := postProcessGroup(group, nil)
		AssertEqual(t, len(result), 0)
	}()
}

func TestPostProcessDup(t *testing.T) {
	group := []file_info.FileInfo{
		file_info.FakeFile("f1", 100, 2001),
		file_info.FakeFile("f2", 100, 2002),
		file_info.FakeFile("f3", 100, 2003),
		file_info.FakeFile("g1", 100, 2001),
		file_info.FakeFile("g2", 100, 2004),
		file_info.FakeFile("g3", 100, 2001),
	}
	{
		// Removing inode duplicates.
		opts.InodeAsDup = false
		result := postProcessGroup(group, nil)
		AssertSliceEqual(t,
			Map(result, func(fi file_info.FileInfo) string { return fi.Path }),
			[]string{"f1", "f2", "f3", "g2"})
	}
	{
		// Not removing inode duplicates.
		opts.InodeAsDup = true
		result := postProcessGroup(group, nil)
		AssertSliceEqual(t,
			Map(result, func(fi file_info.FileInfo) string { return fi.Path }),
			[]string{"f1", "f2", "f3", "g1", "g2", "g3"})
	}
	{
		// Not removing inode duplicates.
		opts.InodeAsDup = true
		// Nodupregex satisfies all files.
		result := postProcessGroup(group, regexp.MustCompile("^(f|g)"))
		AssertEqual(t, len(result), 0)
	}
	{
		// Not removing inode duplicates.
		opts.InodeAsDup = true
		// Nodupregex satisfies all file starting with g.
		result := postProcessGroup(group, regexp.MustCompile("^g"))
		AssertSliceEqual(t,
			Map(result, func(fi file_info.FileInfo) string { return fi.Path }),
			[]string{"g1", "f1", "f2", "f3"})
	}
}
