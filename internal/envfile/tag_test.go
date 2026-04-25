package envfile

import (
	"strings"
	"testing"
)

var baseTagEnv = map[string]string{
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"API_KEY":     "secret",
	"LOG_LEVEL":   "info",
	"FEATURE_FLAG": "true",
}

func TestTag_AssignsTagsToKeys(t *testing.T) {
	opts := TagOptions{
		Tags: map[string][]string{
			"database": {"DB_HOST", "DB_PORT"},
			"security": {"API_KEY"},
		},
	}
	r := Tag(baseTagEnv, nil, opts)
	if r.Tagged["DB_HOST"] != "database" {
		t.Errorf("expected DB_HOST tagged as database, got %q", r.Tagged["DB_HOST"])
	}
	if r.Tagged["DB_PORT"] != "database" {
		t.Errorf("expected DB_PORT tagged as database, got %q", r.Tagged["DB_PORT"])
	}
	if r.Tagged["API_KEY"] != "security" {
		t.Errorf("expected API_KEY tagged as security, got %q", r.Tagged["API_KEY"])
	}
}

func TestTag_SkipsExistingWithoutOverwrite(t *testing.T) {
	existing := map[string]string{"DB_HOST": "infra"}
	opts := TagOptions{
		Tags:      map[string][]string{"database": {"DB_HOST"}},
		Overwrite: false,
	}
	r := Tag(baseTagEnv, existing, opts)
	if r.Tagged["DB_HOST"] != "infra" {
		t.Errorf("expected original tag infra preserved, got %q", r.Tagged["DB_HOST"])
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "DB_HOST" {
		t.Errorf("expected DB_HOST in skipped, got %v", r.Skipped)
	}
}

func TestTag_OverwritesExistingWhenEnabled(t *testing.T) {
	existing := map[string]string{"DB_HOST": "infra"}
	opts := TagOptions{
		Tags:      map[string][]string{"database": {"DB_HOST"}},
		Overwrite: true,
	}
	r := Tag(baseTagEnv, existing, opts)
	if r.Tagged["DB_HOST"] != "database" {
		t.Errorf("expected DB_HOST overwritten to database, got %q", r.Tagged["DB_HOST"])
	}
	if len(r.Skipped) != 0 {
		t.Errorf("expected no skipped keys, got %v", r.Skipped)
	}
}

func TestTag_TracksMissingKeys(t *testing.T) {
	opts := TagOptions{
		Tags: map[string][]string{"ghost": {"NONEXISTENT_KEY"}},
	}
	r := Tag(baseTagEnv, nil, opts)
	if len(r.Missing) != 1 || r.Missing[0] != "NONEXISTENT_KEY" {
		t.Errorf("expected NONEXISTENT_KEY in missing, got %v", r.Missing)
	}
}

func TestTag_DoesNotMutateExisting(t *testing.T) {
	existing := map[string]string{"LOG_LEVEL": "ops"}
	opts := TagOptions{
		Tags:      map[string][]string{"logging": {"LOG_LEVEL"}},
		Overwrite: true,
	}
	Tag(baseTagEnv, existing, opts)
	if existing["LOG_LEVEL"] != "ops" {
		t.Error("Tag mutated the existing tags map")
	}
}

func TestFormatTagResult_NoTags(t *testing.T) {
	r := TagResult{}
	out := FormatTagResult(r)
	if !strings.Contains(out, "No tags applied") {
		t.Errorf("expected 'No tags applied', got %q", out)
	}
}

func TestFormatTagResult_ShowsAllSections(t *testing.T) {
	r := TagResult{
		Tagged:  map[string]string{"DB_HOST": "database"},
		Skipped: []string{"API_KEY"},
		Missing: []string{"GHOST"},
	}
	out := FormatTagResult(r)
	for _, want := range []string{"Tagged", "Skipped", "Missing", "DB_HOST", "API_KEY", "GHOST"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}
