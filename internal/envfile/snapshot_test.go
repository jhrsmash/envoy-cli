package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewSnapshot_CopiesMap(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	s := NewSnapshot("test", env)
	env["A"] = "mutated"
	if s.Env["A"] != "1" {
		t.Errorf("expected snapshot to be independent, got %s", s.Env["A"])
	}
	if s.Label != "test" {
		t.Errorf("expected label 'test', got %s", s.Label)
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s := NewSnapshot("prod", env)

	if err := SaveSnapshot(path, s); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}

	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if loaded.Label != "prod" {
		t.Errorf("expected label 'prod', got %s", loaded.Label)
	}
	if loaded.Env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", loaded.Env["FOO"])
	}
}

func TestLoadSnapshot_InvalidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadSnapshot(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
