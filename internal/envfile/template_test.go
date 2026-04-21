package envfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderTemplate_NoPlaceholders(t *testing.T) {
	result := RenderTemplate("hello world", map[string]string{"FOO": "bar"})
	if result.Rendered != "hello world" {
		t.Errorf("expected 'hello world', got %q", result.Rendered)
	}
	if len(result.Missing) != 0 {
		t.Errorf("expected no missing keys, got %v", result.Missing)
	}
}

func TestRenderTemplate_ReplacesKnownKeys(t *testing.T) {
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "8080"}
	result := RenderTemplate("http://${APP_HOST}:${APP_PORT}/api", env)
	if result.Rendered != "http://localhost:8080/api" {
		t.Errorf("unexpected rendered output: %q", result.Rendered)
	}
	if len(result.Missing) != 0 {
		t.Errorf("expected no missing keys")
	}
}

func TestRenderTemplate_MissingKeys(t *testing.T) {
	result := RenderTemplate("connect to ${DB_HOST}:${DB_PORT}", map[string]string{})
	if !strings.Contains(result.Rendered, "${DB_HOST}") {
		t.Errorf("expected placeholder to remain, got %q", result.Rendered)
	}
	if len(result.Missing) != 2 {
		t.Errorf("expected 2 missing keys, got %v", result.Missing)
	}
}

func TestRenderTemplate_DeduplicatesMissing(t *testing.T) {
	result := RenderTemplate("${FOO} and ${FOO} again", map[string]string{})
	if len(result.Missing) != 1 {
		t.Errorf("expected 1 unique missing key, got %v", result.Missing)
	}
}

func TestRenderTemplate_PartialEnv(t *testing.T) {
	env := map[string]string{"KNOWN": "value"}
	result := RenderTemplate("${KNOWN} and ${UNKNOWN}", env)
	if !strings.Contains(result.Rendered, "value") {
		t.Errorf("expected KNOWN to be replaced")
	}
	if !strings.Contains(result.Rendered, "${UNKNOWN}") {
		t.Errorf("expected UNKNOWN placeholder to remain")
	}
	if len(result.Missing) != 1 || result.Missing[0] != "UNKNOWN" {
		t.Errorf("expected [UNKNOWN] missing, got %v", result.Missing)
	}
}

func TestRenderTemplateFile_ReadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "tmpl.env")
	_ = os.WriteFile(path, []byte("HOST=${APP_HOST}\nPORT=${APP_PORT}\n"), 0644)

	env := map[string]string{"APP_HOST": "example.com", "APP_PORT": "443"}
	result, err := RenderTemplateFile(path, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Rendered, "example.com") {
		t.Errorf("expected rendered host, got %q", result.Rendered)
	}
}

func TestRenderTemplateFile_MissingFile(t *testing.T) {
	_, err := RenderTemplateFile("/nonexistent/path/tmpl.env", map[string]string{})
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFormatTemplateResult_NoMissing(t *testing.T) {
	r := TemplateResult{Rendered: "output", Missing: nil}
	out := FormatTemplateResult(r)
	if strings.Contains(out, "WARNING") {
		t.Error("expected no warning for empty missing list")
	}
}

func TestFormatTemplateResult_WithMissing(t *testing.T) {
	r := TemplateResult{Rendered: "output", Missing: []string{"DB_PASS"}}
	out := FormatTemplateResult(r)
	if !strings.Contains(out, "WARNING") {
		t.Error("expected WARNING in output")
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Error("expected DB_PASS in warning")
	}
}
