package envfile

import "fmt"

// CloneResult holds the outcome of a clone operation.
type CloneResult struct {
	Source      string
	Destination string
	Keys        []string
	Redacted    []string
}

// CloneOptions configures the behaviour of Clone.
type CloneOptions struct {
	// RedactSensitive replaces sensitive values with a placeholder in the destination.
	RedactSensitive bool
	// ExtraRedactPatterns adds additional key patterns to redact.
	ExtraRedactPatterns []string
	// OverwriteExisting controls whether existing keys in dst are overwritten.
	OverwriteExisting bool
}

// Clone copies all key/value pairs from src into a new map, optionally
// redacting sensitive values before returning. It never mutates src.
func Clone(src map[string]string, opts CloneOptions) (map[string]string, CloneResult) {
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}

	var redactedKeys []string
	if opts.RedactSensitive {
		redacted := Redact(dst, opts.ExtraRedactPatterns)
		redactedKeys = RedactedKeys(dst, opts.ExtraRedactPatterns)
		dst = redacted
	}

	keys := sortedKeys(dst)

	return dst, CloneResult{
		Keys:     keys,
		Redacted: redactedKeys,
	}
}

// FormatCloneResult returns a human-readable summary of a clone operation.
func FormatCloneResult(r CloneResult) string {
	out := fmt.Sprintf("Cloned %d key(s)", len(r.Keys))
	if r.Source != "" || r.Destination != "" {
		out += fmt.Sprintf(" from %q to %q", r.Source, r.Destination)
	}
	if len(r.Redacted) > 0 {
		out += fmt.Sprintf("\nRedacted %d sensitive key(s):", len(r.Redacted))
		for _, k := range r.Redacted {
			out += fmt.Sprintf("\n  - %s", k)
		}
	}
	return out
}
