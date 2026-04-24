package envfile

import (
	"fmt"
	"strings"
)

// TransformOp represents a transformation operation to apply to env values.
type TransformOp string

const (
	TransformUppercase  TransformOp = "uppercase"
	TransformLowercase  TransformOp = "lowercase"
	TransformTrimPrefix TransformOp = "trim_prefix"
	TransformTrimSuffix TransformOp = "trim_suffix"
	TransformReplace    TransformOp = "replace"
)

// TransformOptions configures a Transform call.
type TransformOptions struct {
	// Op is the operation to apply.
	Op TransformOp
	// Keys restricts transformation to specific keys; empty means all keys.
	Keys []string
	// Arg1 is the first string argument (e.g. prefix to trim, old value for replace).
	Arg1 string
	// Arg2 is the second string argument (e.g. new value for replace).
	Arg2 string
}

// TransformResult holds the outcome of a Transform call.
type TransformResult struct {
	Output  map[string]string
	Changed []string
}

// Transform applies a value transformation to an env map and returns a new map.
func Transform(env map[string]string, opts TransformOptions) (TransformResult, error) {
	target := make(map[string]bool, len(opts.Keys))
	for _, k := range opts.Keys {
		target[k] = true
	}

	out := make(map[string]string, len(env))
	var changed []string

	for k, v := range env {
		if len(target) > 0 && !target[k] {
			out[k] = v
			continue
		}

		var newVal string
		switch opts.Op {
		case TransformUppercase:
			newVal = strings.ToUpper(v)
		case TransformLowercase:
			newVal = strings.ToLower(v)
		case TransformTrimPrefix:
			newVal = strings.TrimPrefix(v, opts.Arg1)
		case TransformTrimSuffix:
			newVal = strings.TrimSuffix(v, opts.Arg1)
		case TransformReplace:
			newVal = strings.ReplaceAll(v, opts.Arg1, opts.Arg2)
		default:
			return TransformResult{}, fmt.Errorf("unknown transform op: %q", opts.Op)
		}

		out[k] = newVal
		if newVal != v {
			changed = append(changed, k)
		}
	}

	return TransformResult{Output: out, Changed: changed}, nil
}
