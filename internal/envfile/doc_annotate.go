// Package envfile provides the Annotate function for attaching inline
// comments to environment variable keys.
//
// # Overview
//
// Annotate accepts an env map, a map of key->comment strings, and options
// that control which keys are targeted and whether existing annotations are
// overwritten.
//
// # Usage
//
//	annotations := map[string]string{
//		"DATABASE_URL": "Primary database connection string",
//		"API_KEY":      "Third-party API key – keep secret",
//	}
//
//	result := envfile.Annotate(env, annotations, envfile.AnnotateOptions{
//		Overwrite: false,
//	})
//	fmt.Print(envfile.FormatAnnotateResult(result))
//
// # Export
//
// Use ExportAnnotations to persist the annotation map to a .env or JSON file.
package envfile
