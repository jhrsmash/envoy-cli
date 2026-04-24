package envfile

import (
	"strings"
	"testing"
)

func baseTransformEnv() map[string]string {
	return map[string]string{
		"APP_NAME": "MyApp",
		"APP_ENV":  "Production",
		"VERSION":  "v1.2.3",
		"PREFIX":   "prod_database",
	}
}

func TestTransform_Uppercase(t *testing.T) {
	env := baseTransformEnv()
	res, err := Transform(env, TransformOptions{Op: TransformUppercase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["APP_NAME"] != "MYAPP" {
		t.Errorf("expected MYAPP, got %q", res.Output["APP_NAME"])
	}
}

func TestTransform_Lowercase(t *testing.T) {
	env := baseTransformEnv()
	res, err := Transform(env, TransformOptions{Op: TransformLowercase})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["APP_ENV"] != "production" {
		t.Errorf("expected production, got %q", res.Output["APP_ENV"])
	}
}

func TestTransform_SelectedKeys(t *testing.T) {
	env := baseTransformEnv()
	res, err := Transform(env, TransformOptions{
		Op:   TransformLowercase,
		Keys: []string{"APP_NAME"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["APP_NAME"] != "myapp" {
		t.Errorf("expected myapp, got %q", res.Output["APP_NAME"])
	}
	// APP_ENV should be untouched
	if res.Output["APP_ENV"] != "Production" {
		t.Errorf("expected Production untouched, got %q", res.Output["APP_ENV"])
	}
	if len(res.Changed) != 1 || res.Changed[0] != "APP_NAME" {
		t.Errorf("expected only APP_NAME changed, got %v", res.Changed)
	}
}

func TestTransform_TrimPrefix(t *testing.T) {
	env := baseTransformEnv()
	res, err := Transform(env, TransformOptions{
		Op:   TransformTrimPrefix,
		Keys: []string{"PREFIX"},
		Arg1: "prod_",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["PREFIX"] != "database" {
		t.Errorf("expected database, got %q", res.Output["PREFIX"])
	}
}

func TestTransform_Replace(t *testing.T) {
	env := map[string]string{"URL": "http://localhost:8080"}
	res, err := Transform(env, TransformOptions{
		Op:   TransformReplace,
		Arg1: "localhost",
		Arg2: "example.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Output["URL"] != "http://example.com:8080" {
		t.Errorf("unexpected URL: %q", res.Output["URL"])
	}
}

func TestTransform_UnknownOp(t *testing.T) {
	env := baseTransformEnv()
	_, err := Transform(env, TransformOptions{Op: TransformOp("explode")})
	if err == nil {
		t.Fatal("expected error for unknown op")
	}
}

func TestTransform_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	_, _ = Transform(env, TransformOptions{Op: TransformUppercase})
	if env["KEY"] != "value" {
		t.Error("input map was mutated")
	}
}

func TestFormatTransformResult_NoChanges(t *testing.T) {
	r := TransformResult{Output: map[string]string{"A": "a"}, Changed: nil}
	out := FormatTransformResult(r, TransformUppercase)
	if !strings.Contains(out, "no values changed") {
		t.Errorf("expected no-change message, got: %s", out)
	}
}

func TestFormatTransformResult_ShowsChanges(t *testing.T) {
	r := TransformResult{
		Output:  map[string]string{"FOO": "BAR"},
		Changed: []string{"FOO"},
	}
	out := FormatTransformResult(r, TransformUppercase)
	if !strings.Contains(out, "FOO") {
		t.Errorf("expected FOO in output, got: %s", out)
	}
	if !strings.Contains(out, "1 value(s) changed") {
		t.Errorf("expected count in output, got: %s", out)
	}
}
