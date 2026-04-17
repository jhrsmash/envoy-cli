package envfile

import (
	"regexp"
	"testing"
)

func TestRedact_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
	}
	result := Redact(env, RedactOptions{})

	if result["DB_PASSWORD"] != redactedValue {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != redactedValue {
		t.Errorf("expected API_KEY to be redacted, got %q", result["API_KEY"])
	}
	if result["APP_NAME"] != "myapp" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result["APP_NAME"])
	}
}

func TestRedact_ExplicitKeys(t *testing.T) {
	env := map[string]string{
		"MY_CUSTOM_FIELD": "sensitive",
		"SAFE":            "visible",
	}
	opts := RedactOptions{Keys: []string{"MY_CUSTOM_FIELD"}}
	result := Redact(env, opts)

	if result["MY_CUSTOM_FIELD"] != redactedValue {
		t.Errorf("expected MY_CUSTOM_FIELD to be redacted")
	}
	if result["SAFE"] != "visible" {
		t.Errorf("expected SAFE to remain visible")
	}
}

func TestRedact_ExtraPatterns(t *testing.T) {
	env := map[string]string{
		"STRIPE_WEBHOOK": "whsec_xyz",
		"LOG_LEVEL":      "info",
	}
	opts := RedactOptions{
		ExtraPatterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)webhook`),
		},
	}
	result := Redact(env, opts)

	if result["STRIPE_WEBHOOK"] != redactedValue {
		t.Errorf("expected STRIPE_WEBHOOK to be redacted")
	}
	if result["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL to remain unchanged")
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"SECRET_KEY": "original",
	}
	_ = Redact(env, RedactOptions{})
	if env["SECRET_KEY"] != "original" {
		t.Errorf("Redact must not mutate the input map")
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	result := Redact(map[string]string{}, RedactOptions{})
	if len(result) != 0 {
		t.Errorf("expected empty result for empty input")
	}
}

func TestRedact_ExplicitKeysCaseInsensitive(t *testing.T) {
	env := map[string]string{
		"my_token": "supersecret",
		"VISIBLE":  "yes",
	}
	opts := RedactOptions{Keys: []string{"MY_TOKEN"}}
	result := Redact(env, opts)

	// Explicit keys should match case-insensitively so that callers don't
	// need to know the exact casing used in the env map.
	if result["my_token"] != redactedValue {
		t.Errorf("expected my_token to be redacted by case-insensitive key match, got %q", result["my_token"])
	}
	if result["VISIBLE"] != "yes" {
		t.Errorf("expected VISIBLE to remain unchanged")
	}
}
