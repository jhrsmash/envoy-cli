// Package envfile provides utilities for parsing, diffing, merging,
// validating, and auditing .env files.
//
// # Audit
//
// The audit sub-feature records changes between two versions of an env map
// and persists them as a JSON log for later review.
//
// Basic usage:
//
//	before, _ := envfile.Parse("staging.env")
//	after, _  := envfile.Parse("production.env")
//
//	log := envfile.Audit(before, after, "production.env")
//	fmt.Print(envfile.FormatAuditLog(log))
//
//	// Persist for later review
//	_ = envfile.SaveAuditLog(".envoy-audit.json", log)
//
// Each AuditEntry captures the timestamp, key, action (added / removed /
// changed), old and new values, and an arbitrary source label so that
// multiple audit passes can be stored in the same file.
package envfile
