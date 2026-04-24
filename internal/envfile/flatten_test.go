package envfile

import (
	"strings"
	"testing"
)

func TestFlatten_NoOp(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"}
	r := Flatten(env, FlattenOptions{})
	if len(r.Renamed) != 0 {
		t.Errorf("expected 0 renames, got %d", len(r.Renamed))
	}
	if len(r.Unchanged) != 2 {
		t.Errorf("expected 2 unchanged, got %d", len(r.Unchanged))
	}
}

func TestFlatten_ReplacesHyphen(t *testing.T) {
	env := map[string]string{"db-host": "localhost"}
	r := Flatten(env, FlattenOptions{Separator: "_"})
	if v, ok := r.Output["db_host"]; !ok || v != "localhost" {
		t.Errorf("expected db_host=localhost, got output=%v", r.Output)
	}
	if r.Renamed["db-host"] != "db_host" {
		t.Errorf("expected rename entry db-host -> db_host")
	}
}

func TestFlatten_ReplacesDot(t *testing.T) {
	env := map[string]string{"app.port": "9000"}
	r := Flatten(env, FlattenOptions{Separator: "_"})
	if _, ok := r.Output["app_port"]; !ok {
		t.Errorf("expected app_port in output, got %v", r.Output)
	}
}

func TestFlatten_Uppercase(t *testing.T) {
	env := map[string]string{"db-host": "localhost"}
	r := Flatten(env, FlattenOptions{Separator: "_", Uppercase: true})
	if _, ok := r.Output["DB_HOST"]; !ok {
		t.Errorf("expected DB_HOST in output, got %v", r.Output)
	}
}

func TestFlatten_Prefix(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	r := Flatten(env, FlattenOptions{Prefix: "APP_"})
	if v, ok := r.Output["APP_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %v", r.Output)
	}
}

func TestFlatten_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"db-host": "localhost"}
	Flatten(env, FlattenOptions{})
	if _, ok := env["db-host"]; !ok {
		t.Error("input map was mutated")
	}
}

func TestFormatFlattenResult_NoChanges(t *testing.T) {
	r := FlattenResult{
		Output:    map[string]string{"KEY": "val"},
		Renamed:   map[string]string{},
		Unchanged: []string{"KEY"},
	}
	out := FormatFlattenResult(r)
	if !strings.Contains(out, "no keys required renaming") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatFlattenResult_ShowsRenames(t *testing.T) {
	r := FlattenResult{
		Output:    map[string]string{"db_host": "localhost"},
		Renamed:   map[string]string{"db-host": "db_host"},
		Unchanged: []string{},
	}
	out := FormatFlattenResult(r)
	if !strings.Contains(out, "db-host -> db_host") {
		t.Errorf("expected rename line in output, got: %q", out)
	}
	if !strings.Contains(out, "1 key(s) renamed") {
		t.Errorf("expected count summary in output, got: %q", out)
	}
}
