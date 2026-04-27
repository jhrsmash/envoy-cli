package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArchive_SavesAndLoads(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}

	res, err := Archive(env, "v1", dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Label != "v1" {
		t.Errorf("expected label v1, got %s", res.Label)
	}
	if res.Keys != 2 {
		t.Errorf("expected 2 keys, got %d", res.Keys)
	}

	entry, err := LoadArchive("v1", dir)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if entry.Env["APP_ENV"] != "production" {
		t.Errorf("expected production, got %s", entry.Env["APP_ENV"])
	}
	if entry.Env["PORT"] != "8080" {
		t.Errorf("expected 8080, got %s", entry.Env["PORT"])
	}
}

func TestArchive_DoesNotMutateInput(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"KEY": "val"}
	_, _ = Archive(env, "snap", dir)
	env["KEY"] = "mutated"

	entry, _ := LoadArchive("snap", dir)
	if entry.Env["KEY"] != "val" {
		t.Errorf("archive mutated: got %s", entry.Env["KEY"])
	}
}

func TestArchive_EmptyLabel(t *testing.T) {
	dir := t.TempDir()
	_, err := Archive(map[string]string{"A": "1"}, "", dir)
	if err == nil {
		t.Error("expected error for empty label")
	}
}

func TestLoadArchive_MissingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadArchive("nonexistent", dir)
	if err == nil {
		t.Error("expected error for missing archive")
	}
}

func TestListArchives_Empty(t *testing.T) {
	dir := t.TempDir()
	labels, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 0 {
		t.Errorf("expected 0 labels, got %d", len(labels))
	}
}

func TestListArchives_ReturnsLabels(t *testing.T) {
	dir := t.TempDir()
	env := map[string]string{"X": "1"}
	_, _ = Archive(env, "alpha", dir)
	_, _ = Archive(env, "beta", dir)

	labels, err := ListArchives(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestListArchives_MissingDir(t *testing.T) {
	labels, err := ListArchives(filepath.Join(os.TempDir(), "envoy-no-such-dir-xyz"))
	if err != nil {
		t.Fatalf("expected nil error for missing dir, got %v", err)
	}
	if labels != nil {
		t.Errorf("expected nil labels, got %v", labels)
	}
}

func TestFormatArchiveList_Empty(t *testing.T) {
	out := FormatArchiveList(nil)
	if out != "No archives found.\n" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatArchiveList_WithLabels(t *testing.T) {
	out := FormatArchiveList([]string{"v1", "v2"})
	if out == "" {
		t.Error("expected non-empty output")
	}
}
