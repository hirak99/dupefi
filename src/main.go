package main

import (
	"fmt"
	"log"
	"nomen_aliud/duphunter/file_info"
	"regexp"
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
	RegexNodup   string `long:"regexnodup" description:"Regular expression to specify files not to be reported as dups"`
	ShowVersion  bool   `long:"version" description:"Show the version and exit"`

	Verbose    bool `short:"v" description:"Verbose to print additional outputs"`
	Checksum   bool `short:"c" description:"Use checksum instead of full compare"`
	InodeAsDup bool `short:"i" description:"Include multiple hardlinks to the same inode as duplicates"`

	Positional struct {
		Directory string
	} `positional-args:"yes" required:"true"`
}

func debugLog(s string, a ...interface{}) {
	if opts.Verbose {
		log.Printf(s, a...)
	}
}

func pluralize(count int, name string) string {
	return fmt.Sprintf("%v %v%v", count, name, If(count > 1, "s", ""))
}

func getDisplayLines(duplicateGroups [][]file_info.FileInfo) <-chan string {
	baseTemplate := opts.BaseTemplate
	outTemplate := opts.OutTemplate
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

func main() {
	_, err := flags.Parse(&opts)
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
	for line := range getDisplayLines(duplicateGroups) {
		fmt.Println(line)
	}
}
