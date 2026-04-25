package envfile

import (
	"strings"
	"testing"
)

func TestGroup_ByPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "envoy",
		"APP_ENV":  "production",
		"SECRET":   "abc",
	}
	r := Group(env, GroupOptions{})

	if len(r.Groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(r.Groups["DB"]))
	}
	if len(r.Groups["APP"]) != 2 {
		t.Errorf("expected 2 APP keys, got %d", len(r.Groups["APP"]))
	}
	if _, ok := r.Ungrouped["SECRET"]; !ok {
		t.Error("expected SECRET in ungrouped")
	}
}

func TestGroup_MinGroupSize(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "staging",
	}
	r := Group(env, GroupOptions{MinGroupSize: 2})

	if _, ok := r.Groups["DB"]; !ok {
		t.Error("expected DB group to exist")
	}
	if _, ok := r.Groups["APP"]; ok {
		t.Error("APP group should not exist (only 1 key)")
	}
	if _, ok := r.Ungrouped["APP_ENV"]; !ok {
		t.Error("APP_ENV should be ungrouped")
	}
}

func TestGroup_CustomDelimiter(t *testing.T) {
	env := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
		"plain":   "value",
	}
	r := Group(env, GroupOptions{Delimiter: "."})

	if len(r.Groups["db"]) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(r.Groups["db"]))
	}
	if _, ok := r.Ungrouped["plain"]; !ok {
		t.Error("expected plain in ungrouped")
	}
}

func TestGroup_EmptyEnv(t *testing.T) {
	r := Group(map[string]string{}, GroupOptions{})
	if len(r.Groups) != 0 {
		t.Errorf("expected no groups, got %d", len(r.Groups))
	}
	if len(r.Ungrouped) != 0 {
		t.Errorf("expected no ungrouped, got %d", len(r.Ungrouped))
	}
}

func TestGroup_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	orig := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	Group(env, GroupOptions{})
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("input mutated at key %s", k)
		}
	}
}

func TestFormatGroupResult_ContainsPrefix(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"NOPFX":   "val",
	}
	r := Group(env, GroupOptions{})
	out := FormatGroupResult(r)

	if !strings.Contains(out, "[DB]") {
		t.Error("expected [DB] group header in output")
	}
	if !strings.Contains(out, "[ungrouped]") {
		t.Error("expected [ungrouped] header in output")
	}
	if !strings.Contains(out, "DB_HOST=localhost") {
		t.Error("expected DB_HOST in output")
	}
}

func TestFormatGroupResult_Empty(t *testing.T) {
	r := GroupResult{Groups: map[string]map[string]string{}, Ungrouped: map[string]string{}}
	out := FormatGroupResult(r)
	if !strings.Contains(out, "(no keys)") {
		t.Errorf("expected '(no keys)' for empty result, got: %s", out)
	}
}
