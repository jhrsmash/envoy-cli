package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// marshalDotenv converts an env map to a dotenv-formatted string.
// Keys are emitted in sorted order for determinism.
func marshalDotenv(env map[string]string) string {
	keys := sortedKeys(env)
	var sb strings.Builder
	for _, k := range keys {
		fmt.Fprintf(&sb, "%s=%s\n", k, env[k])
	}
	return sb.String()
}

// SaveEncrypted encrypts env with passphrase and writes the result as
// a JSON file to path. The file is written with mode 0600 to restrict
// access to the owning user only.
func SaveEncrypted(path string, env map[string]string, passphrase string) error {
	enc, err := Encrypt(env, passphrase)
	if err != nil {
		return fmt.Errorf("save encrypted: %w", err)
	}

	data, err := json.MarshalIndent(enc, "", "  ")
	if err != nil {
		return fmt.Errorf("save encrypted: marshal: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("save encrypted: write %s: %w", path, err)
	}
	return nil
}

// LoadEncrypted reads a JSON-encoded EncryptedEnv from path and
// decrypts it using passphrase, returning the original env map.
func LoadEncrypted(path, passphrase string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load encrypted: read %s: %w", path, err)
	}

	var enc EncryptedEnv
	if err := json.Unmarshal(data, &enc); err != nil {
		return nil, fmt.Errorf("load encrypted: unmarshal: %w", err)
	}

	env, err := Decrypt(&enc, passphrase)
	if err != nil {
		return nil, fmt.Errorf("load encrypted: %w", err)
	}
	return env, nil
}

// RekeyEncrypted re-encrypts the env file at path using a new passphrase.
// It decrypts with oldPassphrase, then re-encrypts and overwrites the file
// with newPassphrase. The operation is atomic with respect to the passphrase
// change — the file is only overwritten if both steps succeed.
func RekeyEncrypted(path, oldPassphrase, newPassphrase string) error {
	env, err := LoadEncrypted(path, oldPassphrase)
	if err != nil {
		return fmt.Errorf("rekey encrypted: %w", err)
	}

	if err := SaveEncrypted(path, env, newPassphrase); err != nil {
		return fmt.Errorf("rekey encrypted: %w", err)
	}
	return nil
}
