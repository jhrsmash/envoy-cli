// Package envfile provides the Rename function for safely renaming keys
// within an environment map.
//
// # Overview
//
// Rename copies the supplied env map, removes oldKey, and inserts newKey with
// the same value. The original map is never mutated.
//
// # Options
//
// RenameOptions.Overwrite controls whether an existing newKey is replaced.
// When false (the default) the operation is skipped and the reason is
// recorded in RenameResult.
//
// # Example
//
//	env := map[string]string{"DB_HOST": "localhost"}
//	out, result := envfile.Rename(env, "DB_HOST", "DATABASE_HOST", envfile.RenameOptions{})
//	if result.Renamed {
//		fmt.Println(envfile.FormatRenameResult(result))
//	}
package envfile
