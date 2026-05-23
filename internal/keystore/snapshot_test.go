package keystore

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
)

func newSnapshotStore(t *testing.T) *Store {
	t.Helper()
	return &Store{recipients: make(map[string]Recipient)}
}

func addSnapshotRecipient(t *testing.T, s *Store, name string) {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	s.recipients[name] = Recipient{Name: name, PublicKey: id.Recipient().String()}
}

func TestTakeAndLoadSnapshot(t *testing.T) {
	s := newSnapshotStore(t)
	addSnapshotRecipient(t, s, "alice")
	addSnapshotRecipient(t, s, "bob")

	dir := t.TempDir()
	path, err := s.TakeSnapshot(dir, "before-rotation")
	if err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	if snap.Label != "before-rotation" {
		t.Errorf("expected label 'before-rotation', got %q", snap.Label)
	}
	if len(snap.Recipients) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(snap.Recipients))
	}
}

func TestRestoreSnapshot(t *testing.T) {
	s := newSnapshotStore(t)
	addSnapshotRecipient(t, s, "alice")

	dir := t.TempDir()
	path, err := s.TakeSnapshot(dir, "")
	if err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	// Mutate store
	addSnapshotRecipient(t, s, "charlie")

	snap, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}

	s.RestoreSnapshot(snap)

	if _, ok := s.recipients["alice"]; !ok {
		t.Error("expected alice after restore")
	}
	if _, ok := s.recipients["charlie"]; ok {
		t.Error("charlie should not exist after restore")
	}
}

func TestLoadSnapshotMissingFile(t *testing.T) {
	_, err := LoadSnapshot(filepath.Join(t.TempDir(), "no-such-file.json"))
	if err == nil {
		t.Error("expected error for missing snapshot file")
	}
}

func TestSnapshotCreatesDirectory(t *testing.T) {
	s := newSnapshotStore(t)
	addSnapshotRecipient(t, s, "alice")

	dir := filepath.Join(t.TempDir(), "nested", "snapshots")
	path, err := s.TakeSnapshot(dir, "")
	if err != nil {
		t.Fatalf("TakeSnapshot: %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Errorf("snapshot file not found: %v", err)
	}
}
