// Package envfile provides snapshot functionality for capturing and comparing
// .env file states over time.
//
// Snapshots allow users to save a labeled, timestamped copy of an env map to
// disk as JSON, then later load and diff two snapshots to see what changed
// between environments or deployments.
//
// Example usage:
//
//	env, _ := envfile.Parse("production.env")
//	snap := envfile.NewSnapshot("prod-2024-06-01", env)
//	envfile.SaveSnapshot("snapshots/prod.json", snap)
//
//	oldSnap, _ := envfile.LoadSnapshot("snapshots/prod-old.json")
//	newSnap, _ := envfile.LoadSnapshot("snapshots/prod.json")
//	result := envfile.DiffSnapshots(oldSnap, newSnap)
//	fmt.Print(envfile.FormatSnapshotDiff(result))
package envfile
