package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// MigrationVersion represents the current config schema version.
const MigrationVersion = 1

// Migrate checks the config version and applies any necessary migrations.
// It returns an error if migration fails, or nil if the config is up to date.
func Migrate(cfg *Config, cfgPath string) error {
	if cfg == nil {
		return fmt.Errorf("cannot migrate nil config")
	}

	if cfg.Version == MigrationVersion {
		return nil
	}

	if cfg.Version == 0 {
		if err := migrateV0ToV1(cfg); err != nil {
			return fmt.Errorf("migration v0->v1 failed: %w", err)
		}
		cfg.Version = 1
	}

	if err := Save(cfg, cfgPath); err != nil {
		return fmt.Errorf("failed to save migrated config: %w", err)
	}

	return nil
}

// migrateV0ToV1 handles the migration from unversioned (v0) configs to v1.
// In v0, EncryptedFile defaulted to empty string; v1 sets a sensible default.
func migrateV0ToV1(cfg *Config) error {
	if cfg.EncryptedFile == "" {
		cfg.EncryptedFile = ".env.age"
	}
	if cfg.EnvFile == "" {
		cfg.EnvFile = ".env"
	}
	return nil
}

// BackupConfig creates a backup of the config file before migration.
func BackupConfig(cfgPath string) error {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read config for backup: %w", err)
	}

	backupPath := cfgPath + ".bak"
	dir := filepath.Dir(backupPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	if err := os.WriteFile(backupPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	return nil
}
