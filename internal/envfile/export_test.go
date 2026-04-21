package envfile

import (
	"encoding/json"
	"strings"
	"testing"
)

var sampleEnv = map[string]string{
	"APP_NAME": "envoy",
	"PORT":     "8080",
	"SECRET":   "s3cr3t",
}

func TestExport_JSON(t *testing.T) {
	out, err := Export(sampleEnv, ExportOptions{Format: FormatJSON, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]string
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if parsed["APP_NAME"] != "envoy" {
		t.Errorf("expected APP_NAME=envoy, got %q", parsed["APP_NAME"])
	}
}

func TestExport_Dotenv(t *testing.T) {
	out, err := Export(map[string]string{"KEY": "value"}, ExportOptions{Format: FormatDotenv})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY=") {
		t.Errorf("expected KEY= in output, got: %s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	out, err := Export(map[string]string{"MY_VAR": "hello"}, ExportOptions{Format: FormatShell})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "export ") {
		t.Errorf("expected shell export prefix, got: %s", out)
	}
}

func TestExport_Redacted(t *testing.T) {
	out, err := Export(sampleEnv, ExportOptions{Format: FormatJSON, Sorted: true, Redacted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "s3cr3t") {
		t.Errorf("expected secret to be redacted, but found plaintext in output")
	}
}

func TestExport_UnsupportedFormat(t *testing.T) {
	_, err := Export(sampleEnv, ExportOptions{Format: ExportFormat("xml")})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestExport_SortedDotenv(t *testing.T) {
	env := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	out, err := Export(env, ExportOptions{Format: FormatDotenv, Sorted: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "A_KEY") {
		t.Errorf("expected first line to start with A_KEY, got: %s", lines[0])
	}
	if !strings.HasPrefix(lines[2], "Z_KEY") {
		t.Errorf("expected last line to start with Z_KEY, got: %s", lines[2])
	}
}
