package envfile

import (
	"strings"
	"testing"
)

var baseExtractEnv = map[string]string{
	"APP_HOST":     "localhost",
	"APP_PORT":     "8080",
	"APP_SECRET":   "topsecret",
	"DB_HOST":      "db.local",
	"DB_PORT":      "5432",
	"UNRELATED_KEY": "value",
}

func TestExtract_AllKeys(t *testing.T) {
	res, err := Extract(baseExtractEnv, ExtractOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Extracted) != len(baseExtractEnv) {
		t.Errorf("expected %d keys, got %d", len(baseExtractEnv), len(res.Extracted))
	}
}

func TestExtract_ByPrefix(t *testing.T) {
	res, err := Extract(baseExtractEnv, ExtractOptions{Prefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Extracted) != 3 {
		t.Errorf("expected 3 APP_ keys, got %d", len(res.Extracted))
	}
	if _, ok := res.Extracted["APP_HOST"]; !ok {
		t.Error("expected APP_HOST in extracted")
	}
}

func TestExtract_StripPrefix(t *testing.T) {
	res, err := Extract(baseExtractEnv, ExtractOptions{Prefix: "APP_", StripPrefix: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Extracted["HOST"]; !ok {
		t.Error("expected HOST after stripping APP_ prefix")
	}
	if _, ok := res.Extracted["APP_HOST"]; ok {
		t.Error("expected APP_HOST to be stripped")
	}
}

func TestExtract_ExplicitKeys(t *testing.T) {
	res, err := Extract(baseExtractEnv, ExtractOptions{Keys: []string{"APP_HOST", "DB_PORT"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Extracted) != 2 {
		t.Errorf("expected 2 keys, got %d", len(res.Extracted))
	}
}

func TestExtract_MissingKeyNoFail(t *testing.T) {
	res, err := Extract(baseExtractEnv, ExtractOptions{
		Keys:          []string{"APP_HOST", "MISSING_KEY"},
		FailOnMissing: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in missing list, got %v", res.Missing)
	}
}

func TestExtract_FailOnMissing(t *testing.T) {
	_, err := Extract(baseExtractEnv, ExtractOptions{
		Keys:          []string{"APP_HOST", "GHOST"},
		FailOnMissing: true,
	})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "GHOST") {
		t.Errorf("expected error to mention GHOST, got: %v", err)
	}
}

func TestExtract_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	_, _ = Extract(env, ExtractOptions{Prefix: "A", StripPrefix: true})
	if _, ok := env["A"]; !ok {
		t.Error("Extract mutated the input map")
	}
}

func TestFormatExtractResult_ContainsKeys(t *testing.T) {
	res := ExtractResult{
		Extracted: map[string]string{"HOST": "localhost", "PORT": "8080"},
		Missing:   []string{"SECRET"},
	}
	out := FormatExtractResult(res)
	if !strings.Contains(out, "HOST") {
		t.Error("expected HOST in format output")
	}
	if !strings.Contains(out, "missing") {
		t.Error("expected 'missing' label in format output")
	}
	if !strings.Contains(out, "SECRET") {
		t.Error("expected SECRET in missing section")
	}
}
