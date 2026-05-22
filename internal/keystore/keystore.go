package keystore

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"filippo.io/age"
)

const defaultKeyFile = ".envcrypt/keys.txt"

// KeyStore manages age recipient public keys for team members.
type KeyStore struct {
	keys map[string]*age.X25519Recipient
}

// New creates an empty KeyStore.
func New() *KeyStore {
	return &KeyStore{keys: make(map[string]*age.X25519Recipient)}
}

// AddRecipient parses and stores a public key under the given alias.
func (ks *KeyStore) AddRecipient(alias, publicKey string) error {
	recipient, err := age.ParseX25519Recipient(strings.TrimSpace(publicKey))
	if err != nil {
		return fmt.Errorf("invalid public key for %q: %w", alias, err)
	}
	ks.keys[alias] = recipient
	return nil
}

// GetRecipient returns the recipient for the given alias.
func (ks *KeyStore) GetRecipient(alias string) (*age.X25519Recipient, error) {
	r, ok := ks.keys[alias]
	if !ok {
		return nil, fmt.Errorf("recipient %q not found", alias)
	}
	return r, nil
}

// AllRecipients returns all stored recipients.
func (ks *KeyStore) AllRecipients() []*age.X25519Recipient {
	result := make([]*age.X25519Recipient, 0, len(ks.keys))
	for _, r := range ks.keys {
		result = append(result, r)
	}
	return result
}

// LoadFromFile reads alias=publickey pairs from a file.
func (ks *KeyStore) LoadFromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading key file: %w", err)
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("line %d: expected alias=publickey format", i+1)
		}
		if err := ks.AddRecipient(strings.TrimSpace(parts[0]), parts[1]); err != nil {
			return err
		}
	}
	return nil
}

// DefaultKeyFilePath returns the default key file path relative to home dir.
func DefaultKeyFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("could not determine home directory")
	}
	return filepath.Join(home, defaultKeyFile), nil
}
