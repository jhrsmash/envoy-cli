package envfile

import (
	"strings"
	"testing"
)

func TestDedup_NoDuplicates(t *testing.T) {
	lines := []string{"FOO=bar", "BAZ=qux"}
	r := Dedup(lines)
	if len(r.Duplicates) != 0 {
		t.Fatalf("expected no duplicates, got %v", r.Duplicates)
	}
	if r.Output["FOO"] != "bar" || r.Output["BAZ"] != "qux" {
		t.Fatalf("unexpected output: %v", r.Output)
	}
}

func TestDedup_DetectsDuplicate(t *testing.T) {
	lines := []string{"FOO=first", "BAR=keep", "FOO=second", "FOO=third"}
	r := Dedup(lines)
	if r.Output["FOO"] != "third" {
		t.Errorf("expected last value 'third', got %q", r.Output["FOO"])
	}
	if len(r.Duplicates["FOO"]) != 3 {
		t.Errorf("expected 3 recorded values, got %d", len(r.Duplicates["FOO"]))
	}
	if _, dup := r.Duplicates["BAR"]; dup {
		t.Error("BAR should not appear in duplicates")
	}
}

func TestDedup_SkipsCommentsAndBlanks(t *testing.T) {
	lines := []string{"# comment", "", "  ", "KEY=val"}
	r := Dedup(lines)
	if len(r.Output) != 1 {
		t.Fatalf("expected 1 key, got %d", len(r.Output))
	}
	if r.Output["KEY"] != "val" {
		t.Errorf("unexpected value: %q", r.Output["KEY"])
	}
}

func TestDedup_StripsQuotes(t *testing.T) {
	lines := []string{`SINGLE='hello'`, `DOUBLE="world"`}
	r := Dedup(lines)
	if r.Output["SINGLE"] != "hello" {
		t.Errorf("expected 'hello', got %q", r.Output["SINGLE"])
	}
	if r.Output["DOUBLE"] != "world" {
		t.Errorf("expected 'world', got %q", r.Output["DOUBLE"])
	}
}

func TestDedup_SkipsInvalidLines(t *testing.T) {
	lines := []string{"NOTAKVPAIR", "VALID=yes"}
	r := Dedup(lines)
	if _, ok := r.Output["NOTAKVPAIR"]; ok {
		t.Error("invalid line should not appear in output")
	}
	if r.Output["VALID"] != "yes" {
		t.Errorf("expected 'yes', got %q", r.Output["VALID"])
	}
}

func TestFormatDedupResult_NoDuplicates(t *testing.T) {
	r := DedupResult{Output: map[string]string{"A": "1"}, Duplicates: map[string][]string{}}
	out := FormatDedupResult(r)
	if !strings.Contains(out, "No duplicate") {
		t.Errorf("expected no-duplicate message, got: %q", out)
	}
}

func TestFormatDedupResult_WithDuplicates(t *testing.T) {
	r := DedupResult{
		Output: map[string]string{"FOO": "c"},
		Duplicates: map[string][]string{"FOO": {"a", "b", "c"}},
	}
	out := FormatDedupResult(r)
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %q", out)
	}
	if !strings.Contains(out, "3 occurrences") {
		t.Errorf("expected occurrence count, got: %q", out)
	}
	if !strings.Contains(out, `"c"`) {
		t.Errorf("expected kept value 'c', got: %q", out)
	}
}
