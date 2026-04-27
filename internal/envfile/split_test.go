package envfile

import (
	"strings"
	"testing"
)

var baseSplitEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"APP_NAME":    "envoy",
	"APP_VERSION": "1.0.0",
	"LOG_LEVEL":   "info",
}

func TestSplit_BasicPrefixes(t *testing.T) {
	result, err := Split(baseSplitEnv, SplitOptions{
		Prefixes: map[string]string{"db": "DB_", "app": "APP_"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Buckets["db"]) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(result.Buckets["db"]))
	}
	if len(result.Buckets["app"]) != 2 {
		t.Errorf("expected 2 app keys, got %d", len(result.Buckets["app"]))
	}
	if len(result.Unmatched) != 1 || result.Unmatched[0] != "LOG_LEVEL" {
		t.Errorf("expected LOG_LEVEL as unmatched, got %v", result.Unmatched)
	}
}

func TestSplit_StripPrefix(t *testing.T) {
	result, err := Split(baseSplitEnv, SplitOptions{
		Prefixes:    map[string]string{"db": "DB_"},
		StripPrefix: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result.Buckets["db"]["HOST"]; !ok {
		t.Error("expected stripped key HOST in db bucket")
	}
	if _, ok := result.Buckets["db"]["PORT"]; !ok {
		t.Error("expected stripped key PORT in db bucket")
	}
}

func TestSplit_CatchAll(t *testing.T) {
	result, err := Split(baseSplitEnv, SplitOptions{
		Prefixes: map[string]string{"db": "DB_", "app": "APP_"},
		CatchAll: "other",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Unmatched) != 0 {
		t.Errorf("expected no unmatched keys, got %v", result.Unmatched)
	}
	if v, ok := result.Buckets["other"]["LOG_LEVEL"]; !ok || v != "info" {
		t.Errorf("expected LOG_LEVEL in catch-all bucket")
	}
}

func TestSplit_NoPrefixes_ReturnsError(t *testing.T) {
	_, err := Split(baseSplitEnv, SplitOptions{})
	if err == nil {
		t.Error("expected error for empty prefixes, got nil")
	}
}

func TestSplit_DoesNotMutateInput(t *testing.T) {
	original := map[string]string{"DB_HOST": "localhost", "APP_NAME": "envoy"}
	copy := map[string]string{"DB_HOST": "localhost", "APP_NAME": "envoy"}

	_, _ = Split(original, SplitOptions{
		Prefixes:    map[string]string{"db": "DB_"},
		StripPrefix: true,
	})

	for k, v := range copy {
		if original[k] != v {
			t.Errorf("input mutated at key %s", k)
		}
	}
}

func TestFormatSplitResult_ContainsBucketNames(t *testing.T) {
	result := SplitResult{
		Buckets: map[string]map[string]string{
			"db":  {"HOST": "localhost"},
			"app": {"NAME": "envoy"},
		},
		Unmatched: []string{"LOG_LEVEL"},
	}
	out := FormatSplitResult(result)
	if !strings.Contains(out, "[db]") {
		t.Error("expected [db] in output")
	}
	if !strings.Contains(out, "[app]") {
		t.Error("expected [app] in output")
	}
	if !strings.Contains(out, "unmatched") {
		t.Error("expected unmatched section in output")
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in unmatched output")
	}
}
