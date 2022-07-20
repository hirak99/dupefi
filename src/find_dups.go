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
	"fmt"
	"nomen_aliud/dupefi/file_info"
	"regexp"
	"sort"
	"strings"

	. "github.com/hirak99/go-sanity"
)

func sameSizeDups(files []file_info.FileInfo) [][]file_info.FileInfo {
	var groups [][]int
	for i := 0; i < len(files); i++ {
		found := false
		for gi := range groups {
			// Compare only with the first element in each group.
			if files[i].IsDupOf(&files[groups[gi][0]], opts.Checksum) {
				// Matched, add to group.
				groups[gi] = append(groups[gi], i)
				found = true
				break
			}
		}
		if !found {
			groups = append(groups, []int{i})
		}
	}

	// Convert it into file groups.
	var result [][]file_info.FileInfo
	for _, group := range groups {
		if len(group) == 1 {
			// Don't return single files.
			continue
		}
		var r []file_info.FileInfo
		for _, i := range group {
			r = append(r, files[i])
		}
		result = append(result, r)
	}
	return result
}

func findDups(files []file_info.FileInfo) [][]file_info.FileInfo {
	filesBySize := make(map[int64][]file_info.FileInfo)
	for _, f := range files {
		filesBySize[f.Size] = append(filesBySize[f.Size], f)
	}

	var sizes []int64
	for size := range filesBySize {
		sizes = append(sizes, size)
	}

	// Sort in descending order of sizes.
	sort.Slice(sizes, func(i, j int) bool {
		return sizes[i] > sizes[j]
	})

	debugLog("Sorted sizes: %v\n", sizes)

	var results [][]file_info.FileInfo
	totalBytes := Sum(Map(sizes,
		func(size int64) int64 {
			return size * int64(len(filesBySize[size])-1)
		}))
	var doneBytes int64
	for _, size := range sizes {
		debugLog("Checking %v file(s) of size %v...\n", len(filesBySize[size]), size)
		r := sameSizeDups(filesBySize[size])
		if len(r) > 0 {
			results = append(results, r...)
		}
		doneBytes += size * int64(len(filesBySize[size])-1)
		if !opts.NoProgress && totalBytes > 0 {
			print(fmt.Sprintf("\rComparing %v%% ", int(10000.0*doneBytes/totalBytes)/100))
		}
	}
	if !opts.NoProgress && totalBytes > 0 {
		println("Done.")
	}
	return results
}

// Process a duplicate group.
func postProcessGroup(group []file_info.FileInfo, rnodup *regexp.Regexp) []file_info.FileInfo {
	var result []file_info.FileInfo
	// Copy the group into result.
	result = append(result, group...)

	// Paths matching those which shouldn't be reported as dups.
	// Essentially if there are any such, only one of them must be reported and as the first element.
	nodupset := MakeSet[string]()
	for _, f := range group {
		if (rnodup != nil && rnodup.MatchString(f.Path)) ||
			(opts.Against != "" && strings.HasPrefix(f.Path, opts.Against)) {
			nodupset.Add(f.Path)
		}
	}
	// Less-than function for sorting.
	lessfn := func(f1, f2 file_info.FileInfo) bool {
		p1 := f1.Path
		p2 := f2.Path
		if nodupset.Has(p1) != nodupset.Has(p2) {
			// If p1 is in nodupset, put it at the top.
			return nodupset.HasInt(p1) > nodupset.HasInt(p2)
		}
		return p1 < p2
	}

	SaneSortSlice(result, lessfn)

	// Check if duplicate inodes should be removed.
	if !opts.InodeAsDup {
		// Remove duplicate inodes without re-sorting.
		seen := MakeSet[uint64]()
		var newResult []file_info.FileInfo
		for _, f := range result {
			if seen.Has(f.Inode) {
				continue
			}
			seen.Add(f.Inode)
			newResult = append(newResult, f)
		}
		result = newResult
	}

	// After the first element,
	// drop everything in nodupset.
	result = Filter(result,
		func(i int, f file_info.FileInfo) bool {
			return i == 0 || !nodupset.Has(f.Path)
		})

	if len(result) <= 1 {
		return nil
	}
	if opts.Against != "" && !strings.HasPrefix(result[0].Path, opts.Against) {
		// Do not report internal duplicats when comparing against a directory.
		return nil
	}
	return result
}

func postProcessDups(dups [][]file_info.FileInfo) [][]file_info.FileInfo {
	var result [][]file_info.FileInfo
	rnodup := If(opts.RegexNodup == "", nil, regexp.MustCompile(opts.RegexNodup))
	for _, group := range dups {
		processed := postProcessGroup(group, rnodup)
		if len(processed) > 0 {
			result = append(result, processed)
		}
	}
	return result
}
