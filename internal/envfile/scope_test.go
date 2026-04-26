package envfile

import (
	"strings"
	"testing"
)

func baseScope() map[string]string {
	return map[string]string{
		"APP_HOST":   "localhost",
		"APP_PORT":   "8080",
		"DB_URL":     "postgres://",
		"DB_PASS":    "secret",
		"LOG_LEVEL":  "info",
		"UNRELATED":  "value",
	}
}

func TestScope_SinglePrefix(t *testing.T) {
	r := Scope(baseScope(), ScopeOptions{Prefixes: []string{"APP_"}})
	if len(r.Scoped) != 2 {
		t.Fatalf("expected 2 scoped keys, got %d", len(r.Scoped))
	}
	if _, ok := r.Scoped["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in scoped")
	}
	if len(r.Excluded) != 4 {
		t.Fatalf("expected 4 excluded keys, got %d", len(r.Excluded))
	}
}

func TestScope_MultiplePrefix(t *testing.T) {
	r := Scope(baseScope(), ScopeOptions{Prefixes: []string{"APP_", "DB_"}})
	if len(r.Scoped) != 4 {
		t.Fatalf("expected 4 scoped keys, got %d", len(r.Scoped))
	}
}

func TestScope_StripPrefix(t *testing.T) {
	r := Scope(baseScope(), ScopeOptions{
		Prefixes:    []string{"APP_"},
		StripPrefix: true,
	})
	if _, ok := r.Scoped["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_ prefix")
	}
	if _, ok := r.Scoped["PORT"]; !ok {
		t.Error("expected PORT after stripping APP_ prefix")
	}
}

func TestScope_Uppercase(t *testing.T) {
	env := map[string]string{"app_host": "localhost", "db_url": "pg"}
	r := Scope(env, ScopeOptions{
		Prefixes:    []string{"app_"},
		StripPrefix: true,
		Uppercase:   true,
	})
	if _, ok := r.Scoped["HOST"]; !ok {
		t.Error("expected uppercase HOST key")
	}
}

func TestScope_NoPrefixes_IncludesAll(t *testing.T) {
	r := Scope(baseScope(), ScopeOptions{})
	if len(r.Scoped) != len(baseScope()) {
		t.Fatalf("expected all keys included when no prefix filter")
	}
	if len(r.Excluded) != 0 {
		t.Fatalf("expected no excluded keys")
	}
}

func TestScope_DoesNotMutateInput(t *testing.T) {
	env := baseScope()
	origLen := len(env)
	Scope(env, ScopeOptions{Prefixes: []string{"APP_"}, StripPrefix: true})
	if len(env) != origLen {
		t.Error("Scope mutated the input map")
	}
}

func TestScope_EmptyEnv(t *testing.T) {
	r := Scope(map[string]string{}, ScopeOptions{Prefixes: []string{"APP_"}})
	if len(r.Scoped) != 0 {
		t.Error("expected empty scoped map")
	}
}

func TestFormatScopeResult_ContainsLabels(t *testing.T) {
	r := Scope(baseScope(), ScopeOptions{Prefixes: []string{"APP_"}})
	out := FormatScopeResult(r)
	if !strings.Contains(out, "included") {
		t.Error("expected 'included' in format output")
	}
	if !strings.Contains(out, "excluded") {
		t.Error("expected 'excluded' in format output")
	}
}
