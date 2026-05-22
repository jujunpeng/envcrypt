package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envcrypt/internal/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.EnvFile != ".env" {
		t.Errorf("expected .env, got %s", cfg.EnvFile)
	}
	if cfg.EncryptedFile != ".env.age" {
		t.Errorf("expected .env.age, got %s", cfg.EncryptedFile)
	}
	if len(cfg.Recipients) != 0 {
		t.Errorf("expected empty recipients")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	cfg := config.DefaultConfig()
	cfg.Recipients = []string{"age1abc123"}

	if err := config.Save(dir, cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := config.Load(dir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Recipients) != 1 || loaded.Recipients[0] != "age1abc123" {
		t.Errorf("unexpected recipients: %v", loaded.Recipients)
	}
}

func TestLoadNotFound(t *testing.T) {
	dir := t.TempDir()
	_, err := config.Load(dir)
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestSaveNilConfig(t *testing.T) {
	dir := t.TempDir()
	if err := config.Save(dir, nil); err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestAddRecipient(t *testing.T) {
	cfg := config.DefaultConfig()
	added := cfg.AddRecipient("age1xyz")
	if !added {
		t.Fatal("expected recipient to be added")
	}
	if len(cfg.Recipients) != 1 {
		t.Fatalf("expected 1 recipient, got %d", len(cfg.Recipients))
	}

	// Adding duplicate should return false
	added = cfg.AddRecipient("age1xyz")
	if added {
		t.Fatal("expected duplicate to not be added")
	}
	if len(cfg.Recipients) != 1 {
		t.Fatalf("expected still 1 recipient, got %d", len(cfg.Recipients))
	}
}

func TestConfigFileLocation(t *testing.T) {
	dir := t.TempDir()
	cfg := config.DefaultConfig()
	_ = config.Save(dir, cfg)

	_, err := os.Stat(filepath.Join(dir, ".envcrypt.json"))
	if err != nil {
		t.Errorf("expected .envcrypt.json to exist: %v", err)
	}
}
