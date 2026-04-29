package envfile

import (
	"strings"
	"testing"
)

func TestSummarize_BasicCounts(t *testing.T) {
	env := map[string]string{
		"APP_NAME":    "myapp",
		"APP_VERSION": "1.0",
		"DB_PASSWORD": "secret",
		"EMPTY_KEY":   "",
	}
	r := Summarize(env, SummaryOptions{})
	if r.TotalKeys != 4 {
		t.Errorf("expected 4 keys, got %d", r.TotalKeys)
	}
	if r.EmptyCount != 1 {
		t.Errorf("expected 1 empty, got %d", r.EmptyCount)
	}
	if r.Sensitive != 1 {
		t.Errorf("expected 1 sensitive, got %d", r.Sensitive)
	}
}

func TestSummarize_HidesValuesByDefault(t *testing.T) {
	env := map[string]string{"KEY": "visible"}
	r := Summarize(env, SummaryOptions{})
	if r.Entries[0].Value != "***" {
		t.Errorf("expected value to be hidden, got %q", r.Entries[0].Value)
	}
}

func TestSummarize_ShowValues(t *testing.T) {
	env := map[string]string{"KEY": "visible"}
	r := Summarize(env, SummaryOptions{ShowValues: true})
	if r.Entries[0].Value != "visible" {
		t.Errorf("expected 'visible', got %q", r.Entries[0].Value)
	}
}

func TestSummarize_RedactSensitive(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "topsecret",
		"APP_NAME": "myapp",
	}
	r := Summarize(env, SummaryOptions{ShowValues: true, RedactSensitive: true})
	for _, e := range r.Entries {
		if e.Key == "API_KEY" && e.Value != "[REDACTED]" {
			t.Errorf("expected API_KEY to be redacted, got %q", e.Value)
		}
		if e.Key == "APP_NAME" && e.Value != "myapp" {
			t.Errorf("expected APP_NAME to be visible, got %q", e.Value)
		}
	}
}

func TestSummarize_PrefixGroups(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "prod",
		"NOGROUP": "val",
	}
	r := Summarize(env, SummaryOptions{PrefixGroups: true})
	if len(r.Groups) != 2 {
		t.Errorf("expected 2 groups, got %d: %v", len(r.Groups), r.Groups)
	}
	for _, e := range r.Entries {
		if e.Key == "NOGROUP" && e.Group != "" {
			t.Errorf("expected NOGROUP to have no group")
		}
	}
}

func TestSummarize_EmptyEnv(t *testing.T) {
	r := Summarize(map[string]string{}, SummaryOptions{})
	if r.TotalKeys != 0 {
		t.Errorf("expected 0 keys")
	}
	if len(r.Entries) != 0 {
		t.Errorf("expected no entries")
	}
}

func TestFormatSummaryResult_ContainsStats(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"SECRET":   "",
	}
	r := Summarize(env, SummaryOptions{})
	out := FormatSummaryResult(r)
	if !strings.Contains(out, "Total keys") {
		t.Error("expected 'Total keys' in output")
	}
	if !strings.Contains(out, "Empty") {
		t.Error("expected 'Empty' in output")
	}
	if !strings.Contains(out, "APP_NAME") {
		t.Error("expected key name in output")
	}
}
