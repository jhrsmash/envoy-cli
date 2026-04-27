// Package envfile provides the Profile feature for layered environment
// composition across deployment profiles.
//
// # Overview
//
// Many projects maintain a base .env file alongside profile-specific overlays
// such as .env.staging or .env.local. The Profile feature merges these layers
// in a defined order, recording which layer each key originated from.
//
// # Usage
//
//	result, err := envfile.Profile(envfile.ProfileOptions{
//		Base:         baseEnv,
//		Profiles:     []map[string]string{stagingEnv, localEnv},
//		ProfileNames: []string{"staging", "local"},
//		Overwrite:    true,
//	})
//
// # Loading from disk
//
//	result, err := envfile.LoadProfile(envfile.ProfileLoadOptions{
//		BaseFile:     ".env",
//		ProfileNames: []string{"staging", "local"},
//		Overwrite:    true,
//		SkipMissing:  true,
//	})
//
// # Provenance
//
// Each key in ProfileResult.Entries carries a Source field indicating whether
// the value came from "base" or a named profile, making it easy to audit which
// layer contributed each configuration value.
package envfile
