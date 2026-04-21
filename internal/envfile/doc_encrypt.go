// Package envfile provides encryption and decryption utilities for
// env maps, allowing sensitive .env files to be stored securely at rest.
//
// # Overview
//
// Encrypt converts a map[string]string into an [EncryptedEnv] value
// using AES-256-GCM. The key is derived from a user-supplied passphrase
// via SHA-256. Each call produces a fresh random nonce, so repeated
// encryptions of the same data produce different ciphertexts.
//
// Decrypt reverses the process. Authentication is built into GCM, so
// any tampering or wrong passphrase returns an error immediately.
//
// # File helpers
//
// [SaveEncrypted] serialises an [EncryptedEnv] to a JSON file with
// permissions 0600. [LoadEncrypted] reads and decrypts such a file.
//
// # Example
//
//	enc, err := envfile.Encrypt(env, passphrase)
//	if err != nil { ... }
//	if err := envfile.SaveEncrypted(".env.enc", env, passphrase); err != nil { ... }
//
//	env, err := envfile.LoadEncrypted(".env.enc", passphrase)
package envfile
