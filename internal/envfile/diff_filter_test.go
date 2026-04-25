package envfile

import (
	"strings"
	"testing"
)

func makeDiffEntries() []DiffEntry {
	return []DiffEntry{
		{Key: "APP_HOST", Status: StatusUnchanged, OldValue: "localhost", NewValue: "localhost"},
		{Key: "APP_PORT", Status: StatusChanged, OldValue: "8080", NewValue: "9090"},
		{Key: "DB_URL", Status: StatusAdded, OldValue: "", NewValue: "postgres://localhost/db"},
		{Key: "LEGACY_KEY", Status: StatusRemoved, OldValue: "old", NewValue: ""},
		{Key: "APP_NAME", Status: StatusUnchanged, OldValue: "myapp", NewValue: "myapp"},
	}
}

func TestFilterDiff_NoOptions(t *testing.T) {
	entries := makeDiffEntries()
	result, err := FilterDiff(entries, DiffFilterOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(result))
	}
}

func TestFilterDiff_ByStatus_Added(t *testing.T) {
	result, err := FilterDiff(makeDiffEntries(), DiffFilterOptions{Statuses: []string{"added"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Key != "DB_URL" {
		t.Errorf("expected only DB_URL (added), got %+v", result)
	}
}

func TestFilterDiff_ExcludeUnchanged(t *testing.T) {
	result, err := FilterDiff(makeDiffEntries(), DiffFilterOptions{ExcludeUnchanged: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range result {
		if e.Status == StatusUnchanged {
			t.Errorf("unexpected unchanged entry: %s", e.Key)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 non-unchanged entries, got %d", len(result))
	}
}

func TestFilterDiff_ByKeyPrefix(t *testing.T) {
	result, err := FilterDiff(makeDiffEntries(), DiffFilterOptions{KeyPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range result {
		if !strings.HasPrefix(e.Key, "APP_") {
			t.Errorf("unexpected key without APP_ prefix: %s", e.Key)
		}
	}
	if len(result) != 3 {
		t.Errorf("expected 3 APP_ entries, got %d", len(result))
	}
}

func TestFilterDiff_InvalidStatus(t *testing.T) {
	_, err := FilterDiff(makeDiffEntries(), DiffFilterOptions{Statuses: []string{"modified"}})
	if err == nil {
		t.Error("expected error for invalid status, got nil")
	}
}

func TestFormatFilteredDiff_Empty(t *testing.T) {
	out := FormatFilteredDiff([]DiffEntry{})
	if !strings.Contains(out, "no matching") {
		t.Errorf("expected 'no matching' message, got: %s", out)
	}
}

func TestFormatFilteredDiff_ShowsSymbols(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: StatusAdded, NewValue: "1"},
		{Key: "B", Status: StatusRemoved, OldValue: "2"},
		{Key: "C", Status: StatusChanged, OldValue: "3", NewValue: "4"},
		{Key: "D", Status: StatusUnchanged, OldValue: "5", NewValue: "5"},
	}
	out := FormatFilteredDiff(entries)
	if !strings.Contains(out, "+ A=") {
		t.Errorf("missing added line: %s", out)
	}
	if !strings.Contains(out, "- B=") {
		t.Errorf("missing removed line: %s", out)
	}
	if !strings.Contains(out, "~ C:") {
		t.Errorf("missing changed line: %s", out)
	}
	if !strings.Contains(out, "  D=") {
		t.Errorf("missing unchanged line: %s", out)
	}
}
