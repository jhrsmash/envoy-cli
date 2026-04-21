// Package envfile provides utilities for parsing, diffing, and managing .env files.
//
// # Promote
//
// The Promote function copies keys from a source environment map into a
// destination environment map, supporting selective promotion and overwrite control.
//
// Basic usage:
//
//	// Promote all keys from staging into production, skipping existing keys.
//	result, pr := envfile.Promote(staging, production, nil, false)
//
//	// Promote specific keys only.
//	result, pr := envfile.Promote(staging, production, []string{"DB_HOST", "DB_PORT"}, false)
//
//	// Promote and overwrite existing keys in destination.
//	result, pr := envfile.Promote(staging, production, nil, true)
//
// The returned PromoteResult categorises each key as Promoted, Skipped, or Missing.
// Use FormatPromoteResult to render a human-readable summary suitable for CLI output.
package envfile
