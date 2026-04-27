package envfile

import (
	"strings"
	"testing"
)

func TestTypecast_BoolNormalisation(t *testing.T) {
	env := map[string]string{
		"FEATURE_A": "1",
		"FEATURE_B": "TRUE",
		"FEATURE_C": "False",
		"FEATURE_D": "0",
	}
	res := Typecast(env, TypecastOptions{Target: CastBool})
	expected := map[string]string{
		"FEATURE_A": "true",
		"FEATURE_B": "true",
		"FEATURE_C": "false",
		"FEATURE_D": "false",
	}
	for k, want := range expected {
		if got := res.Output[k]; got != want {
			t.Errorf("key %s: got %q, want %q", k, got, want)
		}
		if _, skipped := res.Skipped[k]; skipped {
			t.Errorf("key %s should not be skipped", k)
		}
	}
}

func TestTypecast_IntNormalisation(t *testing.T) {
	env := map[string]string{
		"PORT":    "  8080  ",
		"TIMEOUT": "30",
		"WORKERS": "not-a-number",
	}
	res := Typecast(env, TypecastOptions{Target: CastInt})

	if got := res.Output["PORT"]; got != "8080" {
		t.Errorf("PORT: got %q, want %q", got, "8080")
	}
	if got := res.Output["TIMEOUT"]; got != "30" {
		t.Errorf("TIMEOUT: got %q, want %q", got, "30")
	}
	if _, ok := res.Skipped["WORKERS"]; !ok {
		t.Error("WORKERS should be in Skipped")
	}
	if res.Output["WORKERS"] != "not-a-number" {
		t.Error("skipped value should be unchanged in Output")
	}
}

func TestTypecast_FloatNormalisation(t *testing.T) {
	env := map[string]string{
		"RATE": "3.14000",
		"BAD":  "abc",
	}
	res := Typecast(env, TypecastOptions{Target: CastFloat})

	if got := res.Output["RATE"]; got != "3.14" {
		t.Errorf("RATE: got %q, want %q", got, "3.14")
	}
	if _, ok := res.Skipped["BAD"]; !ok {
		t.Error("BAD should be skipped")
	}
}

func TestTypecast_SelectedKeys(t *testing.T) {
	env := map[string]string{
		"ENABLED": "TRUE",
		"DEBUG":   "1",
		"NAME":    "envoy",
	}
	res := Typecast(env, TypecastOptions{
		Target: CastBool,
		Keys:   []string{"ENABLED"},
	})

	if res.Output["ENABLED"] != "true" {
		t.Errorf("ENABLED should be cast to 'true'")
	}
	// DEBUG not in Keys, should be unchanged
	if res.Output["DEBUG"] != "1" {
		t.Errorf("DEBUG should be unchanged")
	}
	if _, ok := res.Cast["DEBUG"]; ok {
		t.Error("DEBUG should not appear in Cast")
	}
}

func TestTypecast_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"FLAG": "TRUE"}
	original := env["FLAG"]
	Typecast(env, TypecastOptions{Target: CastBool})
	if env["FLAG"] != original {
		t.Error("Typecast mutated the input map")
	}
}

func TestTypecast_StringTrim(t *testing.T) {
	env := map[string]string{"NAME": "  hello  "}
	res := Typecast(env, TypecastOptions{Target: CastString})
	if res.Output["NAME"] != "hello" {
		t.Errorf("expected trimmed value, got %q", res.Output["NAME"])
	}
}

func TestFormatTypecastResult_NoChanges(t *testing.T) {
	res := TypecastResult{
		Cast:    map[string]string{},
		Skipped: map[string]string{},
		Output:  map[string]string{},
	}
	out := FormatTypecastResult(res)
	if !strings.Contains(out, "no changes") {
		t.Errorf("expected 'no changes', got: %s", out)
	}
}

func TestFormatTypecastResult_ShowsCastAndSkipped(t *testing.T) {
	res := TypecastResult{
		Cast:    map[string]string{"FLAG": "true"},
		Skipped: map[string]string{"BAD": "nope"},
		Output:  map[string]string{"FLAG": "true", "BAD": "nope"},
	}
	out := FormatTypecastResult(res)
	if !strings.Contains(out, "cast") {
		t.Errorf("expected 'cast' section, got: %s", out)
	}
	if !strings.Contains(out, "skipped") {
		t.Errorf("expected 'skipped' section, got: %s", out)
	}
	if !strings.Contains(out, "FLAG") {
		t.Errorf("expected FLAG in output")
	}
	if !strings.Contains(out, "BAD") {
		t.Errorf("expected BAD in output")
	}
}
