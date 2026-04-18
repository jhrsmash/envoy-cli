package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// interpolatePattern matches ${VAR} and $VAR style references.
var interpolatePattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateError is returned when a referenced variable cannot be resolved.
type InterpolateError struct {
	Key string
	Missing string
}

func (e *InterpolateError) Error() string {
	return fmt.Sprintf("key %q references undefined variable %q", e.Key, e.Missing)
}

// Interpolate expands variable references within env values using the same map.
// References to undefined variables return an InterpolateError.
// Already-resolved keys are used as the source of truth; no external env lookup
// is performed.
func Interpolate(env map[string]string) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	for k, v := range result {
		expanded, err := expandValue(v, result)
		if err != nil {
			return nil, &InterpolateError{Key: k, Missing: err.Error()}
		}
		result[k] = expanded
	}
	return result, nil
}

func expandValue(val string, env map[string]string) (string, error) {
	var expandErr string
	expanded := interpolatePattern.ReplaceAllStringFunc(val, func(match string) string {
		if expandErr != "" {
			return match
		}
		name := strings.TrimPrefix(strings.Trim(match, "${}"), "$")
		// handle ${VAR} form
		name = strings.TrimSuffix(strings.TrimPrefix(name, "{"), "}")
		submatches := interpolatePattern.FindStringSubmatch(match)
		if len(submatches) == 3 {
			if submatches[1] != "" {
				name = submatches[1]
			} else {
				name = submatches[2]
			}
		}
		resolved, ok := env[name]
		if !ok {
			expandErr = name
			return match
		}
		return resolved
	})
	if expandErr != "" {
		return "", fmt.Errorf("%s", expandErr)
	}
	return expanded, nil
}
