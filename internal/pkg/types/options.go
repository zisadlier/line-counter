package types

import (
	"os"
	"regexp"
	"strings"
	"golang.org/x/exp/slices"
)

var DefaultDirectory = "."
var Wildcard = "*"

type Options struct {
	// List of valid file extensions to count lines of
	Extensions []string
	// Directory to count lines in
	Directory string
	// File to count the lines of
	File string
	// If true, prints out extra information
	Verbose bool
	// If true, lines that are purely whitespace are also counted
	CountWhiteSpace bool
	// If true, tells the runner to just print the current version
	Version bool
	// Regex that will cause a line to be skipped if it matches it
	SkipRegexStr string
	SkipRegex *regexp.Regexp
	// Regex for which only lines that match it will be counted
	IncludeRegexStr string
	IncludeRegex *regexp.Regexp
}

func (o *Options) Validate() (errors []*OptionsValidationError) {
	filePresent := o.HasFile()
	directoryPresent := o.HasDirectory()

	if filePresent && directoryPresent {
		errors = append(errors, &OptionsValidationError{
			FieldName: "Directory & File",
			Reason: "Directory and file cannot be used together",
		})
	}

	if filePresent {
		if _, err := os.Stat(o.File); os.IsNotExist(err) {
			errors = append(errors, &OptionsValidationError{
				FieldName: "File",
				Reason: "Please pass in a valid file",
			})
		}
	}

	if directoryPresent {
		if _, err := os.Stat(o.Directory); os.IsNotExist(err) {
			errors = append(errors, &OptionsValidationError{
				FieldName: "Directory",
				Reason: "Please pass in a valid directory path",
			})
		}
	}

	skipRegex, err := regexp.Compile(o.SkipRegexStr)
	if err == nil && len(strings.TrimSpace(o.SkipRegexStr)) > 0 {
		o.SkipRegex = skipRegex
	} else if err != nil {
		errors = append(errors, &OptionsValidationError{
			FieldName: "Skip Regex",
			Reason: "Please pass in a valid regular expression",
		})
	}

	includeRegex, err := regexp.Compile(o.IncludeRegexStr)
	if err == nil && len(strings.TrimSpace(o.IncludeRegexStr)) > 0 {
		o.IncludeRegex = includeRegex
	} else if err != nil {
		errors = append(errors, &OptionsValidationError{
			FieldName: "Include Regex",
			Reason: "Please pass in a valid regular expression",
		})
	}

	return errors
}

func (o *Options) HasDirectory() bool {
	return o.Directory != "" && o.Directory != DefaultDirectory
}

func (o *Options) HasFile() bool {
	return o.File != ""
}

func (o *Options) ShouldFileBeProcessed(file string) bool {
	return slices.IndexFunc(o.Extensions, func(e string) bool { return strings.HasSuffix(file, e) }) != -1 || o.IsExtensionsWildcard()
}

func (o *Options) IsExtensionsWildcard() bool {
	return len(o.Extensions) == 1 && o.Extensions[0] == Wildcard
}