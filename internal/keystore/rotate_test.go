package keystore

import (
	"testing"

	"filippo.io/age"
)

func newKeyForRotate(t *testing.T) string {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return id.Recipient().String()
}

func TestRotateRecipient(t *testing.T) {
	ks := New()
	origKey := newKeyForRotate(t)
	newKey := newKeyForRotate(t)

	if err := ks.AddRecipient("alice", origKey); err != nil {
		t.Fatalf("AddRecipient: %v", err)
	}

	rec, err := ks.RotateRecipient("alice", newKey)
	if err != nil {
		t.Fatalf("RotateRecipient: %v", err)
	}

	if rec.OldKey != origKey {
		t.Errorf("expected old key %q, got %q", origKey, rec.OldKey)
	}
	if rec.NewKey != newKey {
		t.Errorf("expected new key %q, got %q", newKey, rec.NewKey)
	}
	if rec.RotatedBy != "alice" {
		t.Errorf("expected RotatedBy alice, got %q", rec.RotatedBy)
	}

	stored, err := ks.GetRecipient("alice")
	if err != nil {
		t.Fatalf("GetRecipient after rotate: %v", err)
	}
	if stored != newKey {
		t.Errorf("stored key mismatch after rotation")
	}
}

func TestRotateRecipientNotFound(t *testing.T) {
	ks := New()
	_, err := ks.RotateRecipient("ghost", newKeyForRotate(t))
	if err == nil {
		t.Fatal("expected error rotating unknown recipient")
	}
}

func TestRotateRecipientEmptyName(t *testing.T) {
	ks := New()
	_, err := ks.RotateRecipient("", newKeyForRotate(t))
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRemoveRecipient(t *testing.T) {
	ks := New()
	key := newKeyForRotate(t)
	if err := ks.AddRecipient("bob", key); err != nil {
		t.Fatalf("AddRecipient: %v", err)
	}
	if err := ks.RemoveRecipient("bob"); err != nil {
		t.Fatalf("RemoveRecipient: %v", err)
	}
	if _, err := ks.GetRecipient("bob"); err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestRemoveRecipientNotFound(t *testing.T) {
	ks := New()
	if err := ks.RemoveRecipient("nobody"); err == nil {
		t.Fatal("expected error removing non-existent recipient")
	}
}
