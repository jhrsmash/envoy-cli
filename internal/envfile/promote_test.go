package envfile

import (
	"strings"
	"testing"
)

func TestPromote_AllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"C": "3"}

	out, pr := Promote(src, dst, nil, false)

	if len(pr.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(pr.Promoted))
	}
	if out["A"] != "1" || out["B"] != "2" || out["C"] != "3" {
		t.Errorf("unexpected output map: %v", out)
	}
}

func TestPromote_SelectedKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}

	out, pr := Promote(src, dst, []string{"A", "C"}, false)

	if len(pr.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(pr.Promoted))
	}
	if _, ok := out["B"]; ok {
		t.Error("B should not have been promoted")
	}
}

func TestPromote_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	out, pr := Promote(src, dst, nil, false)

	if len(pr.Skipped) != 1 || pr.Skipped[0] != "A" {
		t.Fatalf("expected A to be skipped, got %v", pr.Skipped)
	}
	if out["A"] != "old" {
		t.Errorf("expected old value to be preserved, got %s", out["A"])
	}
}

func TestPromote_OverwritesExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}

	out, pr := Promote(src, dst, nil, true)

	if len(pr.Promoted) != 1 {
		t.Fatalf("expected 1 promoted, got %d", len(pr.Promoted))
	}
	if out["A"] != "new" {
		t.Errorf("expected new value, got %s", out["A"])
	}
}

func TestPromote_MissingKeys(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}

	_, pr := Promote(src, dst, []string{"A", "MISSING"}, false)

	if len(pr.Missing) != 1 || pr.Missing[0] != "MISSING" {
		t.Errorf("expected MISSING in missing list, got %v", pr.Missing)
	}
}

func TestPromote_DoesNotMutateDst(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"B": "2"}
	orig := map[string]string{"B": "2"}

	Promote(src, dst, nil, false)

	for k, v := range orig {
		if dst[k] != v {
			t.Errorf("dst was mutated: key %s changed", k)
		}
	}
	if _, ok := dst["A"]; ok {
		t.Error("dst should not have been mutated with new key A")
	}
}

func TestFormatPromoteResult_Summary(t *testing.T) {
	pr := PromoteResult{
		Promoted: []string{"DB_HOST"},
		Skipped:  []string{"API_KEY"},
		Missing:  []string{"GONE"},
	}
	out := FormatPromoteResult(pr, "staging", "production")

	if !strings.Contains(out, "staging → production") {
		t.Error("expected labels in output")
	}
	if !strings.Contains(out, "promoted: DB_HOST") {
		t.Error("expected promoted key in output")
	}
	if !strings.Contains(out, "skipped") {
		t.Error("expected skipped in output")
	}
	if !strings.Contains(out, "missing") {
		t.Error("expected missing in output")
	}
}

func TestFormatPromoteResult_Empty(t *testing.T) {
	pr := PromoteResult{}
	out := FormatPromoteResult(pr, "dev", "prod")
	if !strings.Contains(out, "nothing to promote") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
