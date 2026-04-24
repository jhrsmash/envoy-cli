package envfile

import (
	"fmt"
	"sort"
)

// PinResult holds the outcome of a pin operation.
type PinResult struct {
	Pinned  []string // keys whose values were pinned (locked)
	Skipped []string // keys not found in the env map
}

// PinOptions controls the behaviour of Pin.
type PinOptions struct {
	// Keys is the explicit list of keys to pin. If empty, all keys are pinned.
	Keys []string
	// PinnedMark is the suffix appended as a comment-style annotation in the
	// returned pinned map. Defaults to "__PINNED".
	PinnedMark string
}

// Pin records a stable snapshot of selected keys from env into a new map that
// can be used to detect unintentional drift. The returned map contains only
// the pinned keys with their current values. Keys that do not exist in env are
// recorded in PinResult.Skipped.
func Pin(env map[string]string, opts PinOptions) (map[string]string, PinResult) {
	mark := opts.PinnedMark
	if mark == "" {
		mark = "__PINNED"
	}

	keys := opts.Keys
	if len(keys) == 0 {
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
	}

	pinned := make(map[string]string, len(keys))
	var result PinResult

	for _, k := range keys {
		v, ok := env[k]
		if !ok {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		pinned[k] = v
		result.Pinned = append(result.Pinned, k)
	}

	return pinned, result
}

// FormatPinResult returns a human-readable summary of a PinResult.
func FormatPinResult(r PinResult) string {
	out := fmt.Sprintf("Pinned: %d key(s)\n", len(r.Pinned))
	for _, k := range r.Pinned {
		out += fmt.Sprintf("  + %s\n", k)
	}
	if len(r.Skipped) > 0 {
		out += fmt.Sprintf("Skipped (not found): %d key(s)\n", len(r.Skipped))
		for _, k := range r.Skipped {
			out += fmt.Sprintf("  - %s\n", k)
		}
	}
	return out
}
