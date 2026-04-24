package envfile

import (
	"strings"
	"testing"
)

var basePinEnv = map[string]string{
	"APP_NAME": "envoy",
	"APP_ENV":  "production",
	"DB_HOST":  "db.example.com",
	"DB_PORT":  "5432",
}

func TestPin_AllKeys(t *testing.T) {
	pinned, result := Pin(basePinEnv, PinOptions{})
	if len(pinned) != len(basePinEnv) {
		t.Fatalf("expected %d pinned keys, got %d", len(basePinEnv), len(pinned))
	}
	if len(result.Skipped) != 0 {
		t.Fatalf("expected no skipped keys, got %v", result.Skipped)
	}
	for k, v := range basePinEnv {
		if pinned[k] != v {
			t.Errorf("key %s: expected %q, got %q", k, v, pinned[k])
		}
	}
}

func TestPin_SelectedKeys(t *testing.T) {
	pinned, result := Pin(basePinEnv, PinOptions{Keys: []string{"APP_NAME", "DB_HOST"}})
	if len(pinned) != 2 {
		t.Fatalf("expected 2 pinned keys, got %d", len(pinned))
	}
	if pinned["APP_NAME"] != "envoy" {
		t.Errorf("unexpected value for APP_NAME: %q", pinned["APP_NAME"])
	}
	if pinned["DB_HOST"] != "db.example.com" {
		t.Errorf("unexpected value for DB_HOST: %q", pinned["DB_HOST"])
	}
	_ = result
}

func TestPin_MissingKeys(t *testing.T) {
	_, result := Pin(basePinEnv, PinOptions{Keys: []string{"APP_NAME", "MISSING_KEY"}})
	if len(result.Skipped) != 1 || result.Skipped[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in skipped, got %v", result.Skipped)
	}
	if len(result.Pinned) != 1 || result.Pinned[0] != "APP_NAME" {
		t.Errorf("expected APP_NAME in pinned, got %v", result.Pinned)
	}
}

func TestPin_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"KEY": "value"}
	pinned, _ := Pin(src, PinOptions{})
	pinned["KEY"] = "mutated"
	if src["KEY"] != "value" {
		t.Error("Pin mutated the source map")
	}
}

func TestPin_EmptyEnv(t *testing.T) {
	pinned, result := Pin(map[string]string{}, PinOptions{})
	if len(pinned) != 0 {
		t.Errorf("expected empty pinned map, got %d keys", len(pinned))
	}
	if len(result.Pinned) != 0 || len(result.Skipped) != 0 {
		t.Errorf("expected empty result, got %+v", result)
	}
}

func TestFormatPinResult_Output(t *testing.T) {
	result := PinResult{
		Pinned:  []string{"APP_NAME", "DB_HOST"},
		Skipped: []string{"GHOST_KEY"},
	}
	out := FormatPinResult(result)
	if !strings.Contains(out, "Pinned: 2") {
		t.Errorf("expected pinned count in output, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Errorf("expected APP_NAME in output, got:\n%s", out)
	}
	if !strings.Contains(out, "GHOST_KEY") {
		t.Errorf("expected GHOST_KEY in skipped output, got:\n%s", out)
	}
}
