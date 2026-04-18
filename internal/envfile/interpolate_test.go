package envfile

import (
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "5432",
	}
	out, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "localhost" || out["PORT"] != "5432" {
		t.Errorf("values should be unchanged, got %v", out)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"DSN": "postgres://${HOST}/db",
	}
	out, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DSN"] != "postgres://localhost/db" {
		t.Errorf("expected expanded DSN, got %q", out["DSN"])
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	env := map[string]string{
		"USER": "admin",
		"GREETING": "hello $USER",
	}
	out, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GREETING"] != "hello admin" {
		t.Errorf("expected 'hello admin', got %q", out["GREETING"])
	}
}

func TestInterpolate_MissingReference(t *testing.T) {
	env := map[string]string{
		"DSN": "postgres://${HOST}/db",
	}
	_, err := Interpolate(env)
	if err == nil {
		t.Fatal("expected error for missing reference, got nil")
	}
	ie, ok := err.(*InterpolateError)
	if !ok {
		t.Fatalf("expected *InterpolateError, got %T", err)
	}
	if ie.Missing != "HOST" {
		t.Errorf("expected missing key 'HOST', got %q", ie.Missing)
	}
}

func TestInterpolate_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"BASE": "http",
		"URL":  "${BASE}://example.com",
	}
	_, err := Interpolate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["URL"] != "${BASE}://example.com" {
		t.Errorf("input map was mutated")
	}
}
