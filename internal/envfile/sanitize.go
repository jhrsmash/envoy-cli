package envfile

import (
	"fmt"
	"strings"
)

// SanitizeOptions controls how sanitization is applied.
type SanitizeOptions struct {
	// StripControlChars removes non-printable ASCII control characters from values.
	StripControlChars bool
	// NormalizeLineEndings converts \r\n and bare \r to \n in values.
	NormalizeLineEndings bool
	// TrimQuotes removes surrounding single or double quotes from values.
	TrimQuotes bool
	// MaxValueLength truncates values longer than this limit (0 = no limit).
	MaxValueLength int
}

// SanitizeResult holds the outcome of a Sanitize call.
type SanitizeResult struct {
	Output   map[string]string
	Changes  []SanitizeChange
}

// SanitizeChange records a single key whose value was modified.
type SanitizeChange struct {
	Key    string
	Before string
	After  string
	Reason string
}

// Sanitize cleans env values according to the provided options.
// It does not mutate the input map.
func Sanitize(env map[string]string, opts SanitizeOptions) SanitizeResult {
	out := make(map[string]string, len(env))
	var changes []SanitizeChange

	for k, v := range env {
		original := v
		var reasons []string

		if opts.NormalizeLineEndings {
			v2 := strings.ReplaceAll(v, "\r\n", "\n")
			v2 = strings.ReplaceAll(v2, "\r", "\n")
			if v2 != v {
				v = v2
				reasons = append(reasons, "normalized line endings")
			}
		}

		if opts.StripControlChars {
			v2 := stripControl(v)
			if v2 != v {
				v = v2
				reasons = append(reasons, "stripped control characters")
			}
		}

		if opts.TrimQuotes {
			v2 := trimSurroundingQuotes(v)
			if v2 != v {
				v = v2
				reasons = append(reasons, "trimmed quotes")
			}
		}

		if opts.MaxValueLength > 0 && len(v) > opts.MaxValueLength {
			v = v[:opts.MaxValueLength]
			reasons = append(reasons, fmt.Sprintf("truncated to %d chars", opts.MaxValueLength))
		}

		out[k] = v
		if len(reasons) > 0 {
			changes = append(changes, SanitizeChange{
				Key:    k,
				Before: original,
				After:  v,
				Reason: strings.Join(reasons, "; "),
			})
		}
	}

	return SanitizeResult{Output: out, Changes: changes}
}

func stripControl(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r == '\t' || r == '\n' || (r >= 0x20 && r != 0x7f) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func trimSurroundingQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

// FormatSanitizeResult returns a human-readable summary of the sanitize operation.
func FormatSanitizeResult(r SanitizeResult) string {
	if len(r.Changes) == 0 {
		return "sanitize: no changes\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "sanitize: %d key(s) modified\n", len(r.Changes))
	for _, c := range r.Changes {
		fmt.Fprintf(&sb, "  %s: %s\n", c.Key, c.Reason)
	}
	return sb.String()
}
