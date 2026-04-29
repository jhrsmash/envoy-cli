package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var baseSearchEnv = map[string]string{
	"APP_NAME":     "myapp",
	"APP_SECRET":   "s3cr3t",
	"DB_HOST":      "localhost",
	"DB_PASSWORD":  "hunter2",
	"FEATURE_FLAG": "true",
}

func TestSearch_MatchesKey(t *testing.T) {
	r, err := Search(baseSearchEnv, SearchOptions{Query: "APP", SearchKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(r.Matches))
	}
}

func TestSearch_MatchesValue(t *testing.T) {
	r, err := Search(baseSearchEnv, SearchOptions{Query: "localhost", SearchValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matches) != 1 || r.Matches[0].Key != "DB_HOST" {
		t.Fatalf("expected DB_HOST match, got %+v", r.Matches)
	}
}

func TestSearch_CaseInsensitiveByDefault(t *testing.T) {
	r, err := Search(baseSearchEnv, SearchOptions{Query: "app", SearchKeys: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(r.Matches))
	}
}

func TestSearch_CaseSensitive_NoMatch(t *testing.T) {
	r, err := Search(baseSearchEnv, SearchOptions{Query: "app", SearchKeys: true, CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matches) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(r.Matches))
	}
}

func TestSearch_RegexMatch(t *testing.T) {
	r, err := Search(baseSearchEnv, SearchOptions{Query: `^DB_`, SearchKeys: true, UseRegex: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(r.Matches))
	}
}

func TestSearch_InvalidRegex(t *testing.T) {
	_, err := Search(baseSearchEnv, SearchOptions{Query: "[invalid", UseRegex: true, SearchKeys: true})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestSearch_DefaultsToKeyAndValue(t *testing.T) {
	r, err := Search(baseSearchEnv, SearchOptions{Query: "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matches) == 0 {
		t.Fatal("expected at least one match")
	}
}

func TestSearch_EmptyEnv(t *testing.T) {
	r, err := Search(map[string]string{}, SearchOptions{Query: "anything"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Matches) != 0 {
		t.Fatalf("expected 0 matches")
	}
}

func TestFormatSearchResult_NoMatches(t *testing.T) {
	r := SearchResult{}
	out := FormatSearchResult(r)
	if !strings.Contains(out, "No matches") {
		t.Errorf("expected no-match message, got: %s", out)
	}
}

func TestFormatSearchResult_WithMatches(t *testing.T) {
	r := SearchResult{
		Matches: []SearchMatch{
			{Key: "APP_NAME", Value: "myapp", MatchedKey: true},
		},
		Options: SearchOptions{Query: "APP"},
	}
	out := FormatSearchResult(r)
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME in output, got: %s", out)
	}
}

func TestExportSearchResult_JSON(t *testing.T) {
	r := SearchResult{
		Matches: []SearchMatch{{Key: "DB_HOST", Value: "localhost", MatchedVal: true}},
		Options: SearchOptions{Query: "localhost"},
	}
	tmp := filepath.Join(t.TempDir(), "result.json")
	if err := ExportSearchResult(r, tmp, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(tmp)
	if !strings.Contains(string(data), "DB_HOST") {
		t.Errorf("expected DB_HOST in JSON output")
	}
}
