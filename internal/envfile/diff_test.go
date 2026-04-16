package envfile

import "testing"

func TestDiff_NoChanges(t *testing.T) {
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(base, target)
	if !result.IsEmpty() {
		t.Errorf("expected no diff, got added=%v removed=%v changed=%v",
			result.Added, result.Removed, result.Changed)
	}
}

func TestDiff_AddedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "NEW_KEY": "newval"}

	result := Diff(base, target)
	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added key, got %d", len(result.Added))
	}
	if result.Added["NEW_KEY"] != "newval" {
		t.Errorf("unexpected added value: %s", result.Added["NEW_KEY"])
	}
}

func TestDiff_RemovedKeys(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD_KEY": "oldval"}
	target := map[string]string{"FOO": "bar"}

	result := Diff(base, target)
	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed key, got %d", len(result.Removed))
	}
	if result.Removed["OLD_KEY"] != "oldval" {
		t.Errorf("unexpected removed value: %s", result.Removed["OLD_KEY"])
	}
}

func TestDiff_ChangedKeys(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	target := map[string]string{"FOO": "new"}

	result := Diff(base, target)
	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(result.Changed))
	}
	pair, ok := result.Changed["FOO"]
	if !ok {
		t.Fatal("expected FOO in changed keys")
	}
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("unexpected changed values: %v", pair)
	}
}

func TestDiff_Mixed(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2", "C": "3"}
	target := map[string]string{"A": "1", "B": "changed", "D": "4"}

	result := Diff(base, target)
	if len(result.Added) != 1 || result.Added["D"] != "4" {
		t.Errorf("unexpected added: %v", result.Added)
	}
	if len(result.Removed) != 1 || result.Removed["C"] != "3" {
		t.Errorf("unexpected removed: %v", result.Removed)
	}
	if len(result.Changed) != 1 || result.Changed["B"] != ([2]string{"2", "changed"}) {
		t.Errorf("unexpected changed: %v", result.Changed)
	}
}
