package keystore

import (
	"testing"
)

func newPinStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func addRecipientForPin(t *testing.T, s *Store, name, key string) {
	t.Helper()
	if err := s.AddRecipient(name, key); err != nil {
		t.Fatalf("AddRecipient(%q): %v", name, err)
	}
}

func TestPinRecipient(t *testing.T) {
	s := newPinStore(t)
	key := generatePublicKey(t)
	addRecipientForPin(t, s, "alice", key)

	if err := s.PinRecipient("alice"); err != nil {
		t.Fatalf("PinRecipient: %v", err)
	}
	r, err := s.GetRecipient("alice")
	if err != nil {
		t.Fatalf("GetRecipient: %v", err)
	}
	if !r.Pinned {
		t.Error("expected recipient to be pinned")
	}
}

func TestPinRecipientAlreadyPinned(t *testing.T) {
	s := newPinStore(t)
	key := generatePublicKey(t)
	addRecipientForPin(t, s, "bob", key)

	_ = s.PinRecipient("bob")
	if err := s.PinRecipient("bob"); err == nil {
		t.Error("expected error when pinning already-pinned recipient")
	}
}

func TestUnpinRecipient(t *testing.T) {
	s := newPinStore(t)
	key := generatePublicKey(t)
	addRecipientForPin(t, s, "carol", key)

	_ = s.PinRecipient("carol")
	if err := s.UnpinRecipient("carol"); err != nil {
		t.Fatalf("UnpinRecipient: %v", err)
	}
	r, err := s.GetRecipient("carol")
	if err != nil {
		t.Fatalf("GetRecipient: %v", err)
	}
	if r.Pinned {
		t.Error("expected recipient to be unpinned")
	}
}

func TestUnpinRecipientNotPinned(t *testing.T) {
	s := newPinStore(t)
	key := generatePublicKey(t)
	addRecipientForPin(t, s, "dave", key)

	if err := s.UnpinRecipient("dave"); err == nil {
		t.Error("expected error when unpinning a non-pinned recipient")
	}
}

func TestPinnedRecipients(t *testing.T) {
	s := newPinStore(t)
	k1, k2 := generatePublicKey(t), generatePublicKey(t)
	addRecipientForPin(t, s, "eve", k1)
	addRecipientForPin(t, s, "frank", k2)

	_ = s.PinRecipient("eve")

	pinned, err := s.PinnedRecipients()
	if err != nil {
		t.Fatalf("PinnedRecipients: %v", err)
	}
	if len(pinned) != 1 {
		t.Fatalf("expected 1 pinned recipient, got %d", len(pinned))
	}
	if pinned[0].Name != "eve" {
		t.Errorf("expected pinned recipient 'eve', got %q", pinned[0].Name)
	}
}

func TestPinEmptyName(t *testing.T) {
	s := newPinStore(t)
	if err := s.PinRecipient(""); err == nil {
		t.Error("expected error for empty name")
	}
	if err := s.UnpinRecipient(""); err == nil {
		t.Error("expected error for empty name")
	}
}
