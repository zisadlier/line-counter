package types

import (
	"fmt"
)

type OptionsValidationError struct {
	FieldName string
	Reason string
}

func (o *OptionsValidationError) Error() string {
	return fmt.Sprintf("Invalid option value for %s. %s", o.FieldName, o.Reason)
}

type OptionsValidationErrorGroup struct {
	Errors []*OptionsValidationError
}

func (o *OptionsValidationErrorGroup) Error() string {
	errStr := ""
	for _, err := range o.Errors {
		errStr += fmt.Sprintf("%s\n", err)
	}

	return errStr
}