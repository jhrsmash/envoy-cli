package envfile

import (
	"strings"
	"testing"
)

func TestClone_CopiesAllKeys(t *testing.T) {
	src := map[string]string{"APP": "myapp", "PORT": "8080"}
	dst, result := Clone(src, CloneOptions{})

	if len(dst) != len(src) {
		t.Fatalf("expected %d keys, got %d", len(src), len(dst))
	}
	if dst["APP"] != "myapp" || dst["PORT"] != "8080" {
		t.Error("cloned values do not match source")
	}
	if len(result.Keys) != 2 {
		t.Errorf("expected 2 keys in result, got %d", len(result.Keys))
	}
}

func TestClone_DoesNotMutateSource(t *testing.T) {
	src := map[string]string{"SECRET_KEY": "supersecret"}
	_, _ = Clone(src, CloneOptions{RedactSensitive: true})

	if src["SECRET_KEY"] != "supersecret" {
		t.Error("Clone mutated the source map")
	}
}

func TestClone_RedactsSensitiveKeys(t *testing.T) {
	src := map[string]string{
		"APP":        "myapp",
		"SECRET_KEY": "hunter2",
		"DB_PASSWORD": "s3cr3t",
	}
	dst, result := Clone(src, CloneOptions{RedactSensitive: true})

	if dst["APP"] != "myapp" {
		t.Error("non-sensitive key should not be redacted")
	}
	if dst["SECRET_KEY"] == "hunter2" {
		t.Error("SECRET_KEY should have been redacted")
	}
	if len(result.Redacted) == 0 {
		t.Error("expected at least one redacted key in result")
	}
}

func TestClone_ExtraRedactPatterns(t *testing.T) {
	src := map[string]string{"MY_CUSTOM_TOKEN": "abc123", "NORMAL": "value"}
	dst, result := Clone(src, CloneOptions{
		RedactSensitive:     true,
		ExtraRedactPatterns: []string{"TOKEN"},
	})

	if dst["MY_CUSTOM_TOKEN"] == "abc123" {
		t.Error("MY_CUSTOM_TOKEN should have been redacted via extra pattern")
	}
	if dst["NORMAL"] != "value" {
		t.Error("NORMAL should not be redacted")
	}
	_ = result
}

func TestClone_EmptySource(t *testing.T) {
	dst, result := Clone(map[string]string{}, CloneOptions{})
	if len(dst) != 0 {
		t.Error("expected empty destination")
	}
	if len(result.Keys) != 0 {
		t.Error("expected no keys in result")
	}
}

func TestFormatCloneResult_Summary(t *testing.T) {
	r := CloneResult{
		Source:      ".env.staging",
		Destination: ".env.prod",
		Keys:        []string{"APP", "PORT"},
		Redacted:    []string{"SECRET_KEY"},
	}
	out := FormatCloneResult(r)
	if !strings.Contains(out, "2 key(s)") {
		t.Errorf("expected key count in output, got: %s", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected redacted key in output, got: %s", out)
	}
}
