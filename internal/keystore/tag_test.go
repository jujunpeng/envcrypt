package keystore

import (
	"os"
	"path/filepath"
	"testing"
)

func newTagStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "keys.json"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func addRecipientForTag(t *testing.T, s *Store, name string) {
	t.Helper()
	pub, _ := generatePublicKey(t)
	if err := s.Add(name, pub); err != nil {
		t.Fatalf("Add(%q): %v", name, err)
	}
}

func TestTagRecipient(t *testing.T) {
	s := newTagStore(t)
	addRecipientForTag(t, s, "alice")

	if err := s.TagRecipient("alice", "backend"); err != nil {
		t.Fatalf("TagRecipient: %v", err)
	}

	r, _ := s.Get("alice")
	if len(r.Tags) != 1 || r.Tags[0] != "backend" {
		t.Errorf("expected tags [backend], got %v", r.Tags)
	}
}

func TestTagRecipientDuplicate(t *testing.T) {
	s := newTagStore(t)
	addRecipientForTag(t, s, "bob")

	_ = s.TagRecipient("bob", "ops")
	if err := s.TagRecipient("bob", "ops"); err != nil {
		t.Fatalf("second TagRecipient should be a no-op, got: %v", err)
	}

	r, _ := s.Get("bob")
	if len(r.Tags) != 1 {
		t.Errorf("expected 1 tag, got %d", len(r.Tags))
	}
}

func TestUntagRecipient(t *testing.T) {
	s := newTagStore(t)
	addRecipientForTag(t, s, "carol")
	_ = s.TagRecipient("carol", "frontend")
	_ = s.TagRecipient("carol", "backend")

	if err := s.UntagRecipient("carol", "frontend"); err != nil {
		t.Fatalf("UntagRecipient: %v", err)
	}

	r, _ := s.Get("carol")
	if len(r.Tags) != 1 || r.Tags[0] != "backend" {
		t.Errorf("expected tags [backend], got %v", r.Tags)
	}
}

func TestRecipientsByTag(t *testing.T) {
	s := newTagStore(t)
	addRecipientForTag(t, s, "dave")
	addRecipientForTag(t, s, "eve")
	addRecipientForTag(t, s, "frank")

	_ = s.TagRecipient("dave", "ops")
	_ = s.TagRecipient("eve", "ops")
	_ = s.TagRecipient("frank", "dev")

	results, err := s.RecipientsByTag("ops")
	if err != nil {
		t.Fatalf("RecipientsByTag: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 recipients with tag ops, got %d", len(results))
	}
}

func TestTagRecipientNotFound(t *testing.T) {
	s := newTagStore(t)
	if err := s.TagRecipient("ghost", "ops"); err == nil {
		t.Error("expected error for unknown recipient")
	}
}

func TestRecipientsByTagEmpty(t *testing.T) {
	s := newTagStore(t)
	addRecipientForTag(t, s, "hank")

	results, err := s.RecipientsByTag("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
