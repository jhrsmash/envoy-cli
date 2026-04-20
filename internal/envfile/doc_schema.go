// Package envfile provides schema validation for .env files.
//
// The schema module allows you to define expected keys, their required/optional
// status, and an optional type pattern (string, int, bool, url). This is useful
// for catching configuration drift early, especially across environments.
//
// # Defining a Schema
//
// A schema can be built programmatically:
//
//	schema := envfile.Schema{
//		Fields: []envfile.SchemaField{
//			{Key: "DATABASE_URL", Required: true, Pattern: "url"},
//			{Key: "PORT",         Required: true, Pattern: "int"},
//			{Key: "DEBUG",        Required: false, Pattern: "bool"},
//		},
//	}
//
// Or loaded from a plain-text schema file via LoadSchema:
//
//	schema, err := envfile.LoadSchema("schema.envschema")
//
// # Validating an Env Map
//
// Pass any map[string]string (e.g. from Parse) to ValidateSchema:
//
//	violations := envfile.ValidateSchema(env, schema)
//	fmt.Print(envfile.FormatSchemaViolations(violations))
package envfile
