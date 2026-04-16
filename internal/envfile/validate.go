package envfile

import (
	"fmt"
	"strings"
)

// ValidationError represents a single validation issue found in an env map.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

// Valid returns true if no validation errors were found.
func (r ValidationResult) Valid() bool {
	return len(r.Errors) == 0
}

// Validate checks an env map for common issues:
//   - Empty values
//   - Keys containing spaces or invalid characters
//   - Duplicate keys are not possible in a map, but we check for blank keys
func Validate(env map[string]string) ValidationResult {
	result := ValidationResult{}

	for key, value := range env {
		if key == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "key must not be empty",
			})
			continue
		}

		if strings.ContainsAny(key, " \t") {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "key must not contain spaces or tabs",
			})
		}

		if strings.HasPrefix(key, "=") {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "key must not start with '='",
			})
		}

		if value == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     key,
				Message: "value is empty",
			})
		}
	}

	return result
}
