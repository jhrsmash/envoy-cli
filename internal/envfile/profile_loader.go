package envfile

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProfileLoadOptions describes how to discover and load profile files from disk.
type ProfileLoadOptions struct {
	// BaseFile is the path to the base .env file.
	BaseFile string
	// ProfileNames lists the profile suffixes to load, e.g. ["staging", "local"].
	// Each name resolves to <dir>/.env.<name> relative to BaseFile.
	ProfileNames []string
	// Overwrite controls merge precedence (passed through to Profile).
	Overwrite bool
	// SkipMissing silently ignores profile files that do not exist.
	SkipMissing bool
}

// LoadProfile reads the base file and each named profile file from disk,
// then merges them using Profile.
func LoadProfile(opts ProfileLoadOptions) (ProfileResult, error) {
	base, err := Parse(opts.BaseFile)
	if err != nil {
		return ProfileResult{}, fmt.Errorf("profile: loading base file %q: %w", opts.BaseFile, err)
	}

	dir := filepath.Dir(opts.BaseFile)
	baseExt := filepath.Base(opts.BaseFile) // e.g. ".env"

	var profiles []map[string]string
	var resolvedNames []string

	for _, name := range opts.ProfileNames {
		path := filepath.Join(dir, baseExt+"."+name)
		env, err := Parse(path)
		if err != nil {
			if os.IsNotExist(err) && opts.SkipMissing {
				continue
			}
			return ProfileResult{}, fmt.Errorf("profile: loading profile %q from %q: %w", name, path, err)
		}
		profiles = append(profiles, env)
		resolvedNames = append(resolvedNames, name)
	}

	return Profile(ProfileOptions{
		Base:         base,
		Profiles:     profiles,
		ProfileNames: resolvedNames,
		Overwrite:    opts.Overwrite,
	})
}
