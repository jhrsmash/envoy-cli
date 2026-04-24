package envfile

import (
	"strings"
	"testing"
)

var baseFilterEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"DB_HOST":      "db.local",
	"DB_PASSWORD":  "secret",
	"LOG_LEVEL":    "info",
	"FEATURE_FLAG": "true",
}

func TestFilter_ByPrefix(t *testing.T) {
	r, err := Filter(baseFilterEnv, FilterOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
	if _, ok := r.Matched["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in matched")
	}
	if _, ok := r.Matched["APP_PORT"]; !ok {
		t.Error("expected APP_PORT in matched")
	}
}

func TestFilter_BySuffix(t *testing.T) {
	r, err := Filter(baseFilterEnv, FilterOptions{Suffix: "_HOST"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestFilter_ByPattern(t *testing.T) {
	r, err := Filter(baseFilterEnv, FilterOptions{Pattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
}

func TestFilter_ByExplicitKeys(t *testing.T) {
	r, err := Filter(baseFilterEnv, FilterOptions{Keys: []string{"LOG_LEVEL", "FEATURE_FLAG"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(r.Matched))
	}
	if len(r.Dropped) != 4 {
		t.Errorf("expected 4 dropped, got %d", len(r.Dropped))
	}
}

func TestFilter_InvalidPattern(t *testing.T) {
	_, err := Filter(baseFilterEnv, FilterOptions{Pattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid pattern")
	}
}

func TestFilter_NoOptions_ReturnsAll(t *testing.T) {
	r, err := Filter(baseFilterEnv, FilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matched) != len(baseFilterEnv) {
		t.Errorf("expected all %d keys matched, got %d", len(baseFilterEnv), len(r.Matched))
	}
}

func TestFilter_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	Filter(env, FilterOptions{Prefix: "A"})
	if len(env) != 2 {
		t.Error("input map was mutated")
	}
}

func TestFormatFilterResult_ContainsKeys(t *testing.T) {
	r := FilterResult{
		Matched: map[string]string{"APP_HOST": "localhost"},
		Dropped: []string{"DB_HOST"},
	}
	out := FormatFilterResult(r)
	if !strings.Contains(out, "APP_HOST") {
		t.Error("expected APP_HOST in output")
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Error("expected DB_HOST in dropped output")
	}
	if !strings.Contains(out, "Matched 1") {
		t.Error("expected matched count in output")
	}
}
