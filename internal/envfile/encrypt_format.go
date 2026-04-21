package envfile

import "fmt"

// FormatEncryptSummary returns a human-readable summary of an encryption
// operation, suitable for CLI output.
func FormatEncryptSummary(path string, env map[string]string, redacted bool) string {
	keys := sortedKeys(env)
	var sensitive int
	for _, k := range keys {
		if isSensitive(k) {
			sensitive++
		}
	}

	lines := fmt.Sprintf("Encrypted %d key(s) → %s\n", len(keys), path)
	if sensitive > 0 {
		lines += fmt.Sprintf("  %d sensitive key(s) detected and protected\n", sensitive)
	}
	if redacted {
		lines += "  Values redacted from output\n"
	}
	return lines
}

// FormatDecryptSummary returns a human-readable summary shown after
// successfully decrypting a file.
func FormatDecryptSummary(path string, env map[string]string) string {
	return fmt.Sprintf("Decrypted %s — %d key(s) loaded\n", path, len(env))
}
