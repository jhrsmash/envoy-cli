package envfile

import (
	"strings"
	"testing"
)

func TestCompare_DisjointMaps(t *testing.T) {
	a := map[string]string{"FOO": "1", "BAR": "2"}
	b := map[string]string{"BAZ": "3", "QUX": "4"}

	r := Compare(a, b)

	if len(r.OnlyInA) != 2 {
		t.Errorf("expected 2 keys only in A, got %d", len(r.OnlyInA))
	}
	if len(r.OnlyInB) != 2 {
		t.Errorf("expected 2 keys only in B, got %d", len(r.OnlyInB))
	}
	if len(r.InBoth) != 0 {
		t.Errorf("expected 0 shared keys, got %d", len(r.InBoth))
	}
}

func TestCompare_OverlappingMaps(t *testing.T) {
	a := map[string]string{"FOO": "1", "SHARED": "x"}
	b := map[string]string{"BAR": "2", "SHARED": "y"}

	r := Compare(a, b)

	if len(r.InBoth) != 1 || r.InBoth[0] != "SHARED" {
		t.Errorf("expected SHARED in both, got %v", r.InBoth)
	}
	if len(r.OnlyInA) != 1 || r.OnlyInA[0] != "FOO" {
		t.Errorf("expected FOO only in A, got %v", r.OnlyInA)
	}
	if len(r.OnlyInB) != 1 || r.OnlyInB[0] != "BAR" {
		t.Errorf("expected BAR only in B, got %v", r.OnlyInB)
	}
}

func TestCompare_IdenticalMaps(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"A": "1", "B": "2"}

	r := Compare(a, b)

	if len(r.OnlyInA) != 0 || len(r.OnlyInB) != 0 {
		t.Errorf("expected no unique keys, got onlyA=%v onlyB=%v", r.OnlyInA, r.OnlyInB)
	}
	if len(r.InBoth) != 2 {
		t.Errorf("expected 2 shared keys, got %d", len(r.InBoth))
	}
}

func TestFormatCompareResult_ContainsLabels(t *testing.T) {
	a := map[string]string{"ONLY_A": "1"}
	b := map[string]string{"ONLY_B": "2"}

	r := Compare(a, b)
	out := FormatCompareResult(r, ".env.staging", ".env.production")

	if !strings.Contains(out, ".env.staging") {
		t.Error("expected label .env.staging in output")
	}
	if !strings.Contains(out, ".env.production") {
		t.Error("expected label .env.production in output")
	}
	if !strings.Contains(out, "ONLY_A") {
		t.Error("expected ONLY_A in output")
	}
}
