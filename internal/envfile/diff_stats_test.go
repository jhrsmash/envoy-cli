package envfile

import (
	"strings"
	"testing"
)

func makeDiffEntries2(statuses ...DiffStatus) []DiffEntry {
	entries := make([]DiffEntry, 0, len(statuses))
	for i, st := range statuses {
		entries = append(entries, DiffEntry{
			Key:    fmt.Sprintf("KEY_%d", i),
			Status: st,
		})
	}
	return entries
}

func TestComputeDiffStats_Empty(t *testing.T) {
	s := ComputeDiffStats(nil)
	if s.Total != 0 {
		t.Errorf("expected Total=0, got %d", s.Total)
	}
}

func TestComputeDiffStats_AllAdded(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: StatusAdded},
		{Key: "B", Status: StatusAdded},
	}
	s := ComputeDiffStats(entries)
	if s.Added != 2 || s.Total != 2 {
		t.Errorf("unexpected stats: %+v", s)
	}
}

func TestComputeDiffStats_Mixed(t *testing.T) {
	entries := []DiffEntry{
		{Key: "A", Status: StatusAdded},
		{Key: "B", Status: StatusRemoved},
		{Key: "C", Status: StatusChanged},
		{Key: "D", Status: StatusUnchanged},
		{Key: "E", Status: StatusUnchanged},
	}
	s := ComputeDiffStats(entries)
	if s.Added != 1 {
		t.Errorf("Added: want 1, got %d", s.Added)
	}
	if s.Removed != 1 {
		t.Errorf("Removed: want 1, got %d", s.Removed)
	}
	if s.Changed != 1 {
		t.Errorf("Changed: want 1, got %d", s.Changed)
	}
	if s.Unchanged != 2 {
		t.Errorf("Unchanged: want 2, got %d", s.Unchanged)
	}
	if s.Total != 5 {
		t.Errorf("Total: want 5, got %d", s.Total)
	}
}

func TestFormatDiffStats_NoEntries(t *testing.T) {
	out := FormatDiffStats(DiffStats{})
	if !strings.Contains(out, "no entries") {
		t.Errorf("expected 'no entries', got: %s", out)
	}
}

func TestFormatDiffStats_NoChanges(t *testing.T) {
	s := DiffStats{Unchanged: 3, Total: 3}
	out := FormatDiffStats(s)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes', got: %s", out)
	}
}

func TestFormatDiffStats_WithChanges(t *testing.T) {
	s := DiffStats{Added: 2, Removed: 1, Changed: 1, Unchanged: 4, Total: 8}
	out := FormatDiffStats(s)
	if !strings.Contains(out, "+2 added") {
		t.Errorf("missing added count: %s", out)
	}
	if !strings.Contains(out, "-1 removed") {
		t.Errorf("missing removed count: %s", out)
	}
	if !strings.Contains(out, "~1 changed") {
		t.Errorf("missing changed count: %s", out)
	}
	if !strings.Contains(out, "4 unchanged") {
		t.Errorf("missing unchanged count: %s", out)
	}
}
