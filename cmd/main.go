package main

import (
	"flag"
	"strings"
	"os"
	"log"
	"fmt"
	"github.com/zisadlier/linecount/internal/runner"
	"github.com/zisadlier/linecount/internal/pkg/types"
)

var help = flag.Bool("help", false, "Show help message")
var exts = types.Wildcard
var directory = types.DefaultDirectory
var file = ""
var verbose = false
var countWhitespace = false
var version = false
var skipRegex = ""
var includeRegex = ""

func main() {
	flag.StringVar(&exts, "e", types.Wildcard, "File extensions to count lines for (comma separated)")
	flag.StringVar(&directory, "d", types.DefaultDirectory, "Directory for scan for files to line count, defaults to current directory. Cannot be used with file")
	flag.StringVar(&file, "f", "", "File to count the lines of, cannot be used with directory")
	flag.BoolVar(&verbose, "v", false, "If true, prints linecounts for individual files")
	flag.BoolVar(&countWhitespace, "w", false, "If true, lines that are empty or whitespace are also counted")
	flag.BoolVar(&version, "V", false, "If true, displays the app version")
	flag.StringVar(&skipRegex, "s", "", "Lines matching this regex will not be included in the count")
	flag.StringVar(&includeRegex, "i", "", "Only lines matching this regex will be included in the count")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	extensions := strings.Split(exts, ",")
	options := &types.Options{
		Extensions: extensions,
		Directory: directory,
		File: file,
		Verbose: verbose,
		CountWhiteSpace: countWhitespace,
		Version: version,
		SkipRegexStr: skipRegex,
		IncludeRegexStr: includeRegex,
	}

	runner, err := runner.New(options)
	if err != nil {
		log.Fatal(err)
	}

	if runner == nil {
		// Version was printed
		os.Exit(0)
	}

	if err := runner.Run(); err != nil {
		log.Fatal(err)
	}

	if runner.Options.HasFile() {
		fmt.Printf("Total matching number of lines in %s: %d\n", runner.Options.File, runner.TotalLines)
	} else if runner.Options.HasDirectory() {
		if runner.Options.Verbose {
			for _, f := range runner.GetFiles() {
				fmt.Printf("%s: %d\n", f.Name, f.Lines)
			}
		}

		extStr := ""
		if !runner.Options.IsExtensionsWildcard() {
			extStr = fmt.Sprintf(" with extensions %s", runner.Options.Extensions)
		}

		fmt.Printf("Total matching number of lines for files in %s%s: %d\n", runner.Options.Directory, extStr, runner.TotalLines)
	}
}