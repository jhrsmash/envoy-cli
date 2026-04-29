// Package envfile provides rollback support for .env environments.
//
// # Rollback
//
// Rollback restores a previously archived snapshot of an environment.
// It uses the Archive/LoadArchive primitives to locate saved states and
// computes a Diff between the current env and the target snapshot so the
// caller can present a clear summary of what will change.
//
// Usage:
//
//	archives, _ := envfile.LoadRollbackIndex(dir)
//	result, err := envfile.Rollback(current, archives, envfile.RollbackOptions{
//		Label:  "v1",
//		DryRun: false,
//	})
//	fmt.Println(envfile.FormatRollbackResult(result))
//
// # Persistence
//
// SaveRollbackIndex and LoadRollbackIndex persist the archive list as a
// JSON index file inside a directory of your choice (typically the same
// directory as your .env files).
package envfile
