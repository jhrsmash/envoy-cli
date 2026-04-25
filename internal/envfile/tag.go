package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// TagOptions controls how keys are tagged.
type TagOptions struct {
	// Tags is a map of tag name -> list of keys to tag.
	Tags map[string][]string
	// Overwrite replaces existing tags on a key if true.
	Overwrite bool
}

// TagResult holds the outcome of a Tag operation.
type TagResult struct {
	// Tagged maps each key to its assigned tag.
	Tagged map[string]string
	// Skipped holds keys that already had a tag and Overwrite was false.
	Skipped []string
	// Missing holds keys referenced in Tags that were not found in the env.
	Missing []string
}

// Tag assigns string labels to env keys without modifying their values.
// The returned TagResult describes what was tagged, skipped, or missing.
func Tag(env map[string]string, existing map[string]string, opts TagOptions) TagResult {
	result := TagResult{
		Tagged: make(map[string]string),
	}

	// Copy existing tags so we don't mutate the caller's map.
	for k, v := range existing {
		result.Tagged[k] = v
	}

	// Build a reverse index: key -> tag from opts.Tags.
	keyToTag := make(map[string]string)
	for tag, keys := range opts.Tags {
		for _, k := range keys {
			keyToTag[k] = tag
		}
	}

	for key, tag := range keyToTag {
		if _, exists := env[key]; !exists {
			result.Missing = append(result.Missing, key)
			continue
		}
		if _, alreadyTagged := result.Tagged[key]; alreadyTagged && !opts.Overwrite {
			result.Skipped = append(result.Skipped, key)
			continue
		}
		result.Tagged[key] = tag
	}

	sort.Strings(result.Skipped)
	sort.Strings(result.Missing)
	return result
}

// FormatTagResult returns a human-readable summary of a TagResult.
func FormatTagResult(r TagResult) string {
	var sb strings.Builder

	if len(r.Tagged) == 0 && len(r.Skipped) == 0 && len(r.Missing) == 0 {
		return "No tags applied.\n"
	}

	if len(r.Tagged) > 0 {
		keys := make([]string, 0, len(r.Tagged))
		for k := range r.Tagged {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("Tagged (%d):\n", len(r.Tagged)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s = %s\n", k, r.Tagged[k]))
		}
	}

	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("Skipped (%d, already tagged):\n", len(r.Skipped)))
		for _, k := range r.Skipped {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	if len(r.Missing) > 0 {
		sb.WriteString(fmt.Sprintf("Missing (%d, not in env):\n", len(r.Missing)))
		for _, k := range r.Missing {
			sb.WriteString(fmt.Sprintf("  %s\n", k))
		}
	}

	return sb.String()
}
