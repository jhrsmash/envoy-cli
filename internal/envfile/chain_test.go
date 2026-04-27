package envfile

import (
	"errors"
	"strings"
	"testing"
)

func uppercaseValuesStep() ChainStep {
	return ChainStep{
		Name: "uppercase-values",
		Apply: func(env map[string]string) (map[string]string, error) {
			out := make(map[string]string, len(env))
			for k, v := range env {
				out[k] = strings.ToUpper(v)
			}
			return out, nil
		},
	}
}

func addKeyStep(key, val string) ChainStep {
	return ChainStep{
		Name: "add-" + key,
		Apply: func(env map[string]string) (map[string]string, error) {
			out := copyMap(env)
			out[key] = val
			return out, nil
		},
	}
}

func errorStep() ChainStep {
	return ChainStep{
		Name: "failing-step",
		Apply: func(env map[string]string) (map[string]string, error) {
			return nil, errors.New("step failed")
		},
	}
}

func TestChain_NoSteps(t *testing.T) {
	env := map[string]string{"A": "1"}
	r := Chain(env, nil)
	if len(r.Steps) != 0 {
		t.Fatalf("expected 0 steps, got %d", len(r.Steps))
	}
	if r.Final["A"] != "1" {
		t.Errorf("expected Final to equal input")
	}
}

func TestChain_SingleStep(t *testing.T) {
	env := map[string]string{"KEY": "hello"}
	r := Chain(env, []ChainStep{uppercaseValuesStep()})
	if r.Final["KEY"] != "HELLO" {
		t.Errorf("expected HELLO, got %s", r.Final["KEY"])
	}
	if len(r.Steps) != 1 {
		t.Fatalf("expected 1 step result")
	}
	if r.Steps[0].Err != nil {
		t.Errorf("unexpected error: %v", r.Steps[0].Err)
	}
}

func TestChain_MultipleSteps(t *testing.T) {
	env := map[string]string{"KEY": "hello"}
	steps := []ChainStep{
		uppercaseValuesStep(),
		addKeyStep("EXTRA", "added"),
	}
	r := Chain(env, steps)
	if r.Final["KEY"] != "HELLO" {
		t.Errorf("expected HELLO, got %s", r.Final["KEY"])
	}
	if r.Final["EXTRA"] != "added" {
		t.Errorf("expected EXTRA=added")
	}
	if len(r.Steps) != 2 {
		t.Fatalf("expected 2 step results")
	}
}

func TestChain_HaltsOnError(t *testing.T) {
	env := map[string]string{"K": "v"}
	steps := []ChainStep{
		addKeyStep("BEFORE", "yes"),
		errorStep(),
		addKeyStep("AFTER", "never"),
	}
	r := Chain(env, steps)
	if len(r.Steps) != 2 {
		t.Fatalf("expected 2 step results (halted at error), got %d", len(r.Steps))
	}
	if r.Steps[1].Err == nil {
		t.Error("expected error on step 2")
	}
	if _, ok := r.Final["AFTER"]; ok {
		t.Error("step after error should not have run")
	}
}

func TestChain_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"K": "original"}
	Chain(env, []ChainStep{uppercaseValuesStep()})
	if env["K"] != "original" {
		t.Error("input map was mutated")
	}
}

func TestFormatChainResult_NoSteps(t *testing.T) {
	r := ChainResult{Final: map[string]string{}}
	out := FormatChainResult(r)
	if !strings.Contains(out, "no steps") {
		t.Errorf("expected 'no steps' in output, got: %s", out)
	}
}

func TestFormatChainResult_WithSteps(t *testing.T) {
	env := map[string]string{"A": "a", "B": "b"}
	r := Chain(env, []ChainStep{addKeyStep("C", "c"), uppercaseValuesStep()})
	out := FormatChainResult(r)
	if !strings.Contains(out, "add-C") {
		t.Errorf("expected step name in output")
	}
	if !strings.Contains(out, "result:") {
		t.Errorf("expected result summary in output")
	}
}
