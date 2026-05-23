package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"filippo.io/age"

	"envcrypt/internal/config"
)

func writeSearchConfig(t *testing.T, names []string) string {
	t.Helper()
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "envcrypt.json")

	cfg := config.DefaultConfig()
	for _, n := range names {
		id, err := age.GenerateX25519Identity()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		cfg.Recipients = append(cfg.Recipients, config.Recipient{
			Name:      n,
			PublicKey: id.Recipient().String(),
		})
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(cfgPath, data, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("ENVCRYPT_CONFIG", cfgPath)
	return cfgPath
}

func TestSearchFound(t *testing.T) {
	writeSearchConfig(t, []string{"alice", "bob", "charlie"})
	var buf bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"search", "alice"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "alice") {
		t.Errorf("expected alice in output, got: %s", buf.String())
	}
}

func TestSearchNotFound(t *testing.T) {
	writeSearchConfig(t, []string{"alice", "bob"})
	var buf bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"search", "zzz_nobody"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No recipients matched") {
		t.Errorf("expected no-match message, got: %s", buf.String())
	}
}

func TestSearchMissingArg(t *testing.T) {
	writeSearchConfig(t, []string{"alice"})
	cmd := newRootCmd()
	cmd.SetArgs([]string{"search"})
	if err := cmd.Execute(); err == nil {
		t.Error("expected error for missing argument")
	}
}
