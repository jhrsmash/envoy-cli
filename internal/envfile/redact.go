package envfile

import (
	"regexp"
	"strings"
)

// defaultSensitivePatterns matches common sensitive key names.
var defaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|api_key|apikey)`),
	regexp.MustCompile(`(?i)(private_key|privatekey)`),
	regexp.MustCompile(`(?i)(auth|credential)`),
}

const redactedValue = "***REDACTED***"

// RedactOptions controls redaction behaviour.
type RedactOptions struct {
	// ExtraPatterns are additional regexp patterns to match sensitive keys.
	ExtraPatterns []*regexp.Regexp
	// Keys is an explicit list of keys to redact (case-insensitive).
	Keys []string
}

// Redact returns a copy of env with sensitive values replaced by a placeholder.
// Matching is performed against key names using defaultSensitivePatterns plus
// any patterns or explicit keys supplied via opts.
func Redact(env map[string]string, opts RedactOptions) map[string]string {
	explicit := make(map[string]struct{}, len(opts.Keys))
	for _, k := range opts.Keys {
		explicit[strings.ToUpper(k)] = struct{}{}
	}

	patterns := make([]*regexp.Regexp, 0, len(defaultSensitivePatterns)+len(opts.ExtraPatterns))
	patterns = append(patterns, defaultSensitivePatterns...)
	patterns = append(patterns, opts.ExtraPatterns...)

	result := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, explicit, patterns) {
			result[k] = redactedValue
		} else {
			result[k] = v
		}
	}
	return result
}

func isSensitive(key string, explicit map[string]struct{}, patterns []*regexp.Regexp) bool {
	if _, ok := explicit[strings.ToUpper(key)]; ok {
		return true
	}
	for _, p := range patterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}
