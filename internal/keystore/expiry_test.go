package keystore

import (
	"testing"
	"time"
)

func newExpiryStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func addRecipientForExpiry(t *testing.T, s *Store, name string) {
	t.Helper()
	pub := generatePublicKey(t)
	if err := s.AddRecipient(name, pub); err != nil {
		t.Fatalf("AddRecipient: %v", err)
	}
}

func TestSetAndGetExpiry(t *testing.T) {
	s := newExpiryStore(t)
	addRecipientForExpiry(t, s, "alice")

	expiry := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	if err := s.SetExpiry("alice", expiry); err != nil {
		t.Fatalf("SetExpiry: %v", err)
	}

	got, set, err := s.GetExpiry("alice")
	if err != nil {
		t.Fatalf("GetExpiry: %v", err)
	}
	if !set {
		t.Fatal("expected expiry to be set")
	}
	if !got.Equal(expiry) {
		t.Errorf("expected %v, got %v", expiry, got)
	}
}

func TestClearExpiry(t *testing.T) {
	s := newExpiryStore(t)
	addRecipientForExpiry(t, s, "bob")

	expiry := time.Now().UTC().Add(48 * time.Hour)
	_ = s.SetExpiry("bob", expiry)

	if err := s.ClearExpiry("bob"); err != nil {
		t.Fatalf("ClearExpiry: %v", err)
	}
	_, set, err := s.GetExpiry("bob")
	if err != nil {
		t.Fatalf("GetExpiry after clear: %v", err)
	}
	if set {
		t.Error("expected expiry to be cleared")
	}
}

func TestIsExpiredFuture(t *testing.T) {
	s := newExpiryStore(t)
	addRecipientForExpiry(t, s, "carol")
	_ = s.SetExpiry("carol", time.Now().UTC().Add(time.Hour))

	expired, err := s.IsExpired("carol")
	if err != nil {
		t.Fatalf("IsExpired: %v", err)
	}
	if expired {
		t.Error("expected key not to be expired")
	}
}

func TestIsExpiredPast(t *testing.T) {
	s := newExpiryStore(t)
	addRecipientForExpiry(t, s, "dave")
	_ = s.SetExpiry("dave", time.Now().UTC().Add(-time.Hour))

	expired, err := s.IsExpired("dave")
	if err != nil {
		t.Fatalf("IsExpired: %v", err)
	}
	if !expired {
		t.Error("expected key to be expired")
	}
}

func TestIsExpiredNotSet(t *testing.T) {
	s := newExpiryStore(t)
	addRecipientForExpiry(t, s, "eve")

	expired, err := s.IsExpired("eve")
	if err != nil {
		t.Fatalf("IsExpired: %v", err)
	}
	if expired {
		t.Error("expected no expiry to mean not expired")
	}
}

func TestSetExpiryNotFound(t *testing.T) {
	s := newExpiryStore(t)
	err := s.SetExpiry("ghost", time.Now())
	if err == nil {
		t.Error("expected error for missing recipient")
	}
}
