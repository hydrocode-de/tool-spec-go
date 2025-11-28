package validate

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	toolspec "github.com/hydrocode-de/tool-spec-go"
)

func ValidateData(spec toolspec.ToolSpec, data map[string]string) (bool, []*ValidationError) {
	var errors []*ValidationError = make([]*ValidationError, 0)

	for name, dataSpec := range spec.Data {
		if _, ok := data[name]; !ok {
			errors = append(errors, &ValidationError{
				Field:    AllowedField(Data),
				Name:     name,
				Type:     ErrorType(Required),
				Expected: "not null",
				Actual:   "null",
				Message:  fmt.Sprintf("%s is a required data entry but was not provided", name),
			})
		}

		if dataSpec.Extensions != nil {
			lowerExts := make([]string, 0, len(dataSpec.Extensions))
			for _, ext := range dataSpec.Extensions {
				ext := strings.ToLower(ext)
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				lowerExts = append(lowerExts, ext)
			}

			dataExt := strings.ToLower(filepath.Ext(data[name]))
			if !slices.Contains(lowerExts, dataExt) {
				errors = append(errors, &ValidationError{
					Field:    AllowedField(Data),
					Name:     name,
					Type:     ErrorType(WrongType),
					Expected: fmt.Sprintf("one of %v", lowerExts),
					Actual:   dataExt,
					Message:  fmt.Sprintf("data file %s has an invalid extension, expected one of %v", name, lowerExts),
				})
			}
		}
	}

	if len(errors) > 0 {
		return true, errors
	}

	return false, nil
}
