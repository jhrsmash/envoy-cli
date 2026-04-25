package envfile

import (
	"strings"
	"testing"
)

func TestResolve_AllPresent(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	opts := ResolveOptions{Required: []string{"HOST", "PORT"}}
	r := Resolve(env, opts)
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", r.Missing)
	}
	if len(r.Defaulted) != 0 {
		t.Errorf("expected no defaulted keys, got %v", r.Defaulted)
	}
}

func TestResolve_MissingKey(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	opts := ResolveOptions{Required: []string{"HOST", "PORT"}}
	r := Resolve(env, opts)
	if len(r.Missing) != 1 || r.Missing[0] != "PORT" {
		t.Errorf("expected [PORT] missing, got %v", r.Missing)
	}
}

func TestResolve_FilledFromDefaults(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	opts := ResolveOptions{
		Required: []string{"HOST", "PORT"},
		Defaults: map[string]string{"PORT": "3000"},
	}
	r := Resolve(env, opts)
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing, got %v", r.Missing)
	}
	if len(r.Defaulted) != 1 || r.Defaulted[0] != "PORT" {
		t.Errorf("expected PORT in defaulted, got %v", r.Defaulted)
	}
	if r.Resolved["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %s", r.Resolved["PORT"])
	}
}

func TestResolve_EmptyValueTreatedAsMissing(t *testing.T) {
	env := map[string]string{"HOST": ""}
	opts := ResolveOptions{
		Required:   []string{"HOST"},
		Defaults:   map[string]string{"HOST": "localhost"},
		AllowEmpty: false,
	}
	r := Resolve(env, opts)
	if len(r.Defaulted) != 1 || r.Defaulted[0] != "HOST" {
		t.Errorf("expected HOST defaulted, got %v", r.Defaulted)
	}
}

func TestResolve_AllowEmpty(t *testing.T) {
	env := map[string]string{"HOST": ""}
	opts := ResolveOptions{
		Required:   []string{"HOST"},
		AllowEmpty: true,
	}
	r := Resolve(env, opts)
	if len(r.Missing) != 0 {
		t.Errorf("expected no missing with AllowEmpty, got %v", r.Missing)
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "1"}
	opts := ResolveOptions{
		Required: []string{"A", "B"},
		Defaults: map[string]string{"B": "2"},
	}
	Resolve(env, opts)
	if _, ok := env["B"]; ok {
		t.Error("Resolve must not mutate the input env map")
	}
}

func TestFormatResolveResult_NoChanges(t *testing.T) {
	r := ResolveResult{Resolved: map[string]string{"A": "1"}}
	out := FormatResolveResult(r)
	if !strings.Contains(out, "all required keys present") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatResolveResult_ShowsMissingAndDefaulted(t *testing.T) {
	r := ResolveResult{
		Resolved:  map[string]string{"PORT": "3000"},
		Missing:   []string{"DB_URL"},
		Defaulted: []string{"PORT"},
	}
	out := FormatResolveResult(r)
	if !strings.Contains(out, "PORT") {
		t.Errorf("expected PORT in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected DB_URL in output, got: %s", out)
	}
	if !strings.Contains(out, "!") {
		t.Errorf("expected missing marker in output, got: %s", out)
	}
}
