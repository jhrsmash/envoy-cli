// Package envfile provides utilities for parsing, diffing, merging,
// validating, and exporting .env files across environments.
//
// # Export
//
// The Export function converts an in-memory env map into a formatted string
// suitable for use in various contexts:
//
//   - FormatJSON   — structured JSON object, useful for APIs or config tools
//   - FormatDotenv — standard KEY="value" pairs, compatible with dotenv loaders
//   - FormatShell  — shell-ready `export KEY="value"` statements
//
// ExportOptions allows callers to control sorting and automatic redaction of
// sensitive values before output.
//
// Example:
//
//	out, err := envfile.Export(env, envfile.ExportOptions{
//		Format:   envfile.FormatJSON,
//		Sorted:   true,
//		Redacted: true,
//	})
package envfile
