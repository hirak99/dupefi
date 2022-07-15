package main

import (
	"nomen_aliud/duphunter/file_info"
	"sort"
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
	for _, size := range sizes {
		// debugLog("Checking %v file(s) for size %v...\n", len(filesBySize[size]), size)
		r := sameSizeDups(filesBySize[size])
		if len(r) > 0 {
			debugLog("Found some dups")
			results = append(results, r...)
		}
	}
	return results
}

// Process a duplicate group.
func postProcessGroup(group []file_info.FileInfo) []file_info.FileInfo {
	var result []file_info.FileInfo
	result = append(result, group...)
	if !opts.InodeAsDup {
		sort.Slice(result, func(i, j int) bool {
			return result[i].Inode < result[j].Inode
		})
		result = Filter(result,
			func(i int, _ file_info.FileInfo) bool {
				return i == 0 || result[i].Inode != result[i-1].Inode
			})
	}
	if len(result) <= 1 {
		return nil
	}
	return result
}

func postProcessDups(dups [][]file_info.FileInfo) [][]file_info.FileInfo {
	var result [][]file_info.FileInfo
	for _, group := range dups {
		processed := postProcessGroup(group)
		if len(processed) > 0 {
			result = append(result, processed)
		}
	}
	return result
}
