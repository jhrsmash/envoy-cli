// Package envfile provides the Pin feature for envoy-cli.
//
// # Pin
//
// Pin captures a stable, point-in-time snapshot of selected keys from an env
// map. The resulting pinned map can be persisted (e.g. via SaveSnapshot) and
// later compared against a live environment to detect unintentional value drift.
//
// Basic usage:
//
//	env, _ := envfile.Parse(".env.production")
//	pinned, result := envfile.Pin(env, envfile.PinOptions{
//	    Keys: []string{"DB_HOST", "DB_PORT"},
//	})
//	fmt.Print(envfile.FormatPinResult(result))
//
// To pin all keys, omit PinOptions.Keys (or leave it nil/empty).
//
// The returned map is a fresh copy and will not reflect future mutations to the
// original env map.
//
// # Drift Detection
//
// After pinning, use CompareSnapshot to diff a previously saved snapshot
// against the current environment:
//
//	old, _ := envfile.LoadSnapshot(".env.pin")
//	current, _ := envfile.Parse(".env.production")
//	drift := envfile.CompareSnapshot(old, current)
//	if len(drift) > 0 {
//	    fmt.Print(envfile.FormatDrift(drift))
//	}
package envfile
