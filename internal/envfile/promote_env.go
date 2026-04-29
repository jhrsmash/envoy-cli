package envfile

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// EnvOverride represents a single environment variable override
// sourced from the process environment (os.Environ).
type EnvOverride struct {
	Key      string
	Value    string
	Previous string
	WasSet   bool
}

// EnvOverrideResult holds the outcome of applying OS environment overrides.
type EnvOverrideResult struct {
	Applied  []EnvOverride
	Skipped  []string
	Env      map[string]string
}

// OverrideFromOS merges values from the current process environment into env.
// Only keys already present in env are eligible for override unless allowNew is true.
// If prefix is non-empty, only OS vars matching that prefix are considered,
// and the prefix is stripped before matching against env keys.
func OverrideFromOS(env map[string]string, prefix string, allowNew bool) EnvOverrideResult {
	result := EnvOverrideResult{
		Env: make(map[string]string, len(env)),
	}
	for k, v := range env {
		result.Env[k] = v
	}

	for _, raw := range os.Environ() {
		idx := strings.IndexByte(raw, '=')
		if idx < 0 {
			continue
		}
		osKey := raw[:idx]
		osVal := raw[idx+1:]

		envKey := osKey
		if prefix != "" {
			if !strings.HasPrefix(osKey, prefix) {
				continue
			}
			envKey = osKey[len(prefix):]
		}

		prev, exists := result.Env[envKey]
		if !exists && !allowNew {
			result.Skipped = append(result.Skipped, envKey)
			continue
		}

		result.Env[envKey] = osVal
		result.Applied = append(result.Applied, EnvOverride{
			Key:      envKey,
			Value:    osVal,
			Previous: prev,
			WasSet:   exists,
		})
	}

	sort.Slice(result.Applied, func(i, j int) bool {
		return result.Applied[i].Key < result.Applied[j].Key
	})
	sort.Strings(result.Skipped)
	return result
}

// FormatOverrideResult returns a human-readable summary of an EnvOverrideResult.
func FormatOverrideResult(r EnvOverrideResult) string {
	var sb strings.Builder
	if len(r.Applied) == 0 && len(r.Skipped) == 0 {
		sb.WriteString("no OS overrides applied\n")
		return sb.String()
	}
	if len(r.Applied) > 0 {
		sb.WriteString(fmt.Sprintf("applied %d override(s):\n", len(r.Applied)))
		for _, ov := range r.Applied {
			if ov.WasSet {
				sb.WriteString(fmt.Sprintf("  ~ %s: %q -> %q\n", ov.Key, ov.Previous, ov.Value))
			} else {
				sb.WriteString(fmt.Sprintf("  + %s = %q\n", ov.Key, ov.Value))
			}
		}
	}
	if len(r.Skipped) > 0 {
		sb.WriteString(fmt.Sprintf("skipped %d key(s) (not in base env):\n", len(r.Skipped)))
		for _, k := range r.Skipped {
			sb.WriteString(fmt.Sprintf("  - %s\n", k))
		}
	}
	return sb.String()
}
