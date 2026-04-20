package envfile

import (
	"fmt"
	"strings"
)

// SchemaField describes a single expected key in an env file.
type SchemaField struct {
	Key      string
	Required bool
	Pattern  string // optional: "string", "int", "bool", "url"
}

// Schema is a collection of expected fields for an env file.
type Schema struct {
	Fields []SchemaField
}

// SchemaViolation represents a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

// ValidateSchema checks that the provided env map conforms to the given schema.
// It returns a list of violations (missing required keys, type mismatches, etc.).
func ValidateSchema(env map[string]string, schema Schema) []SchemaViolation {
	var violations []SchemaViolation

	for _, field := range schema.Fields {
		val, exists := env[field.Key]

		if !exists || val == "" {
			if field.Required {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: "required key is missing or empty",
				})
			}
			continue
		}

		if field.Pattern != "" {
			if msg := checkPattern(val, field.Pattern); msg != "" {
				violations = append(violations, SchemaViolation{
					Key:     field.Key,
					Message: fmt.Sprintf("expected %s: %s", field.Pattern, msg),
				})
			}
		}
	}

	return violations
}

// checkPattern validates a value against a named pattern type.
func checkPattern(val, pattern string) string {
	switch strings.ToLower(pattern) {
	case "bool":
		switch strings.ToLower(val) {
		case "true", "false", "1", "0", "yes", "no":
			return ""
		}
		return fmt.Sprintf("%q is not a valid boolean", val)
	case "int":
		for _, c := range val {
			if c < '0' || c > '9' {
				return fmt.Sprintf("%q is not a valid integer", val)
			}
		}
		return ""
	case "url":
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			return fmt.Sprintf("%q does not look like a URL", val)
		}
		return ""
	case "string":
		return ""
	}
	return ""
}
