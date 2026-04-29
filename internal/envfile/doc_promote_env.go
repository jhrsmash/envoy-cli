// Package envfile provides utilities for managing .env files.
//
// # OS Environment Override
//
// OverrideFromOS applies values from the current process environment (os.Environ)
// onto an existing env map. This is useful for runtime injection of secrets or
// environment-specific configuration without modifying .env files on disk.
//
// # Usage
//
//	// Override any key in env that has a matching OS variable prefixed with "APP_".
//	result := envfile.OverrideFromOS(env, "APP_", false)
//	fmt.Print(envfile.FormatOverrideResult(result))
//
// # Prefix Stripping
//
// When a prefix is provided, it is stripped from the OS variable name before
// matching against the env map. For example, an OS variable "APP_PORT" with
// prefix "APP_" will match the env key "PORT".
//
// # allowNew Flag
//
// When allowNew is false (the default), only keys already present in the base
// env map can be overridden. Unknown OS variables are tracked in the Skipped
// field of the result. Set allowNew to true to permit injection of new keys.
package envfile
