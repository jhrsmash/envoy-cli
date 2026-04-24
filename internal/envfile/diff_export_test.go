package envfile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var sampleDiffEntries = []DiffEntry{
	{Key: "ADDED_KEY", Status: StatusAdded, NewValue: "hello"},
	{Key: "REMOVED_KEY", Status: StatusRemoved, OldValue: "bye"},
	{Key: "CHANGED_KEY", Status: StatusChanged, OldValue: "old", NewValue: "new"},
}

func TestExportDiff_JSON(t *testing.T) {
	out, err := ExportDiff(sampleDiffEntries, DiffExportJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed []map[string]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(parsed) != 3 {
		t.Errorf("expected 3 entries, got %d", len(parsed))
	}
}

func TestExportDiff_Patch(t *testing.T) {
	out, err := ExportDiff(sampleDiffEntries, DiffExportPatch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+ ADDED_KEY=hello") {
		t.Errorf("missing added line: %s", out)
	}
	if !strings.Contains(out, "- REMOVED_KEY=bye") {
		t.Errorf("missing removed line: %s", out)
	}
	if !strings.Contains(out, "- CHANGED_KEY=old") || !strings.Contains(out, "+ CHANGED_KEY=new") {
		t.Errorf("missing changed lines: %s", out)
	}
}

func TestExportDiff_UnsupportedFormat(t *testing.T) {
	_, err := ExportDiff(sampleDiffEntries, "xml")
	if err == nil {
		t.Error("expected error for unsupported format")
	}
}

func TestExportDiffToFile_JSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "diff.json")
	if err := ExportDiffToFile(sampleDiffEntries, tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(tmp)
	var parsed []map[string]string
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatalf("invalid JSON in file: %v", err)
	}
}

func TestExportDiffToFile_Patch(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "diff.patch")
	if err := ExportDiffToFile(sampleDiffEntries, tmp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := os.ReadFile(tmp)
	if !strings.Contains(string(b), "+ ADDED_KEY") {
		t.Errorf("expected patch content in file: %s", string(b))
	}
}
