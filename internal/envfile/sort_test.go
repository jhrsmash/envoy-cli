package envfile

import (
	"strings"
	"testing"
)

func TestSort_Ascending(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}
	r := Sort(env, SortOptions{})
	if len(r.Keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(r.Keys))
	}
	if r.Keys[0] != "ALPHA" || r.Keys[1] != "MANGO" || r.Keys[2] != "ZEBRA" {
		t.Errorf("unexpected order: %v", r.Keys)
	}
}

func TestSort_Descending(t *testing.T) {
	env := map[string]string{
		"ZEBRA": "z",
		"ALPHA": "a",
		"MANGO": "m",
	}
	r := Sort(env, SortOptions{Reverse: true})
	if r.Keys[0] != "ZEBRA" || r.Keys[1] != "MANGO" || r.Keys[2] != "ALPHA" {
		t.Errorf("unexpected reverse order: %v", r.Keys)
	}
}

func TestSort_GroupByPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "localhost",
		"APP_PORT": "8080",
		"DB_PORT":  "5432",
		"APP_NAME": "envoy",
	}
	r := Sort(env, SortOptions{GroupByPrefix: true})
	// APP_* should come before DB_*
	if r.Keys[0] != "APP_NAME" || r.Keys[1] != "APP_PORT" {
		t.Errorf("expected APP keys first, got: %v", r.Keys)
	}
	if r.Keys[2] != "DB_HOST" || r.Keys[3] != "DB_PORT" {
		t.Errorf("expected DB keys last, got: %v", r.Keys)
	}
}

func TestSort_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"B": "2", "A": "1"}
	original := map[string]string{"B": "2", "A": "1"}
	Sort(env, SortOptions{})
	for k, v := range original {
		if env[k] != v {
			t.Errorf("input mutated at key %s", k)
		}
	}
}

func TestSort_EmptyMap(t *testing.T) {
	r := Sort(map[string]string{}, SortOptions{})
	if len(r.Keys) != 0 {
		t.Errorf("expected empty keys slice")
	}
}

func TestFormatSortResult_ContainsAllKeys(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Sort(env, SortOptions{})
	out := FormatSortResult(r)
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("output missing FOO=bar: %s", out)
	}
	if !strings.Contains(out, "BAZ=qux") {
		t.Errorf("output missing BAZ=qux: %s", out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
