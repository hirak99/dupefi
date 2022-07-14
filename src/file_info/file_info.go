package file_info

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"syscall"
)

type FileInfo struct {
	Path string
	Size int64

	info  os.FileInfo
	inode uint64
	// Don't access directly, call getChecksum() instead.
	checksum *string
}

// How many bytes to read at once, for check-summing or file comparisons.
const bufferSize = 1024 * 4 // 4 KiB

// Total number of checksums computed.
var NChecksums int

// Total number of full comparisons done
var NFullComparisons int

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
	if f1.inode != 0 && f1.inode == f2.inode {
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

func ScanDir(dir string, minSize int64, r *regexp.Regexp) <-chan FileInfo {
	out := make(chan FileInfo)
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
					return errors.New(fmt.Sprintf("Could not perform syscall on %v", path))
				} else {
					inode = stat.Ino
				}
				if r == nil || r.MatchString(path) {
					out <- FileInfo{Path: path, info: info, inode: inode, Size: info.Size()}
				}
				return nil
			})
		close(out)
		if err != nil {
			log.Println(err)
		}
	}()
	return out
}
