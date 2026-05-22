package config

import (
	"errors"
	"fmt"
	"strings"
)

// Validate checks that the Config fields are valid.
func (c *Config) Validate() error {
	if c.EnvFile == "" {
		return errors.New("env_file must not be empty")
	}
	if c.EncryptedFile == "" {
		return errors.New("encrypted_file must not be empty")
	}
	if c.EnvFile == c.EncryptedFile {
		return errors.New("env_file and encrypted_file must be different paths")
	}
	for i, r := range c.Recipients {
		if err := validateRecipient(r); err != nil {
			return fmt.Errorf("recipient[%d] %q: %w", i, r, err)
		}
	}
	return nil
}

// validateRecipient performs basic sanity checks on an age public key.
func validateRecipient(key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return errors.New("recipient key must not be empty")
	}
	if !strings.HasPrefix(key, "age1") {
		return errors.New("recipient key must start with 'age1'")
	}
	if len(key) < 10 {
		return errors.New("recipient key is too short")
	}
	return nil
}
