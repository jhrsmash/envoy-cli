// Package envfile provides the Scope function for narrowing an environment map
// to a specific set of key prefixes.
//
// # Overview
//
// Scope is useful when a single .env file contains keys for multiple services
// or subsystems (e.g. APP_, DB_, CACHE_) and you need to extract only the
// relevant subset for a given component.
//
// # Basic Usage
//
//	result := envfile.Scope(env, envfile.ScopeOptions{
//		Prefixes:    []string{"APP_"},
//		StripPrefix: true,
//	})
//	// result.Scoped contains keys without the APP_ prefix
//
// # Options
//
//   - Prefixes      — one or more key prefixes to include
//   - StripPrefix   — remove the matched prefix from output keys
//   - Uppercase     — convert output keys to uppercase
package envfile
