package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	CurrentVersion  = 1
	DefaultFileName = ".envcrypt.json"
)

// Config holds the envcrypt project configuration.
type Config struct {
	Version       int      `json:"version"`
	EnvFile       string   `json:"env_file"`
	EncryptedFile string   `json:"encrypted_file"`
	Recipients    []string `json:"recipients"`
	AuditLog      string   `json:"audit_log,omitempty"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Version:       CurrentVersion,
		EnvFile:       ".env",
		EncryptedFile: ".env.age",
		Recipients:    []string{},
		AuditLog:      ".envcrypt-audit.log",
	}
}

// Load reads a Config from the given file path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: read %q: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parse %q: %w", path, err)
	}
	return &cfg, nil
}

// Save writes the Config to the given file path, creating directories as needed.
func Save(path string, cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config: cannot save nil config")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("config: create directory: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("config: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("config: write %q: %w", path, err)
	}
	return nil
}

// AddRecipient appends a public key to the Recipients list if not already present.
func (c *Config) AddRecipient(pubkey string) {
	for _, r := range c.Recipients {
		if r == pubkey {
			return
		}
	}
	c.Recipients = append(c.Recipients, pubkey)
}
