package validate

import (
	"fmt"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

type AllowedField string

const (
	Parameters AllowedField = "parameters"
	Data       AllowedField = "data"
)

type ErrorType string

const (
	WrongType       ErrorType = "wrong-type"
	NotArray        ErrorType = "not-array"
	OutOfRange      ErrorType = "out-of-range"
	NotInEnum       ErrorType = "not-in-enum"
	InvalidDateTime ErrorType = "invalid-datetime"
	Required        ErrorType = "required"
	NotAllowed      ErrorType = "not-allowed"
)

type ValidationError struct {
	Field    AllowedField `json:"field"`
	Name     string       `json:"name"`
	Type     ErrorType    `json:"type"`
	Expected string       `json:"expected"`
	Actual   string       `json:"actual"`
	Message  string       `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s (expected %s, got %s)", e.Name, e.Message, e.Expected, e.Actual)
}

func ValidateInputs(spec toolspec.ToolSpec, inputs toolspec.ToolInput) (bool, []*ValidationError) {
	var errors []*ValidationError = make([]*ValidationError, 0)

	didError, errs := ValidateParameters(spec, inputs.Parameters, false)
	if didError {
		errors = append(errors, errs...)
	}
	didError, errs = ValidateData(spec, inputs.Datasets)
	if didError {
		errors = append(errors, errs...)
	}

	if len(errors) > 0 {
		return true, errors
	}
	return false, nil
}
