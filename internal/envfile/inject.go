package envfile

import (
	"fmt"
	"os"
	"sort"
)

// InjectOptions controls how environment variables are injected into the process.
type InjectOptions struct {
	// Overwrite replaces existing OS environment variables with values from env.
	Overwrite bool
	// Keys restricts injection to only the specified keys. If empty, all keys are injected.
	Keys []string
}

// InjectResult summarises what was injected, skipped, or overwritten.
type InjectResult struct {
	Injected    []string
	Skipped     []string
	Overwritten []string
}

// Inject sets keys from env into the current process environment via os.Setenv.
// It returns an InjectResult describing what happened and any first error encountered.
func Inject(env map[string]string, opts InjectOptions) (InjectResult, error) {
	var result InjectResult

	target := env
	if len(opts.Keys) > 0 {
		target = make(map[string]string, len(opts.Keys))
		for _, k := range opts.Keys {
			if v, ok := env[k]; ok {
				target[k] = v
			}
		}
	}

	keys := make([]string, 0, len(target))
	for k := range target {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := target[k]
		_, exists := os.LookupEnv(k)
		if exists && !opts.Overwrite {
			result.Skipped = append(result.Skipped, k)
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			return result, fmt.Errorf("inject: failed to set %q: %w", k, err)
		}
		if exists {
			result.Overwritten = append(result.Overwritten, k)
		} else {
			result.Injected = append(result.Injected, k)
		}
	}

	return result, nil
}
