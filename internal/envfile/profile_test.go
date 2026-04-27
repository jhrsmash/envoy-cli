package envfile

import (
	"strings"
	"testing"
)

func TestProfile_BaseOnly(t *testing.T) {
	base := map[string]string{"HOST": "localhost", "PORT": "5432"}
	res, err := Profile(ProfileOptions{Base: base})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", res.Env["HOST"])
	}
	if len(res.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(res.Entries))
	}
	for _, e := range res.Entries {
		if e.Source != "base" {
			t.Errorf("expected source=base for key %q, got %q", e.Key, e.Source)
		}
	}
}

func TestProfile_OverlayAddsKeys(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	overlay := map[string]string{"DEBUG": "true"}
	res, err := Profile(ProfileOptions{
		Base:         base,
		Profiles:     []map[string]string{overlay},
		ProfileNames: []string{"dev"},
		Overwrite:    false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", res.Env["DEBUG"])
	}
	for _, e := range res.Entries {
		if e.Key == "DEBUG" && e.Source != "dev" {
			t.Errorf("expected source=dev for DEBUG, got %q", e.Source)
		}
	}
}

func TestProfile_OverwriteFalse_BaseWins(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	overlay := map[string]string{"HOST": "remotehost"}
	res, err := Profile(ProfileOptions{
		Base:      base,
		Profiles:  []map[string]string{overlay},
		Overwrite: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost (base wins), got %q", res.Env["HOST"])
	}
}

func TestProfile_OverwriteTrue_OverlayWins(t *testing.T) {
	base := map[string]string{"HOST": "localhost"}
	overlay := map[string]string{"HOST": "remotehost"}
	res, err := Profile(ProfileOptions{
		Base:      base,
		Profiles:  []map[string]string{overlay},
		Overwrite: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["HOST"] != "remotehost" {
		t.Errorf("expected HOST=remotehost (overlay wins), got %q", res.Env["HOST"])
	}
}

func TestProfile_MultipleProfiles_OrderMatters(t *testing.T) {
	base := map[string]string{"A": "base"}
	p1 := map[string]string{"A": "p1", "B": "p1"}
	p2 := map[string]string{"A": "p2", "B": "p2", "C": "p2"}
	res, err := Profile(ProfileOptions{
		Base:         base,
		Profiles:     []map[string]string{p1, p2},
		ProfileNames: []string{"staging", "local"},
		Overwrite:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["A"] != "p2" {
		t.Errorf("expected A=p2 (last overwrite wins), got %q", res.Env["A"])
	}
	if res.Env["C"] != "p2" {
		t.Errorf("expected C=p2, got %q", res.Env["C"])
	}
}

func TestProfile_MismatchedNames_ReturnsError(t *testing.T) {
	base := map[string]string{"A": "1"}
	_, err := Profile(ProfileOptions{
		Base:         base,
		Profiles:     []map[string]string{{"B": "2"}},
		ProfileNames: []string{"x", "y"},
	})
	if err == nil {
		t.Fatal("expected error for mismatched ProfileNames length")
	}
}

func TestFormatProfileResult_ContainsSources(t *testing.T) {
	res := ProfileResult{
		Entries: []ProfileEntry{
			{Key: "HOST", Value: "localhost", Source: "base"},
			{Key: "DEBUG", Value: "true", Source: "dev"},
		},
		Env: map[string]string{"HOST": "localhost", "DEBUG": "true"},
	}
	out := FormatProfileResult(res)
	if !strings.Contains(out, "base") {
		t.Errorf("expected output to contain 'base', got: %s", out)
	}
	if !strings.Contains(out, "dev") {
		t.Errorf("expected output to contain 'dev', got: %s", out)
	}
	if !strings.Contains(out, "2 keys resolved") {
		t.Errorf("expected summary line, got: %s", out)
	}
}

func TestFormatProfileResult_Empty(t *testing.T) {
	out := FormatProfileResult(ProfileResult{})
	if !strings.Contains(out, "no keys") {
		t.Errorf("expected 'no keys' message, got: %s", out)
	}
}
