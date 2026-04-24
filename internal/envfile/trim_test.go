package envfile

import (
	"strings"
	"testing"
)

func TestTrim_NoOp(t *testing.T) {
	env := map[string]string{"KEY": "value", "OTHER": "123"}
	r := Trim(env, TrimOptions{})

	if len(r.Output) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(r.Output))
	}
	if len(r.Trimmed) != 0 || len(r.Removed) != 0 {
		t.Error("expected no trimmed or removed keys")
	}
}

func TestTrim_TrimsWhitespace(t *testing.T) {
	env := map[string]string{
		"A": "  hello  ",
		"B": "world",
		"C": "\t spaced\t ",
	}
	r := Trim(env, TrimOptions{TrimValues: true})

	if r.Output["A"] != "hello" {
		t.Errorf("A: expected 'hello', got %q", r.Output["A"])
	}
	if r.Output["B"] != "world" {
		t.Errorf("B: expected 'world', got %q", r.Output["B"])
	}
	if r.Output["C"] != "spaced" {
		t.Errorf("C: expected 'spaced', got %q", r.Output["C"])
	}
	if len(r.Trimmed) != 2 {
		t.Errorf("expected 2 trimmed keys, got %d", len(r.Trimmed))
	}
}

func TestTrim_RemovesEmptyValues(t *testing.T) {
	env := map[string]string{
		"PRESENT": "ok",
		"EMPTY":   "",
		"BLANK":   "   ",
	}
	r := Trim(env, TrimOptions{TrimValues: true, RemoveEmpty: true})

	if _, ok := r.Output["EMPTY"]; ok {
		t.Error("EMPTY should have been removed")
	}
	if _, ok := r.Output["BLANK"]; ok {
		t.Error("BLANK should have been removed")
	}
	if r.Output["PRESENT"] != "ok" {
		t.Errorf("PRESENT: expected 'ok', got %q", r.Output["PRESENT"])
	}
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed keys, got %d", len(r.Removed))
	}
}

func TestTrim_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"X": "  val  "}
	_ = Trim(env, TrimOptions{TrimValues: true})
	if env["X"] != "  val  " {
		t.Error("Trim must not mutate the input map")
	}
}

func TestFormatTrimResult_NoChanges(t *testing.T) {
	r := TrimResult{Output: map[string]string{"A": "b"}}
	out := FormatTrimResult(r)
	if !strings.Contains(out, "nothing to do") {
		t.Errorf("expected 'nothing to do' message, got: %s", out)
	}
}

func TestFormatTrimResult_ShowsChanges(t *testing.T) {
	r := TrimResult{
		Output:  map[string]string{"A": "val"},
		Trimmed: []string{"A"},
		Removed: []string{"EMPTY"},
	}
	out := FormatTrimResult(r)
	if !strings.Contains(out, "trimmed") {
		t.Errorf("expected 'trimmed' in output, got: %s", out)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected 'removed' in output, got: %s", out)
	}
	if !strings.Contains(out, "EMPTY") {
		t.Errorf("expected 'EMPTY' in output, got: %s", out)
	}
}
