package envfile

import (
	"fmt"
	"strings"
)

// FormatSchemaViolations returns a human-readable summary of schema violations.
func FormatSchemaViolations(violations []SchemaViolation) string {
	if len(violations) == 0 {
		return "schema validation passed: no violations found\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "schema validation failed: %d violation(s) found\n", len(violations))
	for _, v := range violations {
		fmt.Fprintf(&sb, "  [%s] %s\n", v.Key, v.Message)
	}
	return sb.String()
}

// FormatSchemaFields returns a human-readable summary of the schema definition.
func FormatSchemaFields(schema Schema) string {
	if len(schema.Fields) == 0 {
		return "schema: (empty)\n"
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "schema: %d field(s) defined\n", len(schema.Fields))
	for _, f := range schema.Fields {
		required := "optional"
		if f.Required {
			required = "required"
		}
		pattern := f.Pattern
		if pattern == "" {
			pattern = "any"
		}
		fmt.Fprintf(&sb, "  %-30s %-10s type=%s\n", f.Key, required, pattern)
	}
	return sb.String()
}
