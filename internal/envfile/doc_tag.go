// Package envfile provides the Tag function for assigning string labels
// (tags) to environment variable keys.
//
// # Overview
//
// Tags are metadata labels that can be attached to env keys to describe their
// purpose, ownership, or sensitivity tier — without altering the key's value.
// Tags are stored as a separate map[string]string alongside the env map.
//
// # Usage
//
//	opts := envfile.TagOptions{
//		Tags: map[string][]string{
//			"database": {"DB_HOST", "DB_PORT", "DB_NAME"},
//			"security": {"API_KEY", "JWT_SECRET"},
//		},
//		Overwrite: false,
//	}
//	result := envfile.Tag(env, existingTags, opts)
//	fmt.Print(envfile.FormatTagResult(result))
//
// # Result
//
// TagResult contains:
//   - Tagged: the full tag map after applying opts (includes pre-existing tags).
//   - Skipped: keys that already had a tag and Overwrite was false.
//   - Missing: keys referenced in opts.Tags not found in the env map.
package envfile
