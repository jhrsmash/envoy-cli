// Package envfile provides the Search feature for querying .env maps.
//
// # Search
//
// Search scans an env map and returns all entries whose keys or values
// match a given query string or regular expression.
//
// Basic usage:
//
//	result, err := envfile.Search(env, envfile.SearchOptions{
//	    Query:        "DB",
//	    SearchKeys:   true,
//	    SearchValues: false,
//	})
//
// Regex usage:
//
//	result, err := envfile.Search(env, envfile.SearchOptions{
//	    Query:    `^APP_`,
//	    UseRegex: true,
//	})
//
// Results can be formatted with FormatSearchResult or exported to a file
// via ExportSearchResult (supports "json" and "text" formats).
package envfile
