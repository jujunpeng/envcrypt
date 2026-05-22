package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"envcrypt/internal/audit"
	"envcrypt/internal/config"
)

func writeAuditConfig(t *testing.T, dir string) string {
	t.Helper()
	cfg := config.DefaultConfig()
	cfg.AuditLog = filepath.Join(dir, "audit.log")
	cfgPath := filepath.Join(dir, ".envcrypt.json")
	if err := config.Save(cfgPath, cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	return cfgPath
}

func seedAuditLog(t *testing.T, path string) {
	t.Helper()
	entry := audit.Entry{
		Timestamp: time.Now().UTC(),
		Event:     audit.EventEncrypt,
		Actor:     "alice",
		Details:   "encrypted .env",
		Success:   true,
	}
	data, _ := json.Marshal(entry)
	if err := os.WriteFile(path, append(data, '\n'), 0600); err != nil {
		t.Fatalf("seed audit log: %v", err)
	}
}

func TestAuditLogEmpty(t *testing.T) {
	dir := t.TempDir()
	cfgPath := writeAuditConfig(t, dir)
	out, err := withArgs("audit", "--config", cfgPath)
	if err != nil {
		t.Fatalf("audit command: %v", err)
	}
	if out == "" {
		t.Error("expected output message for empty log")
	}
}

func TestAuditLogWithEntries(t *testing.T) {
	dir := t.TempDir()
	cfgPath := writeAuditConfig(t, dir)
	cfg, _ := config.Load(cfgPath)
	seedAuditLog(t, cfg.AuditLog)
	out, err := withArgs("audit", "--config", cfgPath)
	if err != nil {
		t.Fatalf("audit command: %v", err)
	}
	if out == "" {
		t.Error("expected table output")
	}
	if !containsStr(out, "alice") {
		t.Errorf("expected actor 'alice' in output, got: %s", out)
	}
}

func TestAuditLogMissingConfig(t *testing.T) {
	_, err := withArgs("audit", "--config", "/nonexistent/.envcrypt.json")
	if err == nil {
		t.Error("expected error for missing config")
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
