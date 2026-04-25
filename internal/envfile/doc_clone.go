// Package envfile provides utilities for managing .env files.
//
// # Clone
//
// The Clone function creates a deep copy of an env map, with optional
// redaction of sensitive values before the copy is handed to the caller.
//
// Basic usage:
//
//	dst, result := envfile.Clone(src, envfile.CloneOptions{
//		RedactSensitive: true,
//	})
//
// Extra patterns can be supplied to redact additional keys beyond the
// built-in heuristics (PASSWORD, SECRET, TOKEN, KEY, etc.):
//
//	dst, result := envfile.Clone(src, envfile.CloneOptions{
//		RedactSensitive:     true,
//		ExtraRedactPatterns: []string{"INTERNAL", "PRIVATE"},
//	})
//
// The original source map is never modified.
//
// # CloneResult
//
// The returned CloneResult value reports how many keys were copied and
// how many values were redacted during the operation:
//
//	dst, result := envfile.Clone(src, envfile.CloneOptions{RedactSensitive: true})
//	fmt.Printf("copied %d keys, redacted %d values\n", result.Copied, result.Redacted)
package envfile
