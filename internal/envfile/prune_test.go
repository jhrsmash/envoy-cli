package envfile

import (
	"strings"
	"testing"
)

func basePruneEnv() map[string]string {
	return map[string]string{
		"APP_NAME":    "myapp",
		"APP_VERSION": "1.0",
		"DB_HOST":     "localhost",
		"DB_PASS":     "",
		"DEBUG":       "true",
		"LEGACY_KEY":  "old",
	}
}

func TestPrune_NoOp(t *testing.T) {
	env := basePruneEnv()
	r := Prune(env, PruneOptions{})
	if len(r.Removed) != 0 {
		t.Errorf("expected no removals, got %v", r.Removed)
	}
	if len(r.Output) != len(env) {
		t.Errorf("expected output length %d, got %d", len(env), len(r.Output))
	}
}

func TestPrune_RemoveEmpty(t *testing.T) {
	r := Prune(basePruneEnv(), PruneOptions{RemoveEmpty: true})
	if len(r.Removed) != 1 || r.Removed[0] != "DB_PASS" {
		t.Errorf("expected [DB_PASS] removed, got %v", r.Removed)
	}
	if _, ok := r.Output["DB_PASS"]; ok {
		t.Error("DB_PASS should not be in output")
	}
}

func TestPrune_ExplicitKeys(t *testing.T) {
	r := Prune(basePruneEnv(), PruneOptions{RemoveKeys: []string{"DEBUG", "LEGACY_KEY"}})
	if len(r.Removed) != 2 {
		t.Errorf("expected 2 removed, got %d", len(r.Removed))
	}
	for _, k := range []string{"DEBUG", "LEGACY_KEY"} {
		if _, ok := r.Output[k]; ok {
			t.Errorf("%s should not be in output", k)
		}
	}
}

func TestPrune_RemovePrefix(t *testing.T) {
	r := Prune(basePruneEnv(), PruneOptions{RemovePrefix: "DB_"})
	for _, k := range r.Removed {
		if !strings.HasPrefix(k, "DB_") {
			t.Errorf("unexpected key removed: %s", k)
		}
	}
	if _, ok := r.Output["DB_HOST"]; ok {
		t.Error("DB_HOST should have been pruned")
	}
	if _, ok := r.Output["DB_PASS"]; ok {
		t.Error("DB_PASS should have been pruned")
	}
}

func TestPrune_DoesNotMutateInput(t *testing.T) {
	env := basePruneEnv()
	orig := make(map[string]string, len(env))
	for k, v := range env {
		orig[k] = v
	}
	Prune(env, PruneOptions{RemoveEmpty: true, RemovePrefix: "APP_"})
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("input mutated at key %s", k)
		}
	}
}

func TestPrune_Mixed(t *testing.T) {
	r := Prune(basePruneEnv(), PruneOptions{
		RemoveEmpty:  true,
		RemoveKeys:   []string{"DEBUG"},
		RemovePrefix: "LEGACY_",
	})
	expected := []string{"DB_PASS", "DEBUG", "LEGACY_KEY"}
	if len(r.Removed) != len(expected) {
		t.Fatalf("expected %d removed, got %d: %v", len(expected), len(r.Removed), r.Removed)
	}
	for i, k := range expected {
		if r.Removed[i] != k {
			t.Errorf("removed[%d]: want %s, got %s", i, k, r.Removed[i])
		}
	}
}

func TestFormatPruneResult_NoRemovals(t *testing.T) {
	r := PruneResult{Output: map[string]string{"A": "1"}, Removed: nil}
	out := FormatPruneResult(r)
	if !strings.Contains(out, "No keys pruned") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatPruneResult_WithRemovals(t *testing.T) {
	r := PruneResult{
		Output:  map[string]string{"A": "1"},
		Removed: []string{"B", "C"},
	}
	out := FormatPruneResult(r)
	if !strings.Contains(out, "Pruned 2 key(s)") {
		t.Errorf("expected count in output, got: %s", out)
	}
	if !strings.Contains(out, "- B") || !strings.Contains(out, "- C") {
		t.Errorf("expected key names in output, got: %s", out)
	}
}
