package envfile

import (
	"strings"
	"testing"
)

func TestLint_Clean(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"DB_HOST":  "localhost",
	}
	result := Lint(env)
	if !result.OK() {
		t.Errorf("expected no issues, got: %v", result.Issues)
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"app_name": "myapp",
	}
	result := Lint(env)
	if result.OK() {
		t.Fatal("expected lint issue for lowercase key")
	}
	if result.Issues[0].Message != "key should be UPPER_SNAKE_CASE" {
		t.Errorf("unexpected message: %s", result.Issues[0].Message)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
	}
	result := Lint(env)
	if result.OK() {
		t.Fatal("expected lint issue for empty value")
	}
	found := false
	for _, issue := range result.Issues {
		if issue.Key == "API_KEY" && issue.Message == "value is empty" {
			found = true
		}
	}
	if !found {
		t.Error("expected empty value issue for API_KEY")
	}
}

func TestLint_ReservedKey(t *testing.T) {
	env := map[string]string{
		"PATH": "/usr/local/bin",
	}
	result := Lint(env)
	if result.OK() {
		t.Fatal("expected lint issue for reserved key")
	}
	found := false
	for _, issue := range result.Issues {
		if strings.Contains(issue.Message, "reserved") {
			found = true
		}
	}
	if !found {
		t.Error("expected reserved key warning")
	}
}

func TestFormatLintResult_NoIssues(t *testing.T) {
	r := &LintResult{}
	out := FormatLintResult(r)
	if !strings.Contains(out, "No lint issues") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatLintResult_WithIssues(t *testing.T) {
	r := &LintResult{
		Issues: []LintIssue{
			{Key: "bad_key", Message: "key should be UPPER_SNAKE_CASE"},
		},
	}
	out := FormatLintResult(r)
	if !strings.Contains(out, "1 lint issue") {
		t.Errorf("unexpected output: %s", out)
	}
	if !strings.Contains(out, "bad_key") {
		t.Errorf("expected key in output: %s", out)
	}
}
