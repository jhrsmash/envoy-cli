package envfile

import (
	"strings"
	"testing"
)

func basePivotEnv() map[string]string {
	return map[string]string{
		"DEV_DB_HOST":  "localhost",
		"DEV_DB_PORT":  "5432",
		"PROD_DB_HOST": "db.example.com",
		"PROD_DB_PORT": "5432",
		"PROD_API_KEY": "secret",
	}
}

func TestPivot_BasicPrefixes(t *testing.T) {
	env := basePivotEnv()
	res, err := Pivot(env, PivotOptions{Prefixes: []string{"DEV", "PROD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.RowKeys) != 3 {
		t.Fatalf("expected 3 row keys, got %d: %v", len(res.RowKeys), res.RowKeys)
	}
	if res.Rows["DB_HOST"]["DEV"] != "localhost" {
		t.Errorf("expected DEV DB_HOST=localhost")
	}
	if res.Rows["DB_HOST"]["PROD"] != "db.example.com" {
		t.Errorf("expected PROD DB_HOST=db.example.com")
	}
}

func TestPivot_MissingColumnValue(t *testing.T) {
	env := basePivotEnv()
	res, err := Pivot(env, PivotOptions{Prefixes: []string{"DEV", "PROD"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// API_KEY only exists under PROD
	if _, ok := res.Rows["API_KEY"]["DEV"]; ok {
		t.Errorf("expected DEV API_KEY to be absent")
	}
	if res.Rows["API_KEY"]["PROD"] != "secret" {
		t.Errorf("expected PROD API_KEY=secret")
	}
}

func TestPivot_NoPrefixes_ReturnsError(t *testing.T) {
	_, err := Pivot(map[string]string{"FOO": "bar"}, PivotOptions{})
	if err == nil {
		t.Fatal("expected error for empty prefixes")
	}
}

func TestPivot_CustomDelimiter(t *testing.T) {
	env := map[string]string{
		"DEV.HOST": "localhost",
		"PRD.HOST": "prod.host",
	}
	res, err := Pivot(env, PivotOptions{Prefixes: []string{"DEV", "PRD"}, Delimiter: "."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Rows["HOST"]["DEV"] != "localhost" {
		t.Errorf("expected DEV HOST=localhost")
	}
}

func TestPivot_DoesNotMutateInput(t *testing.T) {
	env := basePivotEnv()
	orig := make(map[string]string, len(env))
	for k, v := range env {
		orig[k] = v
	}
	_, _ = Pivot(env, PivotOptions{Prefixes: []string{"DEV", "PROD"}})
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("input mutated at key %s", k)
		}
	}
}

func TestFormatPivotResult_ContainsHeaders(t *testing.T) {
	env := basePivotEnv()
	res, _ := Pivot(env, PivotOptions{Prefixes: []string{"DEV", "PROD"}})
	out := FormatPivotResult(res)
	if !strings.Contains(out, "DEV") || !strings.Contains(out, "PROD") {
		t.Errorf("expected column headers in output: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected row key DB_HOST in output: %s", out)
	}
}

func TestFormatPivotResult_EmptyResult(t *testing.T) {
	res := PivotResult{}
	out := FormatPivotResult(res)
	if !strings.Contains(out, "no matching") {
		t.Errorf("expected empty message, got: %s", out)
	}
}
