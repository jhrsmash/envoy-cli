package envfile

import (
	"testing"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	override := map[string]string{"C": "3"}
	res := Merge(base, override, PreferBase)
	if res.Merged["A"] != "1" || res.Merged["B"] != "2" || res.Merged["C"] != "3" {
		t.Error("expected all keys to be present")
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %v", res.Conflicts)
	}
	if len(res.Added) != 1 || res.Added[0] != "C" {
		t.Errorf("expected Added=[C], got %v", res.Added)
	}
}

func TestMerge_PreferBase(t *testing.T) {
	base := map[string]string{"A": "base"}
	override := map[string]string{"A": "override"}
	res := Merge(base, override, PreferBase)
	if res.Merged["A"] != "base" {
		t.Errorf("expected base value, got %s", res.Merged["A"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_PreferOverride(t *testing.T) {
	base := map[string]string{"A": "base"}
	override := map[string]string{"A": "override"}
	res := Merge(base, override, PreferOverride)
	if res.Merged["A"] != "override" {
		t.Errorf("expected override value, got %s", res.Merged["A"])
	}
	if len(res.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %d", len(res.Conflicts))
	}
}

func TestMerge_Mixed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	override := map[string]string{"B": "changed", "D": "4"}
	res := Merge(base, override, PreferOverride)
	if res.Merged["A"] != "1" {
		t.Error("A should be unchanged")
	}
	if res.Merged["B"] != "changed" {
		t.Error("B should be overridden")
	}
	if res.Merged["D"] != "4" {
		t.Error("D should be added")
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0] != "B" {
		t.Errorf("expected conflict on B, got %v", res.Conflicts)
	}
	if len(res.Added) != 1 || res.Added[0] != "D" {
		t.Errorf("expected added D, got %v", res.Added)
	}
}
