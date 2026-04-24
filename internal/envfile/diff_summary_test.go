package envfile

import (
	"strings"
	"testing"
)

func TestSummarizeDiff_Empty(t *testing.T) {
	s := SummarizeDiff(nil)
	if s.Total != 0 || s.Added != 0 || s.Removed != 0 || s.Changed != 0 {
		t.Errorf("expected all zeros, got %+v", s)
	}
}

func TestSummarizeDiff_AllStatuses(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: StatusAdded, NewValue: "1"},
		{Key: "B", Status: StatusRemoved, OldValue: "2"},
		{Key: "C", Status: StatusChanged, OldValue: "3", NewValue: "4"},
		{Key: "D", Status: StatusAdded, NewValue: "5"},
	}
	s := SummarizeDiff(entries)
	if s.Added != 2 {
		t.Errorf("expected 2 added, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", s.Changed)
	}
	if s.Total != 4 {
		t.Errorf("expected total 4, got %d", s.Total)
	}
}

func TestSummarizeDiff_KeysSorted(t *testing.T) {
	entries := []DiffEntry{
		{Key: "Z", Status: StatusAdded},
		{Key: "A", Status: StatusAdded},
		{Key: "M", Status: StatusAdded},
	}
	s := SummarizeDiff(entries)
	if s.AddedKeys[0] != "A" || s.AddedKeys[1] != "M" || s.AddedKeys[2] != "Z" {
		t.Errorf("expected sorted keys, got %v", s.AddedKeys)
	}
}

func TestFormatDiffSummary_NoChanges(t *testing.T) {
	out := FormatDiffSummary(DiffSummary{})
	if out != "No differences found." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatDiffSummary_WithChanges(t *testing.T) {
	s := DiffSummary{
		Added:    1,
		Removed:  1,
		Changed:  1,
		Total:    3,
		AddedKeys:   []string{"NEW_KEY"},
		RemovedKeys: []string{"OLD_KEY"},
		ChangedKeys: []string{"MOD_KEY"},
	}
	out := FormatDiffSummary(s)
	if !strings.Contains(out, "3 change(s)") {
		t.Errorf("expected total count in output: %s", out)
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Errorf("expected added key in output: %s", out)
	}
	if !strings.Contains(out, "OLD_KEY") {
		t.Errorf("expected removed key in output: %s", out)
	}
	if !strings.Contains(out, "MOD_KEY") {
		t.Errorf("expected changed key in output: %s", out)
	}
}

func TestFormatDiffSummary_OnlyAdded(t *testing.T) {
	s := DiffSummary{Added: 2, Total: 2, AddedKeys: []string{"FOO", "BAR"}}
	out := FormatDiffSummary(s)
	if strings.Contains(out, "Removed") || strings.Contains(out, "Changed") {
		t.Errorf("should not show removed/changed sections: %s", out)
	}
	if !strings.Contains(out, "Added") {
		t.Errorf("should show added section: %s", out)
	}
}
