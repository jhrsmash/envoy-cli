package envfile

import (
	"strings"
	"testing"
)

func basePatchEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "myapp",
		"APP_ENV":  "development",
		"DB_HOST":  "localhost",
	}
}

func TestPatch_SetNewKey(t *testing.T) {
	out, result := Patch(basePatchEnv(), []PatchOp{
		{Action: "set", Key: "NEW_KEY", Value: "hello"},
	})
	if out["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", out["NEW_KEY"])
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(result.Applied))
	}
}

func TestPatch_SetOverwritesExisting(t *testing.T) {
	out, _ := Patch(basePatchEnv(), []PatchOp{
		{Action: "set", Key: "APP_ENV", Value: "production"},
	})
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", out["APP_ENV"])
	}
}

func TestPatch_DeleteExistingKey(t *testing.T) {
	out, result := Patch(basePatchEnv(), []PatchOp{
		{Action: "delete", Key: "DB_HOST"},
	})
	if _, ok := out["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be deleted")
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(result.Applied))
	}
}

func TestPatch_DeleteMissingKey(t *testing.T) {
	_, result := Patch(basePatchEnv(), []PatchOp{
		{Action: "delete", Key: "NONEXISTENT"},
	})
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped op, got %d", len(result.Skipped))
	}
	if len(result.Errors) == 0 {
		t.Error("expected an error message for missing key delete")
	}
}

func TestPatch_RenameKey(t *testing.T) {
	out, result := Patch(basePatchEnv(), []PatchOp{
		{Action: "rename", Key: "APP_NAME", NewKey: "SERVICE_NAME"},
	})
	if _, ok := out["APP_NAME"]; ok {
		t.Error("expected APP_NAME to be removed after rename")
	}
	if out["SERVICE_NAME"] != "myapp" {
		t.Errorf("expected SERVICE_NAME=myapp, got %q", out["SERVICE_NAME"])
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(result.Applied))
	}
}

func TestPatch_UnknownAction(t *testing.T) {
	_, result := Patch(basePatchEnv(), []PatchOp{
		{Action: "upsert", Key: "FOO"},
	})
	if len(result.Skipped) != 1 {
		t.Errorf("expected 1 skipped op, got %d", len(result.Skipped))
	}
}

func TestPatch_DoesNotMutateInput(t *testing.T) {
	env := basePatchEnv()
	Patch(env, []PatchOp{
		{Action: "set", Key: "APP_NAME", Value: "changed"},
		{Action: "delete", Key: "DB_HOST"},
	})
	if env["APP_NAME"] != "myapp" {
		t.Error("input map was mutated by Patch")
	}
	if _, ok := env["DB_HOST"]; !ok {
		t.Error("input map was mutated: DB_HOST was deleted")
	}
}

func TestFormatPatchResult_Applied(t *testing.T) {
	_, result := Patch(basePatchEnv(), []PatchOp{
		{Action: "set", Key: "FOO", Value: "bar"},
		{Action: "delete", Key: "APP_ENV"},
		{Action: "rename", Key: "APP_NAME", NewKey: "SERVICE_NAME"},
	})
	out := FormatPatchResult(result)
	if !strings.Contains(out, "SET") {
		t.Error("expected SET in output")
	}
	if !strings.Contains(out, "DELETE") {
		t.Error("expected DELETE in output")
	}
	if !strings.Contains(out, "RENAME") {
		t.Error("expected RENAME in output")
	}
}

func TestFormatPatchResult_Empty(t *testing.T) {
	out := FormatPatchResult(PatchResult{})
	if !strings.Contains(out, "No patch operations") {
		t.Error("expected empty message")
	}
}
