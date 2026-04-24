package envfile

import (
	"strings"
	"testing"
)

func TestCopy_AllKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	out, r := Copy(src, dst, CopyOptions{})
	if out["A"] != "1" || out["B"] != "2" {
		t.Fatalf("expected all keys copied, got %v", out)
	}
	if len(r.Copied) != 2 {
		t.Fatalf("expected 2 copied, got %d", len(r.Copied))
	}
}

func TestCopy_SelectedKeys(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2", "C": "3"}
	dst := map[string]string{}
	out, r := Copy(src, dst, CopyOptions{Keys: []string{"A", "C"}})
	if out["A"] != "1" || out["C"] != "3" {
		t.Fatalf("unexpected output: %v", out)
	}
	if _, ok := out["B"]; ok {
		t.Fatal("B should not be copied")
	}
	if len(r.Copied) != 2 {
		t.Fatalf("expected 2 copied, got %d", len(r.Copied))
	}
}

func TestCopy_SkipsExistingWithoutOverwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, r := Copy(src, dst, CopyOptions{})
	if out["A"] != "old" {
		t.Fatalf("expected old value preserved, got %s", out["A"])
	}
	if len(r.Skipped) != 1 || r.Skipped[0] != "A" {
		t.Fatalf("expected A skipped, got %v", r.Skipped)
	}
}

func TestCopy_OverwritesExisting(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	out, r := Copy(src, dst, CopyOptions{Overwrite: true})
	if out["A"] != "new" {
		t.Fatalf("expected new value, got %s", out["A"])
	}
	if len(r.Copied) != 1 {
		t.Fatalf("expected 1 copied, got %d", len(r.Copied))
	}
}

func TestCopy_MissingKeys(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	_, r := Copy(src, dst, CopyOptions{Keys: []string{"A", "MISSING"}})
	if len(r.Missing) != 1 || r.Missing[0] != "MISSING" {
		t.Fatalf("expected MISSING in missing list, got %v", r.Missing)
	}
}

func TestCopy_DoesNotMutateSrc(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{}
	Copy(src, dst, CopyOptions{})
	if len(src) != 1 {
		t.Fatal("src was mutated")
	}
}

func TestFormatCopyResult_Empty(t *testing.T) {
	r := CopyResult{}
	out := FormatCopyResult(r)
	if !strings.Contains(out, "nothing to do") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestFormatCopyResult_ShowsAll(t *testing.T) {
	r := CopyResult{
		Copied:  []string{"A"},
		Skipped: []string{"B"},
		Missing: []string{"C"},
	}
	out := FormatCopyResult(r)
	if !strings.Contains(out, "copied") || !strings.Contains(out, "skipped") || !strings.Contains(out, "missing") {
		t.Fatalf("unexpected output: %s", out)
	}
}
