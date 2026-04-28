// Package envfile provides the Prune function for removing unwanted keys
// from an environment map based on configurable criteria.
//
// # Overview
//
// Prune accepts a map of environment variables and a PruneOptions struct.
// It returns a PruneResult containing the cleaned map and the list of
// keys that were removed. The original map is never modified.
//
// # Options
//
//   - RemoveEmpty: drop keys whose value is an empty string.
//   - RemoveKeys: explicitly name keys to drop.
//   - RemovePrefix: drop every key that starts with the given prefix.
//
// Options may be combined freely; a key matching any criterion is removed.
//
// # Example
//
//	opts := envfile.PruneOptions{
//		RemoveEmpty:  true,
//		RemovePrefix: "LEGACY_",
//	}
//	result := envfile.Prune(env, opts)
//	fmt.Print(envfile.FormatPruneResult(result))
package envfile
