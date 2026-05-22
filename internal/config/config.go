package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultConfigFile = ".envcrypt.json"

// Config holds the envcrypt project configuration.
type Config struct {
	Recipients []string `json:"recipients"`
	EncryptedFile string   `json:"encrypted_file"`
	EnvFile       string   `json:"env_file"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Recipients:    []string{},
		EncryptedFile: ".env.age",
		EnvFile:       ".env",
	}
}

// Load reads the config file from the given directory.
func Load(dir string) (*Config, error) {
	path := filepath.Join(dir, defaultConfigFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New("config file not found; run 'envcrypt init' first")
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("invalid config file: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to the given directory.
func Save(dir string, cfg *Config) error {
	if cfg == nil {
		return errors.New("config must not be nil")
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(dir, defaultConfigFile)
	return os.WriteFile(path, data, 0644)
}

// AddRecipient appends a public key to the recipient list if not already present.
func (c *Config) AddRecipient(pubKey string) bool {
	for _, r := range c.Recipients {
		if r == pubKey {
			return false
		}
	}
	c.Recipients = append(c.Recipients, pubKey)
	return true
}
