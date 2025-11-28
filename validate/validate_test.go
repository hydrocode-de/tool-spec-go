package validate

import (
	"os"
	"path/filepath"
	"testing"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

func TestValidateValidConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	specPath := filepath.Join(wd, "..", "test", "data", "valid", "src", "tool.yml")
	inputsPath := filepath.Join(wd, "..", "test", "data", "valid", "in", "inputs.json")

	specData, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("failed to read spec file: %v", err)
	}

	specFile, err := toolspec.LoadToolSpec(specData)
	if err != nil {
		t.Fatalf("failed to parse spec file: %v", err)
	}

	toolSpec, err := specFile.GetTool("foobar")
	if err != nil {
		t.Fatalf("failed to get tool spec: %v", err)
	}

	inputsData, err := os.ReadFile(inputsPath)
	if err != nil {
		t.Fatalf("failed to read inputs file: %v", err)
	}

	inputFile, err := toolspec.LoadInputs(inputsData)
	if err != nil {
		t.Fatalf("failed to parse inputs file: %v", err)
	}

	toolInput, err := inputFile.GetToolInput("foobar")
	if err != nil {
		t.Fatalf("failed to get tool input: %v", err)
	}

	hasErrors, errors := ValidateInputs(toolSpec, toolInput)
	if hasErrors {
		t.Errorf("validation failed with %d errors:", len(errors))
		for _, err := range errors {
			t.Errorf("  - %v", err)
		}
	}
}

func TestValidateInvalidConfig(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	specPath := filepath.Join(wd, "..", "test", "data", "invalid", "src", "tool.yml")
	inputsPath := filepath.Join(wd, "..", "test", "data", "invalid", "in", "inputs.json")

	specData, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("failed to read spec file: %v", err)
	}

	specFile, err := toolspec.LoadToolSpec(specData)
	if err != nil {
		t.Fatalf("failed to parse spec file: %v", err)
	}

	toolSpec, err := specFile.GetTool("foobar")
	if err != nil {
		t.Fatalf("failed to get tool spec: %v", err)
	}

	inputsData, err := os.ReadFile(inputsPath)
	if err != nil {
		t.Fatalf("failed to read inputs file: %v", err)
	}

	inputFile, err := toolspec.LoadInputs(inputsData)
	if err != nil {
		t.Fatalf("failed to parse inputs file: %v", err)
	}

	toolInput, err := inputFile.GetToolInput("foobar")
	if err != nil {
		t.Fatalf("failed to get tool input: %v", err)
	}

	hasErrors, errors := ValidateInputs(toolSpec, toolInput)
	if !hasErrors {
		t.Fatal("expected validation to fail, but it passed")
	}

	if len(errors) == 0 {
		t.Fatal("expected validation errors, but got none")
	}

	// Check for specific error types
	errorTypes := make(map[ErrorType]bool)
	errorNames := make(map[string]bool)
	for _, err := range errors {
		errorTypes[err.Type] = true
		errorNames[err.Name] = true
	}

	// Verify we have multiple types of errors
	expectedErrorTypes := []ErrorType{
		WrongType,
		OutOfRange,
		NotInEnum,
		NotAllowed,
	}

	foundTypes := 0
	for _, expectedType := range expectedErrorTypes {
		if errorTypes[expectedType] {
			foundTypes++
		}
	}

	if foundTypes < 3 {
		t.Errorf("expected at least 3 different error types, found %d. Errors: %v", foundTypes, errors)
	}

	// Verify specific errors are present
	expectedErrors := map[string]bool{
		"foo_int":     false, // wrong type (float with decimals)
		"foo_float":   false, // out of range
		"foo_string":  false, // wrong type (number instead of string)
		"foo_enum":    false, // invalid enum value
		"foo_array":   false, // array element wrong type
		"foo_matrix":  false, // wrong extension
		"extra_param": false, // extra parameter
		"extra_data":  false, // extra data entry
	}

	for _, err := range errors {
		if _, exists := expectedErrors[err.Name]; exists {
			expectedErrors[err.Name] = true
		}
	}

	// Check that we found at least some of the expected errors
	foundErrors := 0
	for name, found := range expectedErrors {
		if found {
			foundErrors++
			t.Logf("Found expected error for: %s", name)
		}
	}

	if foundErrors < 5 {
		t.Errorf("expected at least 5 specific errors, found %d. All errors: %v", foundErrors, errors)
	}

	t.Logf("Validation correctly identified %d errors:", len(errors))
	for _, err := range errors {
		t.Logf("  - %v", err)
	}
}
