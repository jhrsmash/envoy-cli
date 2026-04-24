package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// MaskOptions controls how values are masked in output.
type MaskOptions struct {
	// ShowLength reveals the original value length as repeated '*' chars.
	ShowLength bool
	// Placeholder replaces every masked value when ShowLength is false.
	// Defaults to "***" if empty.
	Placeholder string
	// ExtraKeys are additional keys to mask beyond the auto-detected ones.
	ExtraKeys []string
}

// MaskResult holds the masked environment map and metadata.
type MaskResult struct {
	Masked  map[string]string
	Keys    []string // keys that were masked, sorted
}

// Mask returns a copy of env where sensitive values are replaced with a
// placeholder or a length-preserving '*' sequence. It relies on isSensitive
// from redact.go for automatic detection and respects ExtraKeys.
func Mask(env map[string]string, opts MaskOptions) MaskResult {
	placeholder := opts.Placeholder
	if placeholder == "" {
		placeholder = "***"
	}

	sensitiveSet := make(map[string]bool)
	for _, k := range opts.ExtraKeys {
		sensitiveSet[strings.ToUpper(k)] = true
	}
	for k := range env {
		if isSensitive(k) {
			sensitiveSet[k] = true
		}
	}

	masked := make(map[string]string, len(env))
	var maskedKeys []string

	for k, v := range env {
		if sensitiveSet[k] || sensitiveSet[strings.ToUpper(k)] {
			if opts.ShowLength && len(v) > 0 {
				masked[k] = strings.Repeat("*", len(v))
			} else {
				masked[k] = placeholder
			}
			maskedKeys = append(maskedKeys, k)
		} else {
			masked[k] = v
		}
	}

	sort.Strings(maskedKeys)
	return MaskResult{Masked: masked, Keys: maskedKeys}
}

// FormatMaskResult returns a human-readable summary of which keys were masked.
func FormatMaskResult(r MaskResult) string {
	if len(r.Keys) == 0 {
		return "No keys masked.\n"
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Masked %d key(s):\n", len(r.Keys))
	for _, k := range r.Keys {
		fmt.Fprintf(&sb, "  - %s\n", k)
	}
	return sb.String()
}
