package envfile

import (
	"strings"
	"testing"
)

func TestDefaults_AppliesMissingKeys(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"A": "99", "B": "2", "C": "3"}

	out, result := Defaults(env, defs, DefaultsOptions{})

	if out["A"] != "1" {
		t.Errorf("expected A=1, got %s", out["A"])
	}
	if out["B"] != "2" {
		t.Errorf("expected B=2, got %s", out["B"])
	}
	if out["C"] != "3" {
		t.Errorf("expected C=3, got %s", out["C"])
	}
	if _, ok := result.Applied["A"]; ok {
		t.Error("A should not be in Applied")
	}
	if result.Applied["B"] != "2" {
		t.Error("B should be in Applied")
	}
	if result.Skipped["A"] != "1" {
		t.Error("A should be in Skipped")
	}
}

func TestDefaults_OverwriteExisting(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"A": "99"}

	out, result := Defaults(env, defs, DefaultsOptions{Overwrite: true})

	if out["A"] != "99" {
		t.Errorf("expected A=99, got %s", out["A"])
	}
	if result.Applied["A"] != "99" {
		t.Error("A should be in Applied")
	}
	if len(result.Skipped) != 0 {
		t.Error("Skipped should be empty")
	}
}

func TestDefaults_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "1"}
	defs := map[string]string{"B": "2"}

	Defaults(env, defs, DefaultsOptions{})

	if _, ok := env["B"]; ok {
		t.Error("env should not be mutated")
	}
}

func TestDefaults_EmptyDefaults(t *testing.T) {
	env := map[string]string{"A": "1"}
	out, result := Defaults(env, map[string]string{}, DefaultsOptions{})

	if len(out) != 1 || out["A"] != "1" {
		t.Error("output should equal input")
	}
	if len(result.Applied) != 0 || len(result.Skipped) != 0 {
		t.Error("result maps should be empty")
	}
}

func TestFormatDefaultsResult_Applied(t *testing.T) {
	_, result := Defaults(
		map[string]string{},
		map[string]string{"PORT": "8080"},
		DefaultsOptions{},
	)
	s := FormatDefaultsResult(result)
	if !strings.Contains(s, "applied") || !strings.Contains(s, "PORT") {
		t.Errorf("unexpected format output: %q", s)
	}
}

func TestFormatDefaultsResult_Nothing(t *testing.T) {
	s := FormatDefaultsResult(DefaultsResult{
		Applied: map[string]string{},
		Skipped: map[string]string{},
	})
	if !strings.Contains(s, "nothing") {
		t.Errorf("expected 'nothing' message, got: %q", s)
	}
}
