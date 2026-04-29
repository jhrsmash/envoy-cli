package envfile

import (
	"strings"
	"testing"
)

func TestCast_Bool(t *testing.T) {
	env := map[string]string{
		"ENABLED": "TRUE",
		"DEBUG":   "1",
		"VERBOSE": "false",
	}
	r, err := Cast(env, CastOptions{TargetType: "bool"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Output["ENABLED"] != "true" {
		t.Errorf("ENABLED: got %q, want %q", r.Output["ENABLED"], "true")
	}
	if r.Output["DEBUG"] != "true" {
		t.Errorf("DEBUG: got %q, want %q", r.Output["DEBUG"], "true")
	}
	if r.Output["VERBOSE"] != "false" {
		t.Errorf("VERBOSE: got %q, want %q", r.Output["VERBOSE"], "false")
	}
	if len(r.Cast) != 3 {
		t.Errorf("expected 3 cast keys, got %d", len(r.Cast))
	}
}

func TestCast_Int(t *testing.T) {
	env := map[string]string{
		"PORT":    "8080",
		"TIMEOUT": "30.0",
	}
	r, err := Cast(env, CastOptions{TargetType: "int"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Output["PORT"] != "8080" {
		t.Errorf("PORT: got %q, want \"8080\"", r.Output["PORT"])
	}
	if r.Output["TIMEOUT"] != "30" {
		t.Errorf("TIMEOUT: got %q, want \"30\"", r.Output["TIMEOUT"])
	}
}

func TestCast_Float(t *testing.T) {
	env := map[string]string{"RATE": "1.5"}
	r, err := Cast(env, CastOptions{TargetType: "float"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Output["RATE"] != "1.5" {
		t.Errorf("RATE: got %q, want \"1.5\"", r.Output["RATE"])
	}
}

func TestCast_SelectedKeys(t *testing.T) {
	env := map[string]string{
		"A": "TRUE",
		"B": "FALSE",
	}
	r, err := Cast(env, CastOptions{TargetType: "bool", Keys: []string{"A"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Cast) != 1 || r.Cast[0] != "A" {
		t.Errorf("expected only A cast, got %v", r.Cast)
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "B" {
		t.Errorf("expected B skipped, got %v", r.Skipped)
	}
}

func TestCast_FailedNonStrict(t *testing.T) {
	env := map[string]string{"X": "not-a-bool"}
	r, err := Cast(env, CastOptions{TargetType: "bool"})
	if err != nil {
		t.Fatalf("unexpected error in non-strict mode: %v", err)
	}
	if len(r.Failed) != 1 || r.Failed[0] != "X" {
		t.Errorf("expected X in failed, got %v", r.Failed)
	}
	// original value preserved
	if r.Output["X"] != "not-a-bool" {
		t.Errorf("expected original value preserved, got %q", r.Output["X"])
	}
}

func TestCast_FailedStrict(t *testing.T) {
	env := map[string]string{"X": "not-an-int"}
	_, err := Cast(env, CastOptions{TargetType: "int", Strict: true})
	if err == nil {
		t.Fatal("expected error in strict mode")
	}
}

func TestCast_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"FLAG": "TRUE"}
	Cast(env, CastOptions{TargetType: "bool"}) //nolint
	if env["FLAG"] != "TRUE" {
		t.Errorf("input mutated: got %q", env["FLAG"])
	}
}

func TestFormatCastResult_NoChanges(t *testing.T) {
	r := CastResult{Output: map[string]string{}}
	out := FormatCastResult(r)
	if !strings.Contains(out, "no keys modified") {
		t.Errorf("expected 'no keys modified', got: %s", out)
	}
}

func TestFormatCastResult_ShowsCastAndFailed(t *testing.T) {
	r := CastResult{
		Output:  map[string]string{"A": "true", "B": "bad"},
		Cast:    []string{"A"},
		Failed:  []string{"B"},
	}
	out := FormatCastResult(r)
	if !strings.Contains(out, "cast") {
		t.Errorf("expected 'cast' in output, got: %s", out)
	}
	if !strings.Contains(out, "failed") {
		t.Errorf("expected 'failed' in output, got: %s", out)
	}
}
