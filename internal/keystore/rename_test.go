package keystore

import (
	"testing"

	"filippo.io/age"
)

func newRenameStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func addRecipientForRename(t *testing.T, s *Store, name string) string {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	pub := id.Recipient().String()
	if err := s.AddRecipient(name, pub); err != nil {
		t.Fatalf("AddRecipient: %v", err)
	}
	return pub
}

func TestRenameRecipient(t *testing.T) {
	s := newRenameStore(t)
	pub := addRecipientForRename(t, s, "alice")

	if err := s.RenameRecipient("alice", "alice2"); err != nil {
		t.Fatalf("RenameRecipient: %v", err)
	}

	if _, err := s.GetRecipient("alice"); err == nil {
		t.Error("expected old name to be gone")
	}
	r, err := s.GetRecipient("alice2")
	if err != nil {
		t.Fatalf("GetRecipient alice2: %v", err)
	}
	if r.PublicKey != pub {
		t.Errorf("public key mismatch: got %s want %s", r.PublicKey, pub)
	}
}

func TestRenameRecipientNotFound(t *testing.T) {
	s := newRenameStore(t)
	if err := s.RenameRecipient("nobody", "someone"); err == nil {
		t.Error("expected error for missing recipient")
	}
}

func TestRenameRecipientEmptyNames(t *testing.T) {
	s := newRenameStore(t)
	addRecipientForRename(t, s, "alice")

	if err := s.RenameRecipient("", "bob"); err == nil {
		t.Error("expected error for empty old name")
	}
	if err := s.RenameRecipient("alice", ""); err == nil {
		t.Error("expected error for empty new name")
	}
}

func TestRenameRecipientSameName(t *testing.T) {
	s := newRenameStore(t)
	addRecipientForRename(t, s, "alice")

	if err := s.RenameRecipient("alice", "alice"); err == nil {
		t.Error("expected error when old and new names are identical")
	}
}

func TestRenameRecipientConflict(t *testing.T) {
	s := newRenameStore(t)
	addRecipientForRename(t, s, "alice")
	addRecipientForRename(t, s, "bob")

	if err := s.RenameRecipient("alice", "bob"); err == nil {
		t.Error("expected error when new name already exists")
	}
}
