// Package envfile provides the Patch function for applying a sequence of
// structured operations to an env map.
//
// # Patch Operations
//
// A patch is a slice of PatchOp values, each describing one of three actions:
//
//   - "set"    — create or overwrite a key with a given value
//   - "delete" — remove a key (skipped with an error if the key is absent)
//   - "rename" — move a key to a new name, preserving its value
//
// Patch never mutates the input map; it always returns a new map.
//
// # Example
//
//	ops := []envfile.PatchOp{
//		{Action: "set",    Key: "APP_ENV",  Value: "production"},
//		{Action: "delete", Key: "DEBUG"},
//		{Action: "rename", Key: "APP_NAME", NewKey: "SERVICE_NAME"},
//	}
//	result, summary := envfile.Patch(env, ops)
//	fmt.Print(envfile.FormatPatchResult(summary))
package envfile
