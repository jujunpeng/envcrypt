package keystore

import (
	"testing"

	"filippo.io/age"
)

func newMergeStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	return s
}

func addMergeRecipient(t *testing.T, s *Store, name string) string {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	pub := id.Recipient().String()
	if err := s.AddRecipient(name, pub); err != nil {
		t.Fatalf("add recipient: %v", err)
	}
	return pub
}

func TestMergeNoOverlap(t *testing.T) {
	dst := newMergeStore(t)
	src := newMergeStore(t)
	addMergeRecipient(t, dst, "alice")
	addMergeRecipient(t, src, "bob")

	res, err := MergeRecipients(dst, src, MergeSkipConflicts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "bob" {
		t.Errorf("expected bob added, got %v", res.Added)
	}
	if len(res.Skipped) != 0 || len(res.Conflict) != 0 {
		t.Errorf("unexpected skipped/conflict: %v %v", res.Skipped, res.Conflict)
	}
}

func TestMergeSkipConflicts(t *testing.T) {
	dst := newMergeStore(t)
	src := newMergeStore(t)
	origKey := addMergeRecipient(t, dst, "alice")
	addMergeRecipient(t, src, "alice")

	res, err := MergeRecipients(dst, src, MergeSkipConflicts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "alice" {
		t.Errorf("expected alice skipped, got %v", res.Skipped)
	}
	r, _ := dst.GetRecipient("alice")
	if r.PublicKey != origKey {
		t.Errorf("key should not have changed on skip")
	}
}

func TestMergeOverwriteConflicts(t *testing.T) {
	dst := newMergeStore(t)
	src := newMergeStore(t)
	addMergeRecipient(t, dst, "alice")
	newKey := addMergeRecipient(t, src, "alice")

	res, err := MergeRecipients(dst, src, MergeOverwriteConflicts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflict) != 1 || res.Conflict[0] != "alice" {
		t.Errorf("expected alice in conflict, got %v", res.Conflict)
	}
	r, _ := dst.GetRecipient("alice")
	if r.PublicKey != newKey {
		t.Errorf("key should have been overwritten")
	}
}

func TestMergeNilStore(t *testing.T) {
	s := newMergeStore(t)
	if _, err := MergeRecipients(nil, s, MergeSkipConflicts); err == nil {
		t.Error("expected error for nil dst")
	}
	if _, err := MergeRecipients(s, nil, MergeSkipConflicts); err == nil {
		t.Error("expected error for nil src")
	}
}
