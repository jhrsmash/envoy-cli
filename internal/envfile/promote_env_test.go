package envfile

import (
	"strings"
	"testing"
)

func TestOverrideFromOS_NoMatchingKeys(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	// Use a prefix that won't appear in the test process environment.
	r := OverrideFromOS(env, "ENVOY_TEST_UNIQUE_PREFIX_XYZ_", false)
	if len(r.Applied) != 0 {
		t.Errorf("expected 0 applied, got %d", len(r.Applied))
	}
	if r.Env["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST unchanged")
	}
}

func TestOverrideFromOS_AllowNew(t *testing.T) {
	t.Setenv("ENVOY_OVERRIDE_NEW_KEY", "newval")
	env := map[string]string{"EXISTING": "old"}
	r := OverrideFromOS(env, "ENVOY_OVERRIDE_", true)
	if v, ok := r.Env["NEW_KEY"]; !ok || v != "newval" {
		t.Errorf("expected NEW_KEY=newval, got %q (ok=%v)", v, ok)
	}
	if len(r.Applied) == 0 {
		t.Error("expected at least one applied override")
	}
}

func TestOverrideFromOS_SkipsNewWhenNotAllowed(t *testing.T) {
	t.Setenv("ENVOY_SKIP_BRAND_NEW", "value")
	env := map[string]string{"OTHER": "x"}
	r := OverrideFromOS(env, "ENVOY_SKIP_", false)
	if _, ok := r.Env["BRAND_NEW"]; ok {
		t.Error("should not have inserted BRAND_NEW")
	}
	found := false
	for _, k := range r.Skipped {
		if k == "BRAND_NEW" {
			found = true
		}
	}
	if !found {
		t.Error("expected BRAND_NEW in Skipped")
	}
}

func TestOverrideFromOS_OverwritesExistingKey(t *testing.T) {
	t.Setenv("ENVOY_OVR_PORT", "9090")
	env := map[string]string{"PORT": "8080"}
	r := OverrideFromOS(env, "ENVOY_OVR_", false)
	if r.Env["PORT"] != "9090" {
		t.Errorf("expected PORT=9090, got %q", r.Env["PORT"])
	}
	if len(r.Applied) != 1 {
		t.Errorf("expected 1 applied, got %d", len(r.Applied))
	}
	if r.Applied[0].Previous != "8080" {
		t.Errorf("expected Previous=8080, got %q", r.Applied[0].Previous)
	}
}

func TestOverrideFromOS_DoesNotMutateInput(t *testing.T) {
	t.Setenv("ENVOY_MUT_HOST", "changed")
	env := map[string]string{"HOST": "original"}
	_ = OverrideFromOS(env, "ENVOY_MUT_", false)
	if env["HOST"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestFormatOverrideResult_NoChanges(t *testing.T) {
	r := EnvOverrideResult{Env: map[string]string{}}
	out := FormatOverrideResult(r)
	if !strings.Contains(out, "no OS overrides") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatOverrideResult_ShowsAppliedAndSkipped(t *testing.T) {
	r := EnvOverrideResult{
		Applied: []EnvOverride{
			{Key: "HOST", Value: "new", Previous: "old", WasSet: true},
			{Key: "NEW_KEY", Value: "val", WasSet: false},
		},
		Skipped: []string{"UNKNOWN"},
		Env:     map[string]string{},
	}
	out := FormatOverrideResult(r)
	if !strings.Contains(out, "applied 2") {
		t.Errorf("expected applied count: %q", out)
	}
	if !strings.Contains(out, "skipped 1") {
		t.Errorf("expected skipped count: %q", out)
	}
	if !strings.Contains(out, "HOST") {
		t.Errorf("expected HOST in output")
	}
	if !strings.Contains(out, "UNKNOWN") {
		t.Errorf("expected UNKNOWN in skipped output")
	}
}
