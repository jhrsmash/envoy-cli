package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// ResolveResult holds the outcome of resolving an env map against a set of
// required keys, distinguishing resolved, missing, and defaulted entries.
type ResolveResult struct {
	Resolved  map[string]string // keys that were present or filled by defaults
	Missing   []string          // required keys absent from both env and defaults
	Defaulted []string          // keys filled in from defaults (were missing in env)
}

// ResolveOptions controls how Resolve behaves.
type ResolveOptions struct {
	// Required lists keys that must be present after resolution.
	Required []string
	// Defaults provides fallback values for missing keys.
	Defaults map[string]string
	// AllowEmpty, when false, treats keys with empty string values as missing.
	AllowEmpty bool
}

// Resolve checks env against the required keys, filling gaps from defaults.
// It returns a ResolveResult describing what was resolved, defaulted, or missing.
func Resolve(env map[string]string, opts ResolveOptions) ResolveResult {
	resolved := make(map[string]string, len(env))
	for k, v := range env {
		resolved[k] = v
	}

	var missing []string
	var defaulted []string

	for _, key := range opts.Required {
		val, ok := resolved[key]
		if ok && (opts.AllowEmpty || val != "") {
			continue
		}
		// Try defaults
		if def, hasDefault := opts.Defaults[key]; hasDefault && (opts.AllowEmpty || def != "") {
			resolved[key] = def
			defaulted = append(defaulted, key)
		} else {
			missing = append(missing, key)
		}
	}

	sort.Strings(missing)
	sort.Strings(defaulted)

	return ResolveResult{
		Resolved:  resolved,
		Missing:   missing,
		Defaulted: defaulted,
	}
}

// FormatResolveResult returns a human-readable summary of a ResolveResult.
func FormatResolveResult(r ResolveResult) string {
	var sb strings.Builder

	if len(r.Missing) == 0 && len(r.Defaulted) == 0 {
		sb.WriteString("resolve: all required keys present\n")
		return sb.String()
	}

	if len(r.Defaulted) > 0 {
		sb.WriteString(fmt.Sprintf("resolve: %d key(s) filled from defaults:\n", len(r.Defaulted)))
		for _, k := range r.Defaulted {
			sb.WriteString(fmt.Sprintf("  ~ %s = %s\n", k, r.Resolved[k]))
		}
	}

	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("resolve: %d required key(s) missing:\n", len(r.Missing)))
		for _, k := range r.Missing {
			sb.WriteString(fmt.Sprintf("  ! %s\n", k))
		}
	}

	return sb.String()
}
