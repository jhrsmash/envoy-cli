// Package envfile provides the Extract function for pulling a targeted subset
// of keys from an env map.
//
// # Overview
//
// Extract is useful when you need to isolate a group of related variables from
// a larger env file — for example, pulling all APP_* keys before passing them
// to a subprocess, or stripping a deployment prefix before writing a scoped
// config file.
//
// # Basic usage
//
//	result, err := envfile.Extract(env, envfile.ExtractOptions{
//	    Prefix:      "APP_",
//	    StripPrefix: true,
//	})
//
// # Explicit key selection
//
//	result, err := envfile.Extract(env, envfile.ExtractOptions{
//	    Keys:          []string{"DATABASE_URL", "REDIS_URL"},
//	    FailOnMissing: true,
//	})
//
// When FailOnMissing is true, Extract returns an error listing every key that
// could not be found in the source map. When false, missing keys are recorded
// in ExtractResult.Missing but no error is returned.
package envfile
