package envfile

import (
	"fmt"
	"sort"
	"strings"
)

// Annotation holds a comment or label attached to a key.
type Annotation struct {
	Key     string
	Comment string
}

// AnnotateOptions controls how annotations are applied.
type AnnotateOptions struct {
	// Overwrite replaces existing annotations when true.
	Overwrite bool
	// Keys restricts annotation to specific keys; empty means all keys.
	Keys []string
}

// AnnotateResult holds the outcome of an Annotate operation.
type AnnotateResult struct {
	Annotations map[string]string // key -> comment
	Skipped     []string          // keys that already had annotations and Overwrite was false
	Missing     []string          // keys requested but not present in env
}

// Annotate attaches comments to keys in env according to the provided mapping.
// The annotations map is key -> comment string.
func Annotate(env map[string]string, annotations map[string]string, opts AnnotateOptions) AnnotateResult {
	targetKeys := map[string]bool{}
	for _, k := range opts.Keys {
		targetKeys[k] = true
	}

	applied := make(map[string]string)
	var skipped, missing []string

	for key, comment := range annotations {
		if len(opts.Keys) > 0 && !targetKeys[key] {
			continue
		}
		if _, exists := env[key]; !exists {
			missing = append(missing, key)
			continue
		}
		if _, already := applied[key]; already && !opts.Overwrite {
			skipped = append(skipped, key)
			continue
		}
		if comment == "" {
			continue
		}
		if !strings.HasPrefix(comment, "# ") {
			comment = fmt.Sprintf("# %s", comment)
		}
		applied[key] = comment
	}

	sort.Strings(skipped)
	sort.Strings(missing)

	return AnnotateResult{
		Annotations: applied,
		Skipped:     skipped,
		Missing:     missing,
	}
}
