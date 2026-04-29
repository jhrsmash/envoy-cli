package envfile

import (
	"os"
	"strings"
	"testing"
	"time"
)

var (
	old = map[string]string{"APP_ENV": "production", "DB_HOST": "db.prod", "LOG_LEVEL": "warn"}
	curr = map[string]string{"APP_ENV": "staging", "DB_HOST": "db.prod", "NEW_KEY": "hello"}
)

func makeArchive(label string, env map[string]string, t time.Time) ArchiveEntry {
	c := make(map[string]string, len(env))
	for k, v := range env {
		c[k] = v
	}
	return ArchiveEntry{Label: label, Env: c, Timestamp: t}
}

func TestRollback_MostRecent(t *testing.T) {
	now := time.Now()
	archives := []ArchiveEntry{
		makeArchive("v1", old, now.Add(-2*time.Hour)),
		makeArchive("v2", map[string]string{"APP_ENV": "dev"}, now.Add(-1*time.Hour)),
	}
	res, err := Rollback(curr, archives, RollbackOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Label != "v2" {
		t.Errorf("expected label v2, got %q", res.Label)
	}
}

func TestRollback_ByLabel(t *testing.T) {
	now := time.Now()
	archives := []ArchiveEntry{
		makeArchive("v1", old, now.Add(-2*time.Hour)),
		makeArchive("v2", map[string]string{"APP_ENV": "dev"}, now.Add(-1*time.Hour)),
	}
	res, err := Rollback(curr, archives, RollbackOptions{Label: "v1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Label != "v1" {
		t.Errorf("expected label v1, got %q", res.Label)
	}
	if res.Applied["LOG_LEVEL"] != "warn" {
		t.Errorf("expected LOG_LEVEL=warn in applied")
	}
}

func TestRollback_LabelNotFound(t *testing.T) {
	archives := []ArchiveEntry{makeArchive("v1", old, time.Now())}
	_, err := Rollback(curr, archives, RollbackOptions{Label: "missing"})
	if err == nil {
		t.Fatal("expected error for missing label")
	}
}

func TestRollback_NoArchives(t *testing.T) {
	_, err := Rollback(curr, nil, RollbackOptions{})
	if err == nil {
		t.Fatal("expected error when no archives")
	}
}

func TestRollback_DryRun(t *testing.T) {
	archives := []ArchiveEntry{makeArchive("v1", old, time.Now())}
	res, err := Rollback(curr, archives, RollbackOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun=true")
	}
	// Applied should still be current in dry-run.
	if res.Applied["APP_ENV"] != "staging" {
		t.Errorf("dry-run should not change applied env")
	}
}

func TestRollbackStore_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	entries := []ArchiveEntry{
		makeArchive("snap1", map[string]string{"K": "v"}, time.Now()),
	}
	if err := SaveRollbackIndex(dir, entries); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadRollbackIndex(dir)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Label != "snap1" {
		t.Errorf("unexpected loaded entries: %+v", loaded)
	}
}

func TestRollbackStore_MissingDir(t *testing.T) {
	entries, err := LoadRollbackIndex(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty slice, got %v", entries)
	}
}

func TestFormatRollbackResult_DryRun(t *testing.T) {
	res := RollbackResult{
		Label:     "v1",
		Timestamp: time.Now(),
		Applied:   curr,
		Diff:      []DiffEntry{{Key: "APP_ENV", Status: StatusChanged, OldValue: "staging", NewValue: "production"}},
		DryRun:    true,
	}
	out := FormatRollbackResult(res)
	if !strings.Contains(out, "dry-run") {
		t.Errorf("expected dry-run in output, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Errorf("expected APP_ENV in diff output")
	}
}

func TestFormatRollbackResult_Applied(t *testing.T) {
	res := RollbackResult{
		Label:     "v2",
		Timestamp: time.Now(),
		Applied:   old,
		Diff:      nil,
		DryRun:    false,
	}
	out := FormatRollbackResult(res)
	if strings.Contains(out, "dry-run") {
		t.Errorf("should not mention dry-run when not set")
	}
	if !strings.Contains(out, "applied") {
		t.Errorf("expected 'applied' in output")
	}
}

func TestRollbackStore_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/.envoy_rollback_index.json"
	os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := LoadRollbackIndex(dir)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
