package envfile

import (
	"testing"
)

func TestValidate_Valid(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
		"DEBUG": "true",
	}
	result := Validate(env)
	if !result.Valid() {
		t.Errorf("expected no errors, got %d: %v", len(result.Errors), result.Errors)
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"HOST": "",
	}
	result := Validate(env)
	if result.Valid() {
		t.Error("expected validation error for empty value")
	}
	if result.Errors[0].Key != "HOST" {
		t.Errorf("expected error key HOST, got %q", result.Errors[0].Key)
	}
}

func TestValidate_KeyWithSpaces(t *testing.T) {
	env := map[string]string{
		"MY KEY": "value",
	}
	result := Validate(env)
	if result.Valid() {
		t.Error("expected validation error for key with space")
	}
}

func TestValidate_KeyWithTab(t *testing.T) {
	env := map[string]string{
		"MY\tKEY": "value",
	}
	result := Validate(env)
	if result.Valid() {
		t.Error("expected validation error for key with tab")
	}
}

func TestValidate_KeyStartsWithEquals(t *testing.T) {
	env := map[string]string{
		"=BAD": "value",
	}
	result := Validate(env)
	if result.Valid() {
		t.Error("expected validation error for key starting with '='")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	env := map[string]string{
		"GOOD":   "value",
		"EMPTY":  "",
		"BAD KEY": "x",
	}
	result := Validate(env)
	if len(result.Errors) != 2 {
		t.Errorf("expected 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationError_Error(t *testing.T) {
	e := ValidationError{Key: "FOO", Message: "value is empty"}
	got := e.Error()
	expected := `key "FOO": value is empty`
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
