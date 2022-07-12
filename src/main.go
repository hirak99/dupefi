package main

import (
	"fmt"
	"log"
	"nomen_aliud/duphunter/file_info"
	"os"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
)

var opts struct {
	MinSize      int64  `long:"min-size" default:"1" description:"Minimum file size to include"`
	Verbose      bool   `short:"v" description:"Make it verbose"`
	OutTemplate  string `long:"out-tmpl" default:"$1 -- $2" description:"Output template"`
	BaseTemplate string `long:"base-tmpl" default:"$1" description:"Template for base file"`
	Positional   struct {
		Directory string
	} `positional-args:"yes" required:"yes"`
}

func debugLog(s string, a ...interface{}) {
	if opts.Verbose {
		log.Printf(s, a...)
	}
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
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

func pluralize(count int, name string) string {
	return fmt.Sprintf("%v %v%v", count, name, If(count > 1, "s", ""))
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		if !flags.WroteHelp(err) {
			panic(err)
		}
		return
	}

	files := file_info.ScanDir(opts.Positional.Directory, opts.MinSize)

	// We could call sameSizeDups here, e.g. -
	// fmt.Println(sameSizeDups(files))
	// But that would not take advantage of size based clustering.

	duplicateGroups := findDups(files)

	defer func() {
		debugLog("Checksums computed: %v, Full comparisons: %v", file_info.NChecksums, file_info.NFullComparisons)
	}()

	if len(duplicateGroups) == 0 {
		// Goes to stderr.
		println(fmt.Sprintf("No duplicates found in %v.", pluralize(len(files), "file")))
		return
	}
	for _, group := range duplicateGroups {
		var basePath string
		for i, f := range group {
			if i == 0 {
				basePath = f.Path
				if opts.BaseTemplate != "" {
					out := strings.ReplaceAll(opts.BaseTemplate, "$1", basePath)
					fmt.Println(out)
				}
			} else {
				if opts.OutTemplate != "" {
					out := strings.ReplaceAll(opts.OutTemplate, "$1", basePath)
					out = strings.ReplaceAll(out, "$2", f.Path)
					fmt.Println(out)
				}
			}
		}
	}
}
