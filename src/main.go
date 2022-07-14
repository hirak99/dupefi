package main

import (
	"fmt"
	"log"
	"nomen_aliud/duphunter/file_info"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
)

// Stores the git hash passed in build script.
var Githash string

var opts struct {
	MinSize      int64  `long:"minsize" description:"Minimum file size to include" default:"1"`
	OutTemplate  string `long:"outtmpl" description:"Output template" default:"$0 -- $1"`
	BaseTemplate string `long:"basetmpl" description:"Template for base file"`
	Regex        string `long:"regex" description:"Regular expression to filter files, e.g. '\\.jpg$'"`
	ShowVersion  bool   `long:"version" description:"Show the version and exit"`

	Verbose    bool `short:"v" description:"Verbose to print additional outputs"`
	Checksum   bool `short:"c" description:"Use checksum instead of full compare"`
	InodeAsDup bool `short:"i" description:"Include multiple hardlinks to the same inode as duplicates"`

	Positional struct {
		Directory string
	} `positional-args:"yes"`
}

func debugLog(s string, a ...interface{}) {
	if opts.Verbose {
		log.Printf(s, a...)
	}
}

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

func pluralize(count int, name string) string {
	return fmt.Sprintf("%v %v%v", count, name, If(count > 1, "s", ""))
}

func getDisplayLines(duplicateGroups [][]file_info.FileInfo, baseTemplate string, outTemplate string) <-chan string {
	out := make(chan string)
	go func() {
		for _, group := range duplicateGroups {
			var basePath string
			for i, f := range group {
				if i == 0 {
					basePath = f.Path
					if baseTemplate != "" {
						out <- strings.ReplaceAll(baseTemplate, "$1", basePath)
					}
				} else {
					if outTemplate != "" {
						line := strings.ReplaceAll(outTemplate, "$1", f.Path)
						line = strings.ReplaceAll(line, "$0", basePath)
						out <- line
					}
				}
			}
		}
		close(out)
	}()
	return out
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

func main() {
	_, err := flags.ParseArgs(&opts, os.Args[1:])
	if err != nil {
		if !flags.WroteHelp(err) {
			panic(err)
		}
		return
	}

	if opts.ShowVersion {
		fmt.Printf("Git commit hash: %s\n", Githash)
		return
	}

	if opts.Positional.Directory == "" {
		opts.Positional.Directory = "."
	}

	regex := If(opts.Regex == "", nil, regexp.MustCompile(opts.Regex))
	files := ChanToSlice(file_info.ScanDir(opts.Positional.Directory, opts.MinSize, regex))

	// We could call sameSizeDups here, e.g. -
	// fmt.Println(sameSizeDups(files))
	// But that would not take advantage of size based clustering.

	duplicateGroups := postProcessDups(findDups(files))

	defer func() {
		debugLog("Checksums computed: %v, Full comparisons: %v", file_info.NChecksums, file_info.NFullComparisons)
	}()

	if len(duplicateGroups) == 0 {
		// Goes to stderr.
		println(fmt.Sprintf("No duplicates found in %v.", pluralize(len(files), "file")))
		return
	}
	for line := range getDisplayLines(duplicateGroups, opts.BaseTemplate, opts.OutTemplate) {
		fmt.Println(line)
	}
}
