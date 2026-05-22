package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "audit.log")
}

func TestLogAndReadAll(t *testing.T) {
	path := tempLogPath(t)
	l, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := l.Log(EventEncrypt, "encrypted .env", true); err != nil {
		t.Fatalf("Log: %v", err)
	}
	if err := l.Log(EventDecrypt, "decrypted .env", false); err != nil {
		t.Fatalf("Log: %v", err)
	}
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Event != EventEncrypt || !entries[0].Success {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Event != EventDecrypt || entries[1].Success {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestLogAs(t *testing.T) {
	path := tempLogPath(t)
	l, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := l.LogAs("alice", EventAddRecipient, "added bob", true); err != nil {
		t.Fatalf("LogAs: %v", err)
	}
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Actor != "alice" {
		t.Errorf("expected actor alice, got %q", entries[0].Actor)
	}
}

func TestReadAllEmptyFile(t *testing.T) {
	path := tempLogPath(t)
	l, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	entries, err := l.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll on missing file: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestNewCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "subdir", "audit.log")
	_, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

func TestTimestampIsSet(t *testing.T) {
	path := tempLogPath(t)
	l, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := l.Log(EventRotate, "rotated key", true); err != nil {
		t.Fatalf("Log: %v", err)
	}
	entries, _ := l.ReadAll()
	if entries[0].Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}
