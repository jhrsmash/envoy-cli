package envfile

import (
	"strings"
	"testing"
)

func baseAnnotateEnv() map[string]string {
	return map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"API_KEY":      "secret",
		"PORT":         "8080",
	}
}

func TestAnnotate_AppliesComments(t *testing.T) {
	env := baseAnnotateEnv()
	annotations := map[string]string{
		"DATABASE_URL": "Primary DB",
		"API_KEY":      "Keep secret",
	}
	r := Annotate(env, annotations, AnnotateOptions{})
	if len(r.Annotations) != 2 {
		t.Fatalf("expected 2 annotations, got %d", len(r.Annotations))
	}
	if !strings.HasPrefix(r.Annotations["DATABASE_URL"], "# ") {
		t.Errorf("expected comment prefix '# ', got %q", r.Annotations["DATABASE_URL"])
	}
}

func TestAnnotate_MissingKey(t *testing.T) {
	env := baseAnnotateEnv()
	r := Annotate(env, map[string]string{"MISSING_KEY": "some note"}, AnnotateOptions{})
	if len(r.Missing) != 1 || r.Missing[0] != "MISSING_KEY" {
		t.Errorf("expected MISSING_KEY in Missing, got %v", r.Missing)
	}
	if len(r.Annotations) != 0 {
		t.Errorf("expected no annotations applied")
	}
}

func TestAnnotate_SelectedKeys(t *testing.T) {
	env := baseAnnotateEnv()
	annotations := map[string]string{
		"DATABASE_URL": "DB note",
		"API_KEY":      "API note",
	}
	r := Annotate(env, annotations, AnnotateOptions{Keys: []string{"DATABASE_URL"}})
	if len(r.Annotations) != 1 {
		t.Fatalf("expected 1 annotation, got %d", len(r.Annotations))
	}
	if _, ok := r.Annotations["DATABASE_URL"]; !ok {
		t.Errorf("expected DATABASE_URL to be annotated")
	}
}

func TestAnnotate_DoesNotMutateInput(t *testing.T) {
	env := baseAnnotateEnv()
	orig := make(map[string]string, len(env))
	for k, v := range env {
		orig[k] = v
	}
	Annotate(env, map[string]string{"PORT": "HTTP port"}, AnnotateOptions{})
	for k, v := range orig {
		if env[k] != v {
			t.Errorf("env mutated at key %s", k)
		}
	}
}

func TestAnnotate_EmptyComment_Skipped(t *testing.T) {
	env := baseAnnotateEnv()
	r := Annotate(env, map[string]string{"PORT": ""}, AnnotateOptions{})
	if len(r.Annotations) != 0 {
		t.Errorf("expected empty comment to be skipped, got %v", r.Annotations)
	}
}

func TestFormatAnnotateResult_NoAnnotations(t *testing.T) {
	r := AnnotateResult{}
	out := FormatAnnotateResult(r)
	if !strings.Contains(out, "No annotations") {
		t.Errorf("expected 'No annotations' in output, got %q", out)
	}
}

func TestFormatAnnotateResult_ShowsAnnotated(t *testing.T) {
	r := AnnotateResult{
		Annotations: map[string]string{"API_KEY": "# Keep secret"},
	}
	out := FormatAnnotateResult(r)
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output, got %q", out)
	}
	if !strings.Contains(out, "Annotated") {
		t.Errorf("expected 'Annotated' header in output")
	}
}
