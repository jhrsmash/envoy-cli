package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	env := map[string]string{
		"API_KEY":  "supersecret",
		"DB_PASS":  "hunter2",
		"APP_NAME": "envoy",
	}

	enc, err := Encrypt(env, "my-passphrase")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	got, err := Decrypt(enc, "my-passphrase")
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	for k, want := range env {
		if got[k] != want {
			t.Errorf("key %s: got %q, want %q", k, got[k], want)
		}
	}
}

func TestEncrypt_EmptyPassphrase(t *testing.T) {
	_, err := Encrypt(map[string]string{"K": "V"}, "")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestDecrypt_WrongPassphrase(t *testing.T) {
	enc, err := Encrypt(map[string]string{"K": "V"}, "correct")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	_, err = Decrypt(enc, "wrong")
	if err == nil {
		t.Fatal("expected error for wrong passphrase")
	}
}

func TestDecrypt_EmptyPassphrase(t *testing.T) {
	enc := &EncryptedEnv{Ciphertext: "abc", Nonce: "def"}
	_, err := Decrypt(enc, "")
	if err == nil {
		t.Fatal("expected error for empty passphrase")
	}
}

func TestSaveAndLoadEncrypted(t *testing.T) {
	env := map[string]string{
		"SECRET": "value123",
		"PORT":   "8080",
	}
	pass := "test-pass-phrase"

	dir := t.TempDir()
	path := filepath.Join(dir, "secrets.enc.json")

	if err := SaveEncrypted(path, env, pass); err != nil {
		t.Fatalf("SaveEncrypted: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}

	got, err := LoadEncrypted(path, pass)
	if err != nil {
		t.Fatalf("LoadEncrypted: %v", err)
	}

	for k, want := range env {
		if got[k] != want {
			t.Errorf("key %s: got %q, want %q", k, got[k], want)
		}
	}
}

func TestLoadEncrypted_MissingFile(t *testing.T) {
	_, err := LoadEncrypted("/nonexistent/path.json", "pass")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
