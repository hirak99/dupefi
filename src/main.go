package main

import (
	"flag"
	"fmt"
	"log"
	"nomen_aliud/duphunter/file_info"
	"sort"
	"strings"
)

var logLevel int

func debugLog(s string, a ...interface{}) {
	if 1 <= logLevel {
		log.Printf(s, a...)
	}
}

func sameSizeDups(files []file_info.FileInfo) [][]file_info.FileInfo {
	var groups [][]int
	for i := 0; i < len(files); i++ {
		found := false
		for gi := range groups {
			// Compare only with the first element in each group.
			if files[i].IsDupOf(&files[groups[gi][0]]) {
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

func main() {
	minSizeFlag := flag.Int64("min-size", 1, "Minimum file size to include")
	verboseFlag := flag.Bool("verbose", false, "Verbosity")
	outTmpl := flag.String("out-tmpl", "$1 -- $2", "Output")
	baseTmpl := flag.String("base-tmpl", "$1", "Template to print base file for each duplicate group")
	flag.Parse()

	if *verboseFlag {
		logLevel = 1
	}

	files := file_info.ScanCurrentDir(*minSizeFlag)

	// We could call sameSizeDups here, e.g. -
	// fmt.Println(sameSizeDups(files))
	// But that would not take advantage of size based clustering.

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
		debugLog("Checking %v file(s) for size %v...\n", len(filesBySize[size]), size)
		r := sameSizeDups(filesBySize[size])
		if len(r) > 0 {
			debugLog("Found some dups")
			results = append(results, r...)
		}
	}

	if len(results) == 0 {
		// Goes to stderr.
		println("No duplicates found.")
		return
	}
	for _, group := range results {
		var basePath string
		for i, f := range group {
			if i == 0 {
				basePath = f.Path
				if len(*baseTmpl) > 0 {
					out := strings.ReplaceAll(*baseTmpl, "$1", basePath)
					fmt.Println(out)
				}
			} else {
				out := strings.ReplaceAll(*outTmpl, "$1", basePath)
				out = strings.ReplaceAll(out, "$2", f.Path)
				fmt.Println(out)
			}
		}
	}
}
