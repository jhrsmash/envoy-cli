package envfile

import (
	"strings"
	"testing"
)

func TestSync_AddMissing(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1"}
	res := Sync(src, dst, SyncOptions{AddMissing: true})
	if len(res.Added) != 1 || res.Added[0] != "B" {
		t.Errorf("expected B added, got %v", res.Added)
	}
	if res.Output["B"] != "2" {
		t.Errorf("expected B=2 in output")
	}
}

func TestSync_RemoveExtra(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"A": "1", "Z": "99"}
	res := Sync(src, dst, SyncOptions{RemoveExtra: true})
	if len(res.Removed) != 1 || res.Removed[0] != "Z" {
		t.Errorf("expected Z removed, got %v", res.Removed)
	}
	if _, ok := res.Output["Z"]; ok {
		t.Error("Z should not be in output")
	}
}

func TestSync_Overwrite(t *testing.T) {
	src := map[string]string{"A": "new"}
	dst := map[string]string{"A": "old"}
	res := Sync(src, dst, SyncOptions{Overwrite: true})
	if len(res.Updated) != 1 || res.Updated[0] != "A" {
		t.Errorf("expected A updated, got %v", res.Updated)
	}
	if res.Output["A"] != "new" {
		t.Errorf("expected A=new, got %s", res.Output["A"])
	}
}

func TestSync_NoOp(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"A": "1"}
	res := Sync(src, dst, SyncOptions{AddMissing: true, Overwrite: true, RemoveExtra: true})
	if len(res.Added)+len(res.Removed)+len(res.Updated) != 0 {
		t.Error("expected no changes")
	}
}

func TestSync_DoesNotMutateInput(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"C": "3"}
	Sync(src, dst, SyncOptions{AddMissing: true, RemoveExtra: true})
	if len(dst) != 1 {
		t.Error("dst was mutated")
	}
}

func TestFormatSyncResult_NoChanges(t *testing.T) {
	r := SyncResult{Output: map[string]string{}}
	out := FormatSyncResult(r)
	if out != "sync: no changes" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatSyncResult_ShowsChanges(t *testing.T) {
	r := SyncResult{
		Added:   []string{"NEW"},
		Updated: []string{"MOD"},
		Removed: []string{"OLD"},
		Output:  map[string]string{},
	}
	out := FormatSyncResult(r)
	for _, want := range []string{"+ NEW", "~ MOD", "- OLD"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output:\n%s", want, out)
		}
	}
}
