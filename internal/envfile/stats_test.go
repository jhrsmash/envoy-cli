package envfile

import (
	"strings"
	"testing"
)

func TestComputeStats_Empty(t *testing.T) {
	s := ComputeStats(map[string]string{})
	if s.Total != 0 || s.Empty != 0 || s.Sensitive != 0 {
		t.Errorf("expected all-zero stats for empty map, got %+v", s)
	}
}

func TestComputeStats_Total(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2", "C": "3"}
	s := ComputeStats(env)
	if s.Total != 3 {
		t.Errorf("expected Total=3, got %d", s.Total)
	}
}

func TestComputeStats_EmptyValues(t *testing.T) {
	env := map[string]string{"A": "", "B": "hello", "C": ""}
	s := ComputeStats(env)
	if s.Empty != 2 {
		t.Errorf("expected Empty=2, got %d", s.Empty)
	}
}

func TestComputeStats_SensitiveKeys(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "secret",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
	}
	s := ComputeStats(env)
	if s.Sensitive < 2 {
		t.Errorf("expected at least 2 sensitive keys, got %d", s.Sensitive)
	}
}

func TestComputeStats_PrefixGroups(t *testing.T) {
	env := map[string]string{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "envoy",
	}
	s := ComputeStats(env)
	if s.PrefixCounts["DB"] != 2 {
		t.Errorf("expected DB prefix count=2, got %d", s.PrefixCounts["DB"])
	}
	if s.PrefixCounts["APP"] != 1 {
		t.Errorf("expected APP prefix count=1, got %d", s.PrefixCounts["APP"])
	}
}

func TestComputeStats_AvgValueLen(t *testing.T) {
	env := map[string]string{"A": "ab", "B": "abcd"} // lengths 2+4=6, avg=3
	s := ComputeStats(env)
	if s.AvgValueLen != 3.0 {
		t.Errorf("expected AvgValueLen=3.0, got %.1f", s.AvgValueLen)
	}
}

func TestFormatStats_ContainsLabels(t *testing.T) {
	env := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "",
		"APP_NAME":    "envoy",
	}
	s := ComputeStats(env)
	out := FormatStats(s)

	for _, want := range []string{"Total", "Empty", "Sensitive", "Avg val len", "DB_*"} {
		if !strings.Contains(out, want) {
			t.Errorf("FormatStats output missing %q\nGot:\n%s", want, out)
		}
	}
}

func TestFormatStats_NoPrefixGroups(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	s := ComputeStats(env)
	out := FormatStats(s)
	if strings.Contains(out, "Prefix groups") {
		t.Errorf("expected no prefix groups section for key without underscore")
	}
}
