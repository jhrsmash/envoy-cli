package envfile

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// EncryptedEnv holds the encrypted representation of an env map.
type EncryptedEnv struct {
	Ciphertext string `json:"ciphertext"`
	Nonce      string `json:"nonce"`
}

// deriveKey produces a 32-byte AES-256 key from a passphrase.
func deriveKey(passphrase string) []byte {
	sum := sha256.Sum256([]byte(passphrase))
	return sum[:]
}

// Encrypt serialises the env map to dotenv format and encrypts it
// with AES-256-GCM using the supplied passphrase.
func Encrypt(env map[string]string, passphrase string) (*EncryptedEnv, error) {
	if passphrase == "" {
		return nil, errors.New("encrypt: passphrase must not be empty")
	}

	plaintext := marshalDotenv(env)

	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)

	return &EncryptedEnv{
		Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
		Nonce:      base64.StdEncoding.EncodeToString(nonce),
	}, nil
}

// Decrypt reverses Encrypt and returns the original env map.
func Decrypt(enc *EncryptedEnv, passphrase string) (map[string]string, error) {
	if passphrase == "" {
		return nil, errors.New("decrypt: passphrase must not be empty")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(enc.Ciphertext)
	if err != nil {
		return nil, errors.New("decrypt: invalid ciphertext encoding")
	}

	nonce, err := base64.StdEncoding.DecodeString(enc.Nonce)
	if err != nil {
		return nil, errors.New("decrypt: invalid nonce encoding")
	}

	block, err := aes.NewCipher(deriveKey(passphrase))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decrypt: authentication failed — wrong passphrase or corrupted data")
	}

	return Parse(string(plaintext))
}
