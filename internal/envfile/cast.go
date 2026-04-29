package envfile

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// CastOptions controls how Cast behaves.
type CastOptions struct {
	// Keys is an optional list of keys to cast. If empty, all keys are processed.
	Keys []string
	// TargetType is the type to cast values to: "string", "int", "float", "bool".
	TargetType string
	// Strict causes Cast to return an error if any value cannot be converted.
	Strict bool
}

// CastResult holds the outcome of a Cast operation.
type CastResult struct {
	Output  map[string]string
	Cast    []string
	Skipped []string
	Failed  []string
}

// Cast coerces values in env to a target type, normalising their string
// representation (e.g. "TRUE" -> "true", "42.0" -> "42").
func Cast(env map[string]string, opts CastOptions) (CastResult, error) {
	if opts.TargetType == "" {
		opts.TargetType = "string"
	}

	output := make(map[string]string, len(env))
	for k, v := range env {
		output[k] = v
	}

	keySet := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		keySet[k] = true
	}

	result := CastResult{Output: output}

	for _, k := range sortedEnvKeys(env) {
		if len(keySet) > 0 && !keySet[k] {
			result.Skipped = append(result.Skipped, k)
			continue
		}

		norm, err := castValue(env[k], opts.TargetType)
		if err != nil {
			result.Failed = append(result.Failed, k)
			if opts.Strict {
				return CastResult{}, fmt.Errorf("cast: key %q value %q: %w", k, env[k], err)
			}
			continue
		}
		output[k] = norm
		result.Cast = append(result.Cast, k)
	}

	return result, nil
}

func castValue(v, targetType string) (string, error) {
	switch targetType {
	case "bool":
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return "", fmt.Errorf("not a bool")
		}
		return strconv.FormatBool(b), nil
	case "int":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("not a number")
		}
		return strconv.FormatInt(int64(f), 10), nil
	case "float":
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("not a float")
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case "string":
		return strings.TrimSpace(v), nil
	default:
		return "", fmt.Errorf("unknown target type %q", targetType)
	}
}

// FormatCastResult returns a human-readable summary of a CastResult.
func FormatCastResult(r CastResult) string {
	var sb strings.Builder
	if len(r.Cast) == 0 && len(r.Failed) == 0 {
		sb.WriteString("cast: no keys modified\n")
		return sb.String()
	}
	sort.Strings(r.Cast)
	for _, k := range r.Cast {
		fmt.Fprintf(&sb, "  cast   %s = %s\n", k, r.Output[k])
	}
	sort.Strings(r.Failed)
	for _, k := range r.Failed {
		fmt.Fprintf(&sb, "  failed %s\n", k)
	}
	return sb.String()
}

func sortedEnvKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
