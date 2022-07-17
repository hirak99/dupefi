// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"nomen_aliud/dupefi/file_info"

	. "github.com/hirak99/go-sanity"
)

func TestEmptyDir(t *testing.T) {
	dir, err := ioutil.TempDir(t.TempDir(), "dupefi_test")
	if err != nil {
		panic(err)
	}

	os.Chdir(dir)

	files := file_info.ScanDirs([]string{"."}, 1, nil)
	AssertEqual(t, len(files), 0)
	dups := findDups(files)
	AssertEqual(t, len(dups), 0)
}

// Create a temp directory, and set up some common files.
// Returns a cleanup function that must be called as defer.
func setupCommonFiles(t *testing.T) {
	dir, err := ioutil.TempDir(t.TempDir(), "dupefi_test")
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
}

func TestRegexpScan(t *testing.T) {
	setupCommonFiles(t)

	// Test regexp.
	AssertSliceEqual(t,
		Map(
			file_info.ScanDirs([]string{"."}, 1, regexp.MustCompile(`\.txt$`)),
			func(f file_info.FileInfo) string { return f.Path }),
		[]string{"subd1/f4.txt", "subd1/f5.txt"})
}

func testDuphunting(t *testing.T) {
	setupCommonFiles(t)

	files := file_info.ScanDirs([]string{"."}, 1, nil)
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
		file_info.FakeFile("f2", 100, 2004),
		file_info.FakeFile("f3", 100, 2003),
		file_info.FakeFile("g1", 100, 2001),
		file_info.FakeFile("g2", 100, 2002),
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
