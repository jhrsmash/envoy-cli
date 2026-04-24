package envfile

import (
	"os"
	"strings"
	"testing"
)

func TestInject_SetsNewKeys(t *testing.T) {
	env := map[string]string{"INJECT_FOO": "bar", "INJECT_BAZ": "qux"}
	t.Cleanup(func() { os.Unsetenv("INJECT_FOO"); os.Unsetenv("INJECT_BAZ") })

	result, err := Inject(env, InjectOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(result.Injected))
	}
	if os.Getenv("INJECT_FOO") != "bar" {
		t.Errorf("expected INJECT_FOO=bar")
	}
}

func TestInject_SkipsExistingWithoutOverwrite(t *testing.T) {
	os.Setenv("INJECT_EXISTING", "original")
	t.Cleanup(func() { os.Unsetenv("INJECT_EXISTING") })

	env := map[string]string{"INJECT_EXISTING": "new"}
	result, err := Inject(env, InjectOptions{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped, got %d", len(result.Skipped))
	}
	if os.Getenv("INJECT_EXISTING") != "original" {
		t.Errorf("expected original value to be preserved")
	}
}

func TestInject_OverwritesExistingWhenEnabled(t *testing.T) {
	os.Setenv("INJECT_OW", "old")
	t.Cleanup(func() { os.Unsetenv("INJECT_OW") })

	env := map[string]string{"INJECT_OW": "new"}
	result, err := Inject(env, InjectOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(result.Overwritten))
	}
	if os.Getenv("INJECT_OW") != "new" {
		t.Errorf("expected value to be overwritten")
	}
}

func TestInject_SelectedKeys(t *testing.T) {
	env := map[string]string{"INJECT_A": "1", "INJECT_B": "2", "INJECT_C": "3"}
	t.Cleanup(func() {
		os.Unsetenv("INJECT_A")
		os.Unsetenv("INJECT_B")
		os.Unsetenv("INJECT_C")
	})

	result, err := Inject(env, InjectOptions{Keys: []string{"INJECT_A", "INJECT_C"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 2 {
		t.Errorf("expected 2 injected, got %d", len(result.Injected))
	}
	if os.Getenv("INJECT_B") != "" {
		t.Errorf("INJECT_B should not have been set")
	}
}

func TestInject_EmptyEnv(t *testing.T) {
	result, err := Inject(map[string]string{}, InjectOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Injected) != 0 {
		t.Errorf("expected 0 injected")
	}
}

func TestFormatInjectResult_NoOp(t *testing.T) {
	out := FormatInjectResult(InjectResult{})
	if !strings.Contains(out, "No variables") {
		t.Errorf("expected no-op message, got: %s", out)
	}
}

func TestFormatInjectResult_ShowsAll(t *testing.T) {
	r := InjectResult{
		Injected:    []string{"FOO"},
		Overwritten: []string{"BAR"},
		Skipped:     []string{"BAZ"},
	}
	out := FormatInjectResult(r)
	for _, want := range []string{"FOO", "BAR", "BAZ", "Injected", "Overwritten", "Skipped"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}
