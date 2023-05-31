package types

import (
	"testing"
)

type TestCase struct {
	input string
	result bool
}

func TestOptions(t *testing.T) {
	testOptions := &Options{
		Extensions: []string { ".txt", ".go" },
		Directory: ".",
		Verbose: false,
		CountWhiteSpace: false,
		SkipRegexStr: "[a-z]{3,}_[A-Z]{2}",
		IncludeRegexStr: "(foo)|(bar)",
	}

	if errs := testOptions.Validate(); len(errs) > 0 {
		t.Fatal(errs)
	}

	// Test that RegEx compiled correctly
	skipCases := []TestCase {
		{ "abcc_FF", true },
		{ "ab_F", false },
		{ "a_f_5", false },
		{ "tew_JSJD", true },
	}

	for _, tc := range skipCases {
		result := testOptions.SkipRegex.MatchString(tc.input)
		if result != tc.result {
			t.Fatalf("%s returned result %t, but expected %t", tc.input, result, tc.result)
		}
	}

	includeCases := []TestCase {
		{ "foo", true },
		{ "bar", true },
		{ "test", false },
		{ "asdnasdnasd", false },
	}

	for _, tc := range includeCases {
		result := testOptions.IncludeRegex.MatchString(tc.input)
		if result != tc.result {
			t.Fatalf("%s returned result %t, but expected %t", tc.input, result, tc.result)
		}
	}

	// Test extension related methods
	if testOptions.IsExtensionsWildcard() {
		t.Fatal("Extensions should not be wildcard")
	}
	
	processCases := []TestCase {
		{ "text.txt", true },
		{ "text", false },
		{ "runner.go", true },
		{ "runnergo", false },
		{ "test.txtt", false },
	}

	for _, tc := range processCases {
		result := testOptions.ShouldFileBeProcessed(tc.input)
		if result != tc.result {
			t.Fatalf("%s returned result %t, but expected %t", tc.input, result, tc.result)
		}
	}
}