package envfile

import (
	"strings"
	"testing"
)

func TestSanitize_NoOp(t *testing.T) {
	env := map[string]string{"KEY": "value", "OTHER": "hello"}
	r := Sanitize(env, SanitizeOptions{})
	if len(r.Changes) != 0 {
		t.Fatalf("expected no changes, got %d", len(r.Changes))
	}
	if r.Output["KEY"] != "value" {
		t.Errorf("unexpected value: %s", r.Output["KEY"])
	}
}

func TestSanitize_NormalizeLineEndings(t *testing.T) {
	env := map[string]string{"KEY": "line1\r\nline2", "CLEAN": "fine"}
	r := Sanitize(env, SanitizeOptions{NormalizeLineEndings: true})
	if r.Output["KEY"] != "line1\nline2" {
		t.Errorf("expected normalized value, got %q", r.Output["KEY"])
	}
	if len(r.Changes) != 1 || r.Changes[0].Key != "KEY" {
		t.Errorf("expected 1 change for KEY, got %v", r.Changes)
	}
	if r.Output["CLEAN"] != "fine" {
		t.Errorf("CLEAN should be unchanged")
	}
}

func TestSanitize_StripControlChars(t *testing.T) {
	env := map[string]string{"KEY": "val\x01ue\x00"}
	r := Sanitize(env, SanitizeOptions{StripControlChars: true})
	if r.Output["KEY"] != "value" {
		t.Errorf("expected control chars stripped, got %q", r.Output["KEY"])
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected 1 change")
	}
}

func TestSanitize_TrimQuotes(t *testing.T) {
	env := map[string]string{
		"DOUBLE": `"hello"`,
		"SINGLE": `'world'`,
		"PLAIN":  "plain",
	}
	r := Sanitize(env, SanitizeOptions{TrimQuotes: true})
	if r.Output["DOUBLE"] != "hello" {
		t.Errorf("expected double quotes trimmed, got %q", r.Output["DOUBLE"])
	}
	if r.Output["SINGLE"] != "world" {
		t.Errorf("expected single quotes trimmed, got %q", r.Output["SINGLE"])
	}
	if r.Output["PLAIN"] != "plain" {
		t.Errorf("PLAIN should be unchanged")
	}
	if len(r.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(r.Changes))
	}
}

func TestSanitize_MaxValueLength(t *testing.T) {
	env := map[string]string{"KEY": "abcdefghij", "SHORT": "hi"}
	r := Sanitize(env, SanitizeOptions{MaxValueLength: 5})
	if r.Output["KEY"] != "abcde" {
		t.Errorf("expected truncated value, got %q", r.Output["KEY"])
	}
	if r.Output["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged")
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected 1 change")
	}
}

func TestSanitize_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": `"quoted"`}
	Sanitize(env, SanitizeOptions{TrimQuotes: true})
	if env["KEY"] != `"quoted"` {
		t.Errorf("input map was mutated")
	}
}

func TestFormatSanitizeResult_NoChanges(t *testing.T) {
	r := SanitizeResult{Output: map[string]string{}, Changes: nil}
	out := FormatSanitizeResult(r)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes', got %q", out)
	}
}

func TestFormatSanitizeResult_WithChanges(t *testing.T) {
	r := SanitizeResult{
		Output: map[string]string{"KEY": "val"},
		Changes: []SanitizeChange{
			{Key: "KEY", Before: `"val"`, After: "val", Reason: "trimmed quotes"},
		},
	}
	out := FormatSanitizeResult(r)
	if !strings.Contains(out, "1 key(s) modified") {
		t.Errorf("expected modification count, got %q", out)
	}
	if !strings.Contains(out, "KEY") {
		t.Errorf("expected KEY in output")
	}
}
