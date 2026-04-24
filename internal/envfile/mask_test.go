package envfile

import (
	"strings"
	"testing"
)

func TestMask_SensitiveKeysReplaced(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret123",
		"APP_HOST":    "localhost",
	}
	r := Mask(env, MaskOptions{})
	if r.Masked["DB_PASSWORD"] != "***" {
		t.Errorf("expected *** got %q", r.Masked["DB_PASSWORD"])
	}
	if r.Masked["APP_HOST"] != "localhost" {
		t.Errorf("expected localhost got %q", r.Masked["APP_HOST"])
	}
}

func TestMask_ShowLength(t *testing.T) {
	env := map[string]string{
		"API_SECRET": "abcde",
	}
	r := Mask(env, MaskOptions{ShowLength: true})
	if r.Masked["API_SECRET"] != "*****" {
		t.Errorf("expected ***** got %q", r.Masked["API_SECRET"])
	}
}

func TestMask_CustomPlaceholder(t *testing.T) {
	env := map[string]string{
		"DB_TOKEN": "tok_live_xyz",
	}
	r := Mask(env, MaskOptions{Placeholder: "<REDACTED>"})
	if r.Masked["DB_TOKEN"] != "<REDACTED>" {
		t.Errorf("expected <REDACTED> got %q", r.Masked["DB_TOKEN"])
	}
}

func TestMask_ExtraKeys(t *testing.T) {
	env := map[string]string{
		"MY_CUSTOM_KEY": "value",
		"NORMAL_KEY":    "visible",
	}
	r := Mask(env, MaskOptions{ExtraKeys: []string{"MY_CUSTOM_KEY"}})
	if r.Masked["MY_CUSTOM_KEY"] != "***" {
		t.Errorf("expected *** got %q", r.Masked["MY_CUSTOM_KEY"])
	}
	if r.Masked["NORMAL_KEY"] != "visible" {
		t.Errorf("expected visible got %q", r.Masked["NORMAL_KEY"])
	}
}

func TestMask_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "original",
	}
	Mask(env, MaskOptions{})
	if env["DB_PASSWORD"] != "original" {
		t.Errorf("input map was mutated")
	}
}

func TestMask_MaskedKeysSorted(t *testing.T) {
	env := map[string]string{
		"SECRET_Z": "z",
		"SECRET_A": "a",
		"SECRET_M": "m",
	}
	r := Mask(env, MaskOptions{})
	for i := 1; i < len(r.Keys); i++ {
		if r.Keys[i] < r.Keys[i-1] {
			t.Errorf("keys not sorted: %v", r.Keys)
		}
	}
}

func TestFormatMaskResult_NoKeys(t *testing.T) {
	r := MaskResult{Masked: map[string]string{}, Keys: nil}
	out := FormatMaskResult(r)
	if !strings.Contains(out, "No keys masked") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatMaskResult_ShowsCount(t *testing.T) {
	r := MaskResult{
		Masked: map[string]string{"DB_PASSWORD": "***", "API_KEY": "***"},
		Keys:   []string{"API_KEY", "DB_PASSWORD"},
	}
	out := FormatMaskResult(r)
	if !strings.Contains(out, "2 key(s)") {
		t.Errorf("expected count in output, got: %q", out)
	}
	if !strings.Contains(out, "API_KEY") || !strings.Contains(out, "DB_PASSWORD") {
		t.Errorf("expected key names in output, got: %q", out)
	}
}
