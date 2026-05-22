package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateAlreadyCurrent(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Version = MigrationVersion

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "envcrypt.json")

	if err := Migrate(cfg, cfgPath); err != nil {
		t.Fatalf("expected no error for current version, got: %v", err)
	}
}

func TestMigrateV0ToV1(t *testing.T) {
	cfg := &Config{
		Version:       0,
		EnvFile:       "",
		EncryptedFile: "",
		Recipients:    []string{},
	}

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "envcrypt.json")

	if err := Migrate(cfg, cfgPath); err != nil {
		t.Fatalf("unexpected error during migration: %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("expected version 1 after migration, got %d", cfg.Version)
	}
	if cfg.EnvFile != ".env" {
		t.Errorf("expected EnvFile '.env', got %q", cfg.EnvFile)
	}
	if cfg.EncryptedFile != ".env.age" {
		t.Errorf("expected EncryptedFile '.env.age', got %q", cfg.EncryptedFile)
	}

	// Verify it was persisted
	loaded, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("failed to load migrated config: %v", err)
	}
	if loaded.Version != 1 {
		t.Errorf("persisted version mismatch: got %d", loaded.Version)
	}
}

func TestMigrateNilConfig(t *testing.T) {
	err := Migrate(nil, "/tmp/irrelevant.json")
	if err == nil {
		t.Error("expected error for nil config, got nil")
	}
}

func TestBackupConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "envcrypt.json")

	original := []byte(`{"version":0}`)
	if err := os.WriteFile(cfgPath, original, 0600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	if err := BackupConfig(cfgPath); err != nil {
		t.Fatalf("unexpected error during backup: %v", err)
	}

	backupPath := cfgPath + ".bak"
	data, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("backup file not found: %v", err)
	}
	if string(data) != string(original) {
		t.Errorf("backup content mismatch: got %q", string(data))
	}
}

func TestBackupConfigMissingFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "nonexistent.json")

	if err := BackupConfig(cfgPath); err != nil {
		t.Errorf("expected no error for missing file, got: %v", err)
	}
}
