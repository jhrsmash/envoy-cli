package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAudit_NoChanges(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	log := Audit(env, env, "test")
	if len(log) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(log))
	}
}

func TestAudit_AddedKey(t *testing.T) {
	before := map[string]string{}
	after := map[string]string{"NEW_KEY": "hello"}
	log := Audit(before, after, "env.test")
	if len(log) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log))
	}
	if log[0].Action != AuditAdded || log[0].Key != "NEW_KEY" {
		t.Errorf("unexpected entry: %+v", log[0])
	}
}

func TestAudit_RemovedKey(t *testing.T) {
	before := map[string]string{"OLD_KEY": "bye"}
	after := map[string]string{}
	log := Audit(before, after, "env.prod")
	if len(log) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log))
	}
	if log[0].Action != AuditRemoved || log[0].OldValue != "bye" {
		t.Errorf("unexpected entry: %+v", log[0])
	}
}

func TestAudit_ChangedKey(t *testing.T) {
	before := map[string]string{"HOST": "localhost"}
	after := map[string]string{"HOST": "prod.example.com"}
	log := Audit(before, after, "env.prod")
	if len(log) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(log))
	}
	e := log[0]
	if e.Action != AuditChanged || e.OldValue != "localhost" || e.NewValue != "prod.example.com" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestFormatAuditLog_Empty(t *testing.T) {
	out := FormatAuditLog(AuditLog{})
	if !strings.Contains(out, "No changes") {
		t.Errorf("expected no-changes message, got: %s", out)
	}
}

func TestFormatAuditLog_ShowsEntries(t *testing.T) {
	before := map[string]string{"A": "1"}
	after := map[string]string{"A": "2", "B": "new"}
	log := Audit(before, after, "myfile")
	out := FormatAuditLog(log)
	if !strings.Contains(out, "myfile") {
		t.Errorf("expected source label in output")
	}
}

func TestSaveAndLoadAuditLog(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	before := map[string]string{"X": "1"}
	after := map[string]string{"X": "2", "Y": "new"}
	log := Audit(before, after, "save-test")

	if err := SaveAuditLog(path, log); err != nil {
		t.Fatalf("SaveAuditLog: %v", err)
	}

	loaded, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(loaded) != len(log) {
		t.Errorf("expected %d entries, got %d", len(log), len(loaded))
	}
}

func TestLoadAuditLog_MissingFile(t *testing.T) {
	_, err := LoadAuditLog("/nonexistent/audit.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected ErrNotExist, got %v", err)
	}
}

func TestSaveAuditLog_AppendsEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")

	log1 := Audit(map[string]string{}, map[string]string{"A": "1"}, "first")
	log2 := Audit(map[string]string{"A": "1"}, map[string]string{"A": "2"}, "second")

	_ = SaveAuditLog(path, log1)
	_ = SaveAuditLog(path, log2)

	loaded, err := LoadAuditLog(path)
	if err != nil {
		t.Fatalf("LoadAuditLog: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("expected 2 entries after append, got %d", len(loaded))
	}
}
