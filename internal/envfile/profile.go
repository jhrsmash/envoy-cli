package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// ProfileOptions controls how a profile merge is performed.
type ProfileOptions struct {
	// Base is the base environment (e.g. .env)
	Base map[string]string
	// Profiles is an ordered list of profile overlays to apply (e.g. .env.staging, .env.local)
	Profiles []map[string]string
	// ProfileNames are human-readable labels for each profile (same order as Profiles)
	ProfileNames []string
	// Overwrite controls whether later profiles overwrite earlier values
	Overwrite bool
}

// ProfileEntry records the resolved value of a key and which layer it came from.
type ProfileEntry struct {
	Key    string
	Value  string
	Source string // "base" or the profile name
}

// ProfileResult holds the merged environment and per-key provenance.
type ProfileResult struct {
	Env     map[string]string
	Entries []ProfileEntry
}

// Profile merges a base environment with one or more named profile overlays.
// Later profiles take precedence when Overwrite is true.
func Profile(opts ProfileOptions) (ProfileResult, error) {
	if len(opts.ProfileNames) > 0 && len(opts.ProfileNames) != len(opts.Profiles) {
		return ProfileResult{}, fmt.Errorf(
			"profile: ProfileNames length (%d) must match Profiles length (%d)",
			len(opts.ProfileNames), len(opts.Profiles),
		)
	}

	merged := make(map[string]string, len(opts.Base))
	source := make(map[string]string, len(opts.Base))

	for k, v := range opts.Base {
		merged[k] = v
		source[k] = "base"
	}

	for i, profile := range opts.Profiles {
		name := fmt.Sprintf("profile[%d]", i)
		if i < len(opts.ProfileNames) && opts.ProfileNames[i] != "" {
			name = opts.ProfileNames[i]
		}
		for k, v := range profile {
			if _, exists := merged[k]; !exists || opts.Overwrite {
				merged[k] = v
				source[k] = name
			}
		}
	}

	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	entries := make([]ProfileEntry, 0, len(keys))
	for _, k := range keys {
		entries = append(entries, ProfileEntry{
			Key:    k,
			Value:  merged[k],
			Source: source[k],
		})
	}

	return ProfileResult{Env: merged, Entries: entries}, nil
}

// FormatProfileResult returns a human-readable summary of the profile merge.
func FormatProfileResult(r ProfileResult) string {
	if len(r.Entries) == 0 {
		return "profile: no keys resolved\n"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("profile: %d keys resolved\n", len(r.Entries)))
	for _, e := range r.Entries {
		sb.WriteString(fmt.Sprintf("  %-30s <- %s\n", e.Key, e.Source))
	}
	return sb.String()
}
