# Missing Validation Logic

This document describes the validation logic that needs to be implemented for the `toolspec` package. This validation is required for both the `gorun` application and future container shim tools.

## Parameter Validation

### Type Validation
Validate that parameter values match their declared types:
- `string`: Must be a string value
- `integer`: Must be an integer number (no decimal part)
- `float`: Must be a numeric value (can have decimal part)
- `boolean`: Must be a boolean value (true/false)
- `datetime`: Must be a valid datetime string (ISO 8601 format recommended)
- `enum`: Must be one of the values in the `Values` slice

### Array Type Validation
When `IsArray` is `true`:
- Value must be an array/slice
- Each element in the array must pass type validation
- Empty arrays should be validated based on `Optional` flag

### Range Validation
For numeric types (`integer`, `float`):
- If `Min` is set, value must be >= `Min`
- If `Max` is set, value must be <= `Max`
- Both `Min` and `Max` can be set to create a range

### Enum Value Validation
When `ToolType` is `enum`:
- Value must be present in the `Values` slice
- Case-sensitive matching (or define case-insensitive option)

### Required Field Validation
- If `Optional` is `false` (or not set), the parameter must be provided
- If `Optional` is `true`, the parameter can be omitted
- When omitted, `Default` value should be used if provided

### Default Value Handling
- If a parameter is not provided and has a `Default` value, use the default
- Default values should also be validated against type and constraints

## Data File Validation

### Required Data Validation
- All data entries defined in `Data` map must be provided (unless optional flag is added)
- Each data entry must have a corresponding value in the inputs

### Path Validation
- Validate that data file paths exist (when running in execution context)
- Paths should be validated against the `Path` field in `DataSpec`

### Extension Validation
- If `Extensions` is defined, validate that the file extension matches one of the allowed extensions
- Extension matching should be case-insensitive
- File extension should be extracted from the file path

## Inputs.json Structure Validation

### Overall Structure
- Validate that inputs.json is valid JSON
- Validate that it's a map/dictionary structure
- Validate that tool names in the inputs file match expected tool names

### Tool Input Validation
- For each tool in the inputs file:
  - Validate that `Parameters` map matches the tool's `ParameterSpec` definitions
  - Validate that `Datasets` map matches the tool's `DataSpec` definitions
  - Run all parameter validations described above
  - Run all data validations described above

### Cross-Validation
- Ensure no extra parameters are provided that aren't defined in the spec
- Ensure no extra data entries are provided that aren't defined in the spec
- Validate that all required parameters are present
- Validate that all required data entries are present

## Proposed Function Signatures

```go
// Validate a single parameter value against its spec
func ValidateParameter(spec ParameterSpec, value interface{}) error

// Validate all parameters against a tool spec
func ValidateParameters(spec ToolSpec, inputs map[string]interface{}) error

// Validate all data entries against a tool spec
func ValidateData(spec ToolSpec, data map[string]string) error

// Validate a complete ToolInput against a ToolSpec
func ValidateInputs(spec ToolSpec, inputs ToolInput) error

// Validate an InputFile against a SpecFile
func ValidateInputFile(specFile SpecFile, inputFile InputFile) error
```

## Error Handling

Validation functions should return descriptive errors that indicate:
- Which parameter/data entry failed validation
- What type of validation failed (type, range, required, etc.)
- What the expected value/constraint was
- What the actual value was (if safe to include)

Consider using structured errors or error wrapping to allow callers to programmatically handle different validation failures.

