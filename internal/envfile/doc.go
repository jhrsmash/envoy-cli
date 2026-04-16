// Package envfile provides utilities for parsing and representing .env files.
//
// A .env file is a plain-text file containing KEY=VALUE pairs, one per line.
// Lines beginning with '#' are treated as comments and ignored. Empty lines
// are also skipped. Values may optionally be wrapped in single or double quotes,
// which will be stripped during parsing.
//
// Example usage:
//
//	env, err := envfile.Parse(".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(env["DATABASE_URL"])
package envfile
