package envfile

import (
	"strings"
	"testing"
)

func TestClassify_AuthKeys(t *testing.T) {
	env := map[string]string{
		"API_TOKEN":   "abc",
		"DB_PASSWORD": "secret",
		"APP_NAME":    "myapp",
	}

	c := Classify(env)

	if _, ok := c.Groups["auth"]; !ok {
		t.Fatal("expected 'auth' group")
	}
	found := false
	for _, k := range c.Groups["auth"] {
		if k == "API_TOKEN" {
			found = true
		}
	}
	if !found {
		t.Error("expected API_TOKEN in auth group")
	}
}

func TestClassify_DatabaseKeys(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}

	c := Classify(env)

	if _, ok := c.Groups["database"]; !ok {
		t.Fatal("expected 'database' group")
	}
	if len(c.Groups["database"]) != 2 {
		t.Errorf("expected 2 database keys, got %d", len(c.Groups["database"]))
	}
}

func TestClassify_Unclassified(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"TIMEZONE": "UTC",
	}

	c := Classify(env)

	if len(c.Unclassified) != 2 {
		t.Errorf("expected 2 unclassified keys, got %d", len(c.Unclassified))
	}
}

func TestClassify_EmptyEnv(t *testing.T) {
	c := Classify(map[string]string{})

	if len(c.Groups) != 0 {
		t.Errorf("expected no groups, got %d", len(c.Groups))
	}
	if len(c.Unclassified) != 0 {
		t.Errorf("expected no unclassified, got %d", len(c.Unclassified))
	}
}

func TestClassify_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "val"}
	orig := map[string]string{"SECRET_KEY": "val"}

	Classify(env)

	if env["SECRET_KEY"] != orig["SECRET_KEY"] {
		t.Error("input map was mutated")
	}
}

func TestFormatClassifyResult_NoKeys(t *testing.T) {
	c := Classification{}
	out := FormatClassifyResult(c)
	if !strings.Contains(out, "no keys") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatClassifyResult_ShowsCategories(t *testing.T) {
	env := map[string]string{
		"DB_HOST":   "localhost",
		"API_TOKEN": "xyz",
		"APP_NAME":  "myapp",
	}

	c := Classify(env)
	out := FormatClassifyResult(c)

	if !strings.Contains(out, "[auth]") {
		t.Error("expected [auth] section in output")
	}
	if !strings.Contains(out, "[database]") {
		t.Error("expected [database] section in output")
	}
	if !strings.Contains(out, "[unclassified]") {
		t.Error("expected [unclassified] section in output")
	}
}
