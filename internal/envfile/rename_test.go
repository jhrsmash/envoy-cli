package envfile

import "testing"

func baseEnv() map[string]string {
	return map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"API_KEY": "secret",
	}
}

func TestRename_Success(t *testing.T) {
	out, r := Rename(baseEnv(), "DB_HOST", "DATABASE_HOST", RenameOptions{})
	if !r.Renamed {
		t.Fatalf("expected rename to succeed, got: %s", r.Reason)
	}
	if _, ok := out["DB_HOST"]; ok {
		t.Error("old key should have been removed")
	}
	if out["DATABASE_HOST"] != "localhost" {
		t.Errorf("new key value mismatch: %q", out["DATABASE_HOST"])
	}
}

func TestRename_MissingKey(t *testing.T) {
	_, r := Rename(baseEnv(), "MISSING", "NEW_KEY", RenameOptions{})
	if r.Renamed {
		t.Fatal("expected rename to fail for missing key")
	}
	if r.Reason == "" {
		t.Error("expected a non-empty reason")
	}
}

func TestRename_ConflictNoOverwrite(t *testing.T) {
	_, r := Rename(baseEnv(), "DB_HOST", "DB_PORT", RenameOptions{Overwrite: false})
	if r.Renamed {
		t.Fatal("expected rename to fail due to conflict")
	}
}

func TestRename_ConflictWithOverwrite(t *testing.T) {
	out, r := Rename(baseEnv(), "DB_HOST", "DB_PORT", RenameOptions{Overwrite: true})
	if !r.Renamed {
		t.Fatalf("expected rename to succeed with overwrite, got: %s", r.Reason)
	}
	if out["DB_PORT"] != "localhost" {
		t.Errorf("expected overwritten value 'localhost', got %q", out["DB_PORT"])
	}
}

func TestRename_SameKey(t *testing.T) {
	_, r := Rename(baseEnv(), "DB_HOST", "DB_HOST", RenameOptions{})
	if r.Renamed {
		t.Fatal("expected rename to be skipped for identical keys")
	}
}

func TestRename_EmptyKeys(t *testing.T) {
	_, r := Rename(baseEnv(), "", "NEW", RenameOptions{})
	if r.Renamed {
		t.Fatal("expected rename to fail for empty old key")
	}
	_, r2 := Rename(baseEnv(), "DB_HOST", "", RenameOptions{})
	if r2.Renamed {
		t.Fatal("expected rename to fail for empty new key")
	}
}

func TestRename_DoesNotMutateInput(t *testing.T) {
	env := baseEnv()
	Rename(env, "DB_HOST", "DATABASE_HOST", RenameOptions{})
	if _, ok := env["DB_HOST"]; !ok {
		t.Error("input map should not be mutated")
	}
}

func TestFormatRenameResult_Success(t *testing.T) {
	r := RenameResult{OldKey: "FOO", NewKey: "BAR", Renamed: true}
	msg := FormatRenameResult(r)
	if msg != `renamed "FOO" → "BAR"` {
		t.Errorf("unexpected message: %s", msg)
	}
}

func TestFormatRenameResult_Skipped(t *testing.T) {
	r := RenameResult{OldKey: "FOO", NewKey: "BAR", Renamed: false, Reason: "key \"FOO\" not found"}
	msg := FormatRenameResult(r)
	if msg == "" {
		t.Error("expected non-empty message for skipped rename")
	}
}
