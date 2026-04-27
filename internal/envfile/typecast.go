package envfile

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// TypecastResult holds the outcome of a typecast operation.
type TypecastResult struct {
	Cast    map[string]string // keys that were successfully recast
	Skipped map[string]string // keys that could not be cast (value unchanged)
	Output  map[string]string // full resulting env map
}

// CastType represents the target type for casting.
type CastType string

const (
	CastBool   CastType = "bool"
	CastInt    CastType = "int"
	CastFloat  CastType = "float"
	CastString CastType = "string"
)

// TypecastOptions configures Typecast behaviour.
type TypecastOptions struct {
	// Keys to target; if empty, all keys are attempted.
	Keys []string
	// Target is the desired type to normalise values into.
	Target CastType
}

// Typecast attempts to normalise env values into a canonical string
// representation of the requested type (e.g. "true"/"false" for bool,
// decimal integers for int). Values that cannot be parsed are left
// unchanged and reported in Skipped.
func Typecast(env map[string]string, opts TypecastOptions) TypecastResult {
	output := make(map[string]string, len(env))
	for k, v := range env {
		output[k] = v
	}

	cast := make(map[string]string)
	skipped := make(map[string]string)

	targetKeys := opts.Keys
	if len(targetKeys) == 0 {
		for k := range env {
			targetKeys = append(targetKeys, k)
		}
	}

	for _, k := range targetKeys {
		v, ok := env[k]
		if !ok {
			continue
		}
		normalised, err := normaliseValue(v, opts.Target)
		if err != nil {
			skipped[k] = v
			continue
		}
		if normalised != v {
			cast[k] = normalised
			output[k] = normalised
		}
	}

	return TypecastResult{
		Cast:    cast,
		Skipped: skipped,
		Output:  output,
	}
}

func normaliseValue(v string, t CastType) (string, error) {
	switch t {
	case CastBool:
		b, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as bool", v)
		}
		if b {
			return "true", nil
		}
		return "false", nil
	case CastInt:
		i, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as int", v)
		}
		return strconv.FormatInt(i, 10), nil
	case CastFloat:
		f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return "", fmt.Errorf("cannot parse %q as float", v)
		}
		return strconv.FormatFloat(f, 'f', -1, 64), nil
	case CastString:
		return strings.TrimSpace(v), nil
	default:
		return "", fmt.Errorf("unknown cast type: %s", t)
	}
}

// FormatTypecastResult returns a human-readable summary of a TypecastResult.
func FormatTypecastResult(r TypecastResult) string {
	var sb strings.Builder

	if len(r.Cast) == 0 && len(r.Skipped) == 0 {
		sb.WriteString("typecast: no changes\n")
		return sb.String()
	}

	if len(r.Cast) > 0 {
		keys := make([]string, 0, len(r.Cast))
		for k := range r.Cast {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("cast (%d):\n", len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s = %s\n", k, r.Cast[k]))
		}
	}

	if len(r.Skipped) > 0 {
		keys := make([]string, 0, len(r.Skipped))
		for k := range r.Skipped {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		sb.WriteString(fmt.Sprintf("skipped (%d):\n", len(keys)))
		for _, k := range keys {
			sb.WriteString(fmt.Sprintf("  %s = %s\n", k, r.Skipped[k]))
		}
	}

	return sb.String()
}
