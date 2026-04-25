package envfile

import (
	"strings"
	"testing"
)

func TestNormalize_NoOp(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	r := Normalize(env, NormalizeOptions{})
	if len(r.Changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(r.Changes))
	}
	if r.Env["KEY"] != "value" {
		t.Errorf("unexpected value: %s", r.Env["KEY"])
	}
}

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "DB_PORT": "5432"}
	r := Normalize(env, NormalizeOptions{UppercaseKeys: true})
	if _, ok := r.Env["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := r.Env["db_host"]; ok {
		t.Error("old lowercase key should be gone")
	}
	// DB_PORT was already uppercase — no change recorded for it
	for _, c := range r.Changes {
		if c.Key == "DB_PORT" {
			t.Error("DB_PORT should not appear in changes")
		}
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	env := map[string]string{"HOST": "  localhost  ", "PORT": "5432"}
	r := Normalize(env, NormalizeOptions{TrimValues: true})
	if r.Env["HOST"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", r.Env["HOST"])
	}
	if len(r.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(r.Changes))
	}
}

func TestNormalize_RemoveEmptyValues(t *testing.T) {
	env := map[string]string{"KEY": "value", "EMPTY": ""}
	r := Normalize(env, NormalizeOptions{TrimValues: true, RemoveEmptyValues: true})
	if _, ok := r.Env["EMPTY"]; ok {
		t.Error("EMPTY key should have been removed")
	}
	if r.Env["KEY"] != "value" {
		t.Error("KEY should remain unchanged")
	}
}

func TestNormalize_QuoteValues(t *testing.T) {
	env := map[string]string{"GREETING": "hello world", "NAME": "Alice"}
	r := Normalize(env, NormalizeOptions{QuoteValues: true})
	if r.Env["GREETING"] != `"hello world"` {
		t.Errorf("expected quoted value, got %s", r.Env["GREETING"])
	}
	if r.Env["NAME"] != "Alice" {
		t.Error("NAME has no spaces and should not be quoted")
	}
}

func TestNormalize_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"key": "  val  "}
	_ = Normalize(env, NormalizeOptions{UppercaseKeys: true, TrimValues: true})
	if _, ok := env["key"]; !ok {
		t.Error("original map was mutated")
	}
	if env["key"] != "  val  " {
		t.Error("original value was mutated")
	}
}

func TestFormatNormalizeResult_NoChanges(t *testing.T) {
	r := NormalizeResult{Env: map[string]string{"A": "1"}, Changes: nil}
	out := FormatNormalizeResult(r)
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected 'No changes' message, got: %s", out)
	}
}

func TestFormatNormalizeResult_WithChanges(t *testing.T) {
	r := NormalizeResult{
		Env: map[string]string{"DB_HOST": "localhost"},
		Changes: []NormalizeChange{
			{Key: "DB_HOST", OldKey: "db_host", OldVal: "localhost", NewVal: "localhost", Reason: "key uppercased"},
		},
	}
	out := FormatNormalizeResult(r)
	if !strings.Contains(out, "RENAMED") {
		t.Errorf("expected RENAMED in output, got: %s", out)
	}
	if !strings.Contains(out, "db_host") {
		t.Errorf("expected old key in output, got: %s", out)
	}
}
