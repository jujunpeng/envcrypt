package sync

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/envcrypt/internal/crypto"
	"github.com/envcrypt/internal/envfile"
	"github.com/envcrypt/internal/keystore"
)

// Manager handles encrypting and decrypting .env files for team sync.
type Manager struct {
	encryptor *crypto.Encryptor
	keystore  *keystore.KeyStore
}

// New creates a new sync Manager with the provided encryptor and keystore.
func New(enc *crypto.Encryptor, ks *keystore.KeyStore) *Manager {
	return &Manager{encryptor: enc, keystore: ks}
}

// Encrypt reads a plaintext .env file, encrypts it for all recipients,
// and writes the ciphertext to outPath.
func (m *Manager) Encrypt(envPath, outPath string) error {
	env, err := envfile.Parse(envPath)
	if err != nil {
		return fmt.Errorf("sync: parse env file: %w", err)
	}

	recipients := m.keystore.AllRecipients()
	if len(recipients) == 0 {
		return fmt.Errorf("sync: no recipients configured")
	}

	plaintext := env.Serialize()
	ciphertext, err := m.encryptor.Encrypt([]byte(plaintext), recipients)
	if err != nil {
		return fmt.Errorf("sync: encrypt: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("sync: create output dir: %w", err)
	}

	if err := os.WriteFile(outPath, ciphertext, 0o600); err != nil {
		return fmt.Errorf("sync: write encrypted file: %w", err)
	}

	return nil
}

// Decrypt reads an encrypted .env file, decrypts it using the encryptor's
// identity, and writes the plaintext to outPath.
func (m *Manager) Decrypt(encPath, outPath string) error {
	ciphertext, err := os.ReadFile(encPath)
	if err != nil {
		return fmt.Errorf("sync: read encrypted file: %w", err)
	}

	plaintext, err := m.encryptor.Decrypt(ciphertext)
	if err != nil {
		return fmt.Errorf("sync: decrypt: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return fmt.Errorf("sync: create output dir: %w", err)
	}

	if err := os.WriteFile(outPath, plaintext, 0o600); err != nil {
		return fmt.Errorf("sync: write decrypted file: %w", err)
	}

	return nil
}
