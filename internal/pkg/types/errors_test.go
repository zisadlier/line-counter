package types

import (
	"testing"
)

func TestOptionValidationError(t *testing.T) {
	testValidationError := &OptionsValidationError{
		FieldName: "Directory",
		Reason: "Directory does not exist",
	}

	error := testValidationError.Error()
	correctError := "Invalid option value for Directory. Directory does not exist"
	if error != correctError {
		t.Fatalf("Error string incorrect, should be %s but got %s", correctError, error)
	}

	testValidationError2 := &OptionsValidationError{
		FieldName: "SkipRegex",
		Reason: "Invalid RegEx",
	}
	testValidationError3 := &OptionsValidationError{
		FieldName: "SkipRegex",
		Reason: "Invalid RegEx",
	}

	errorGroup := &OptionsValidationErrorGroup{
		Errors: []*OptionsValidationError{testValidationError, testValidationError2, testValidationError3},
	}

	errors := errorGroup.Error()
	correctErrors := "Invalid option value for Directory. Directory does not exist\n" +
					  "Invalid option value for SkipRegex. Invalid RegEx\n" +
					  "Invalid option value for SkipRegex. Invalid RegEx\n"

	if errors != correctErrors {
		t.Fatalf("Error string incorrect, should be %s but got %s", correctErrors, errors)
	}
}