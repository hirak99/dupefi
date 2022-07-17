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

package file_info

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"syscall"

	. "github.com/hirak99/go-sanity"
)

type FileInfo struct {
	Path string
	Size int64

	info  os.FileInfo
	Inode uint64
	// Don't access directly, call getChecksum() instead.
	checksum *string
}

// How many bytes to read at once, for check-summing or file comparisons.
const bufferSize = 1024 * 4 // 4 KiB

// Total number of checksums computed.
var NChecksums int

// Total number of full comparisons done
var NFullComparisons int

// Used for mocking in tests.
func FakeFile(path string, size int64, inode uint64) FileInfo {
	return FileInfo{Path: path, Size: size, Inode: inode}
}

func (f *FileInfo) getChecksum() string {
	if f.checksum == nil {
		NChecksums++

		hasher := sha256.New()

		handle, err := os.Open(f.Path)
		if err != nil {
			log.Fatal(err)
		}

		for {
			buf := make([]byte, bufferSize)
			n, err := handle.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			hasher.Write(buf[:n])
		}

		f.checksum = new(string)
		*f.checksum = fmt.Sprintf("%x", hasher.Sum(nil))
	}
	return *f.checksum
}

// Returns true if files are same.
func compare(path1, path2 string) bool {
	NFullComparisons++

	f1, err := os.Open(path1)
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()
	f2, err := os.Open(path2)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer f2.Close()

	for {
		b1 := make([]byte, bufferSize)
		_, err1 := f1.Read(b1)

		b2 := make([]byte, bufferSize)
		_, err2 := f2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return true
			} else if err1 != nil || err2 != nil {
				log.Fatal(err1, err2)
				return false
			} else {
				return false
			}
		}
		if !bytes.Equal(b1, b2) {
			return false
		}
	}

}

func (f1 *FileInfo) IsDupOf(f2 *FileInfo, useChecksum bool) bool {
	if f1.Inode != 0 && f1.Inode == f2.Inode {
		return true
	}
	if f1.Size != f2.Size {
		return false
	}
	if useChecksum {
		return f1.getChecksum() == f2.getChecksum()
	}
	return compare(f1.Path, f2.Path)
}

func scanDir(dir string, minSize int64, r *regexp.Regexp, showProgress bool) <-chan FileInfo {
	out := make(chan FileInfo)
	nfiles := 0
	go func() {
		err := filepath.Walk(dir,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() || info.Size() < minSize {
					return nil
				}
				var inode uint64
				stat, ok := info.Sys().(*syscall.Stat_t)
				if !ok {
					return fmt.Errorf("could not perform syscall on %v", path)
				} else {
					inode = stat.Ino
				}
				if r == nil || r.MatchString(path) {
					out <- FileInfo{Path: path, info: info, Inode: inode, Size: info.Size()}
					nfiles += 1
					if showProgress {
						print(fmt.Sprintf("\rFiles found: %v", nfiles))
					}
				}
				return nil
			})
		close(out)
		if showProgress {
			println()
		}
		if err != nil {
			log.Println(err)
		}
	}()
	return out
}

func ScanDirs(dirs []string, minSize int64, r *regexp.Regexp, showProgress bool) []FileInfo {
	var files []FileInfo
	seen := MakeSet[string]()
	for _, dir := range dirs {
		newFiles := ChanToSlice(scanDir(dir, minSize, r, showProgress))
		for _, f := range newFiles {
			if seen.Has(f.Path) {
				continue
			}
			files = append(files, f)
			seen.Add(f.Path)
		}
	}
	return files
}
