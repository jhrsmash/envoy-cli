package envfile

import (
	"strings"
	"testing"
)

func makeSnap(label string, env map[string]string) Snapshot {
	return NewSnapshot(label, env)
}

func TestDiffSnapshots_NoChanges(t *testing.T) {
	a := makeSnap("v1", map[string]string{"A": "1"})
	b := makeSnap("v2", map[string]string{"A": "1"})
	r := DiffSnapshots(a, b)
	if len(r.Changes) != 0 {
		t.Errorf("expected no changes, got %d", len(r.Changes))
	}
}

func TestDiffSnapshots_DetectsChanges(t *testing.T) {
	a := makeSnap("v1", map[string]string{"A": "1", "B": "old"})
	b := makeSnap("v2", map[string]string{"A": "1", "B": "new", "C": "added"})
	r := DiffSnapshots(a, b)
	if len(r.Changes) != 2 {
		t.Errorf("expected 2 changes, got %d", len(r.Changes))
	}
	if r.From != "v1" || r.To != "v2" {
		t.Errorf("unexpected labels: %s %s", r.From, r.To)
	}
}

func TestFormatSnapshotDiff_NoChanges(t *testing.T) {
	a := makeSnap("a", map[string]string{"X": "1"})
	b := makeSnap("b", map[string]string{"X": "1"})
	out := FormatSnapshotDiff(DiffSnapshots(a, b))
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes', got: %s", out)
	}
}

func TestFormatSnapshotDiff_ShowsChanges(t *testing.T) {
	a := makeSnap("staging", map[string]string{"PORT": "8080"})
	b := makeSnap("prod", map[string]string{"PORT": "443", "NEW": "val"})
	out := FormatSnapshotDiff(DiffSnapshots(a, b))
	if !strings.Contains(out, "staging") || !strings.Contains(out, "prod") {
		t.Errorf("expected labels in output")
	}
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output")
	}
	if !strings.Contains(out, "+ NEW") {
		t.Errorf("expected added key in output")
	}
}
