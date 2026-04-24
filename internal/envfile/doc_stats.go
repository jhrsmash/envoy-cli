// Package envfile provides utilities for managing .env files.
//
// # Stats
//
// The Stats feature computes aggregate metrics about an env map, useful for
// auditing, reporting, and understanding the shape of a configuration set.
//
// Usage:
//
//	env, _ := envfile.Parse(".env")
//	stats := envfile.ComputeStats(env)
//	fmt.Print(envfile.FormatStats(stats))
//
// Metrics reported:
//   - Total number of keys
//   - Number of keys with empty values
//   - Number of keys detected as sensitive (passwords, tokens, secrets)
//   - Average value length in characters
//   - Key counts grouped by underscore-separated prefix (e.g. DB_*, APP_*)
package envfile
