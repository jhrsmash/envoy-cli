package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateSchema_AllPresent(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "https://db.example.com",
		"PORT":         "8080",
		"DEBUG":        "true",
	}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "DATABASE_URL", Required: true, Pattern: "url"},
			{Key: "PORT", Required: true, Pattern: "int"},
			{Key: "DEBUG", Required: false, Pattern: "bool"},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	env := map[string]string{"PORT": "3000"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "DATABASE_URL", Required: true, Pattern: "url"},
			{Key: "PORT", Required: true, Pattern: "int"},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "DATABASE_URL" {
		t.Errorf("expected violation for DATABASE_URL, got %q", violations[0].Key)
	}
}

func TestValidateSchema_BadInt(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "PORT", Required: true, Pattern: "int"},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_BadBool(t *testing.T) {
	env := map[string]string{"DEBUG": "maybe"}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "DEBUG", Required: false, Pattern: "bool"},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestValidateSchema_OptionalMissing(t *testing.T) {
	env := map[string]string{}
	schema := Schema{
		Fields: []SchemaField{
			{Key: "OPTIONAL_KEY", Required: false, Pattern: "string"},
		},
	}
	violations := ValidateSchema(env, schema)
	if len(violations) != 0 {
		t.Fatalf("expected no violations for optional missing key, got %d", len(violations))
	}
}

func TestLoadSchema_ParsesFile(t *testing.T) {
	content := "# comment\nDATABASE_URL required url\nPORT optional int\nSECRET required\n"
	tmp := filepath.Join(t.TempDir(), "schema.txt")
	if err := os.WriteFile(tmp, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	schema, err := LoadSchema(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(schema.Fields))
	}
	if !schema.Fields[0].Required {
		t.Error("DATABASE_URL should be required")
	}
	if schema.Fields[0].Pattern != "url" {
		t.Errorf("expected pattern 'url', got %q", schema.Fields[0].Pattern)
	}
}

func TestFormatSchemaViolations_NoViolations(t *testing.T) {
	out := FormatSchemaViolations(nil)
	if out == "" {
		t.Error("expected non-empty output")
	}
	if !containsStr(out, "passed") {
		t.Errorf("expected 'passed' in output, got: %q", out)
	}
}

func TestFormatSchemaViolations_WithViolations(t *testing.T) {
	v := []SchemaViolation{{Key: "FOO", Message: "required key is missing or empty"}}
	out := FormatSchemaViolations(v)
	if !containsStr(out, "FOO") {
		t.Errorf("expected key FOO in output, got: %q", out)
	}
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && (
		func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		}()))
}
