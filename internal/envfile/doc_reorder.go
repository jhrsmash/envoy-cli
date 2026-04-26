// Package envfile — reorder module
//
// Reorder allows callers to define an explicit key ordering for an env map.
// This is useful when generating .env files that must follow a canonical
// layout (e.g. grouped by service, with secrets last).
//
// # Basic usage
//
//	out, result := envfile.Reorder(env, envfile.ReorderOptions{
//		Keys:      []string{"APP_NAME", "APP_VERSION", "DB_HOST"},
//		AlphaTail: true,
//	})
//
// # Exporting
//
// ExportReordered serialises the map in the resolved order:
//
//	contents, err := envfile.ExportReordered(out, result, envfile.ReorderFormatDotenv)
//
// ExportReorderedToFile writes directly to a file:
//
//	err := envfile.ExportReorderedToFile(".env.ordered", out, result, envfile.ReorderFormatJSON)
//
// # Notes
//
//   - Keys listed in ReorderOptions.Keys that are absent from the env map are
//     reported in ReorderResult.Missing and silently skipped.
//   - Keys present in the env map but not listed in ReorderOptions.Keys appear
//     in ReorderResult.Unlisted and are appended after the explicit list.
//   - When AlphaTail is true, unlisted keys are sorted alphabetically before
//     being appended.
package envfile
