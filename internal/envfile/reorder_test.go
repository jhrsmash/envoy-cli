package envfile

import (
	"strings"
	"testing"
)

var baseReorderEnv = map[string]string{
	"APP_NAME":    "myapp",
	"APP_VERSION": "1.0",
	"DB_HOST":     "localhost",
	"DB_PORT":     "5432",
	"SECRET_KEY":  "s3cr3t",
}

func TestReorder_ExplicitOrder(t *testing.T) {
	opts := ReorderOptions{
		Keys: []string{"SECRET_KEY", "APP_NAME", "APP_VERSION"},
	}
	_, result := Reorder(baseReorderEnv, opts)

	if len(result.Ordered) != len(baseReorderEnv) {
		t.Fatalf("expected %d ordered keys, got %d", len(baseReorderEnv), len(result.Ordered))
	}
	if result.Ordered[0] != "SECRET_KEY" {
		t.Errorf("expected first key SECRET_KEY, got %s", result.Ordered[0])
	}
	if result.Ordered[1] != "APP_NAME" {
		t.Errorf("expected second key APP_NAME, got %s", result.Ordered[1])
	}
}

func TestReorder_MissingKeys(t *testing.T) {
	opts := ReorderOptions{
		Keys: []string{"APP_NAME", "DOES_NOT_EXIST"},
	}
	_, result := Reorder(baseReorderEnv, opts)

	if len(result.Missing) != 1 || result.Missing[0] != "DOES_NOT_EXIST" {
		t.Errorf("expected missing [DOES_NOT_EXIST], got %v", result.Missing)
	}
}

func TestReorder_UnlistedKeys(t *testing.T) {
	opts := ReorderOptions{
		Keys: []string{"APP_NAME"},
	}
	_, result := Reorder(baseReorderEnv, opts)

	if len(result.Unlisted) != 4 {
		t.Errorf("expected 4 unlisted keys, got %d", len(result.Unlisted))
	}
}

func TestReorder_AlphaTail(t *testing.T) {
	opts := ReorderOptions{
		Keys:      []string{"SECRET_KEY"},
		AlphaTail: true,
	}
	_, result := Reorder(baseReorderEnv, opts)

	// Tail should be sorted alphabetically.
	for i := 2; i < len(result.Unlisted); i++ {
		if result.Unlisted[i] < result.Unlisted[i-1] {
			t.Errorf("unlisted keys not sorted: %v", result.Unlisted)
			break
		}
	}
}

func TestReorder_DoesNotMutateInput(t *testing.T) {
	orig := map[string]string{"A": "1", "B": "2"}
	opts := ReorderOptions{Keys: []string{"B", "A"}}
	out, _ := Reorder(orig, opts)
	out["C"] = "3"

	if _, ok := orig["C"]; ok {
		t.Error("Reorder mutated the input map")
	}
}

func TestReorder_EmptyEnv(t *testing.T) {
	_, result := Reorder(map[string]string{}, ReorderOptions{Keys: []string{"X"}})
	if len(result.Ordered) != 0 {
		t.Errorf("expected 0 ordered keys for empty env, got %d", len(result.Ordered))
	}
	if len(result.Missing) != 1 {
		t.Errorf("expected 1 missing key, got %d", len(result.Missing))
	}
}

func TestFormatReorderResult_ContainsSummary(t *testing.T) {
	r := ReorderResult{
		Ordered:  []string{"A", "B", "C"},
		Missing:  []string{"GHOST"},
		Unlisted: []string{"B", "C"},
	}
	out := FormatReorderResult(r)
	if !strings.Contains(out, "Reordered: 3 keys") {
		t.Errorf("expected reorder count in output, got: %s", out)
	}
	if !strings.Contains(out, "GHOST") {
		t.Errorf("expected missing key GHOST in output, got: %s", out)
	}
}
