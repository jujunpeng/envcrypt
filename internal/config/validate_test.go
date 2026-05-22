package config_test

import (
	"testing"

	"github.com/yourorg/envcrypt/internal/config"
)

func TestValidateValid(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Recipients = []string{"age1qyqszqgpqyqszqgpqyqszqgp"}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got: %v", err)
	}
}

func TestValidateEmptyEnvFile(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.EnvFile = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty env_file")
	}
}

func TestValidateEmptyEncryptedFile(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.EncryptedFile = ""
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty encrypted_file")
	}
}

func TestValidateSameFilePaths(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.EnvFile = ".env"
	cfg.EncryptedFile = ".env"
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when env_file equals encrypted_file")
	}
}

func TestValidateInvalidRecipientPrefix(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Recipients = []string{"ssh-rsa AAAAB3NzaC1"}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for non-age recipient key")
	}
}

func TestValidateEmptyRecipientKey(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Recipients = []string{""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty recipient key")
	}
}
