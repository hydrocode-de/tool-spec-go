package validate

import (
	"fmt"
	"math"
	"slices"
	"time"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

func ValidateParameters(spec toolspec.ToolSpec, inputs map[string]interface{}, failOnExtra bool) (bool, []error) {
	var errors []error = make([]error, 0)

	// check for everything in the inputs map
	for name, value := range inputs {
		paramSpec, ok := spec.Parameters[name]
		if !ok && failOnExtra {
			allowedNames := make([]string, 0, len(spec.Parameters))
			for name := range spec.Parameters {
				allowedNames = append(allowedNames, name)
			}
			errors = append(errors, &ValidationError{
				Field:    AllowedField(Parameters),
				Name:     name,
				Type:     ErrorType(NotAllowed),
				Expected: fmt.Sprintf("one of %v", allowedNames),
				Actual:   name,
				Message:  fmt.Sprintf("parameter %s is not allowed, allowed parameters are: %v", name, allowedNames),
			})
		}

		if err := ValidateParameter(paramSpec, value); err != nil {
			errors = append(errors, err)
		}
	}

	for name, paramSpec := range spec.Parameters {
		if _, ok := inputs[name]; !ok {
			if !paramSpec.Optional && paramSpec.Default == nil {
				errors = append(errors, &ValidationError{
					Field:    AllowedField(Parameters),
					Name:     name,
					Type:     ErrorType(Required),
					Expected: "not null",
					Actual:   "null",
					Message:  fmt.Sprintf("%s is a required parameter but was not provided", name),
				})
			}
		}
	}

	if len(errors) > 0 {
		return true, errors
	}
	return false, nil
}

func ValidateParameter(spec toolspec.ParameterSpec, value interface{}) error {
	// in case the value is marked as an array, we first check and then validate recursively
	if spec.IsArray {
		if _, ok := value.([]interface{}); !ok {
			return &ValidationError{
				Field:    AllowedField(Parameters),
				Name:     spec.Name,
				Type:     ErrorType(NotArray),
				Expected: fmt.Sprintf("[]%s", spec.ToolType),
				Actual:   fmt.Sprintf("%T", value),
				Message:  fmt.Sprintf("expected %s to be an array of %s", spec.Name, spec.ToolType),
			}
		}
		// Create a modified spec with IsArray=false for element validation
		elementSpec := spec
		elementSpec.IsArray = false
		for _, v := range value.([]interface{}) {
			if err := ValidateParameter(elementSpec, v); err != nil {
				return err
			}
		}
		return nil
	}

	// check if the parameter is an empty interface
	if value == nil {
		if spec.Optional {
			return nil
		} else {
			return &ValidationError{
				Field:    AllowedField(Parameters),
				Name:     spec.Name,
				Type:     ErrorType(Required),
				Expected: "not nil",
				Actual:   "nil",
				Message:  fmt.Sprintf("%s is required", spec.Name),
			}
		}
	}

	// Normalize float64 to int when spec expects integer and value is a whole number
	// JSON unmarshals all numbers as float64, so we need to convert them early
	normalizedValue := value
	if spec.ToolType == "integer" {
		if floatVal, ok := value.(float64); ok {
			if floatVal == math.Trunc(floatVal) {
				// Convert to int for all subsequent checks
				normalizedValue = int(floatVal)
			} else {
				return &ValidationError{
					Field:    AllowedField(Parameters),
					Name:     spec.Name,
					Type:     ErrorType(WrongType),
					Expected: "integer",
					Actual:   "float",
					Message:  fmt.Sprintf("expected %s to be an integer", spec.Name),
				}
			}
		}
	}

	// Check the type of normalized value
	valueType := getValueType(normalizedValue)
	if valueType != typeMap[spec.ToolType] {
		return &ValidationError{
			Field:    AllowedField(Parameters),
			Name:     spec.Name,
			Type:     ErrorType(WrongType),
			Expected: typeMap[spec.ToolType],
			Actual:   valueType,
			Message:  fmt.Sprintf("expected %s to be a %s", spec.Name, typeMap[spec.ToolType]),
		}
	}

	if spec.ToolType == "integer" || spec.ToolType == "float" {
		var rangeVal float64
		switch spec.ToolType {
		case "integer":
			rangeVal = float64(normalizedValue.(int))
		case "float":
			rangeVal = normalizedValue.(float64)
		}

		if spec.Min != nil && rangeVal < *spec.Min {
			return &ValidationError{
				Field:    AllowedField(Parameters),
				Name:     spec.Name,
				Type:     ErrorType(OutOfRange),
				Expected: fmt.Sprintf(">= %.2f", *spec.Min),
				Actual:   fmt.Sprintf("%.2f", rangeVal),
				Message:  fmt.Sprintf("%s must be >= %.2f", spec.Name, *spec.Min),
			}
		}

		if spec.Max != nil && rangeVal > *spec.Max {
			return &ValidationError{
				Field:    AllowedField(Parameters),
				Name:     spec.Name,
				Type:     ErrorType(OutOfRange),
				Expected: fmt.Sprintf("<= %.2f", *spec.Max),
				Actual:   fmt.Sprintf("%.2f", rangeVal),
				Message:  fmt.Sprintf("%s must be <= %.2f", spec.Name, *spec.Max),
			}
		}
	}

	if spec.ToolType == "enum" {
		if !slices.Contains(spec.Values, normalizedValue.(string)) {
			return &ValidationError{
				Field:    AllowedField(Parameters),
				Name:     spec.Name,
				Type:     ErrorType(NotInEnum),
				Expected: fmt.Sprintf("one of %v", spec.Values),
				Actual:   normalizedValue.(string),
				Message:  fmt.Sprintf("%s must be one of %v", spec.Name, spec.Values),
			}
		}
	}

	if spec.ToolType == "datetime" || spec.ToolType == "date" || spec.ToolType == "time" {
		if _, err := time.Parse(time.RFC3339, normalizedValue.(string)); err != nil {
			return &ValidationError{
				Field:    AllowedField(Parameters),
				Name:     spec.Name,
				Type:     ErrorType(InvalidDateTime),
				Expected: "a valid ISO 8601 datetime string",
				Actual:   normalizedValue.(string),
				Message:  fmt.Sprintf("%s must be a valid ISO 8601 datetime string", spec.Name),
			}
		}
	}

	return nil
}

var typeMap = map[string]string{
	"string":   "string",
	"asset":    "string",
	"enum":     "string",
	"integer":  "integer",
	"float":    "float",
	"boolean":  "boolean",
	"datetime": "datetime",
	"date":     "datetime",
	"time":     "datetime",
}

func getValueType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64:
		return "integer"
	case float32, float64:
		return "float"
	case bool:
		return "boolean"
	case time.Time:
		return "datetime"
	default:
		return "unknown"
	}
}
