package utils

import (
	"regexp"
)

// ShouldLineBeCounted determines whether a given line of text should be counted based on the provided regex patterns
func ShouldLineBeCounted(line string, skip *regexp.Regexp, include *regexp.Regexp) bool {
	result := true

	if skip != nil {
		result = result && !skip.MatchString(line)
	}

	if include != nil {
		result = result && include.MatchString(line)
	}

	return result
}