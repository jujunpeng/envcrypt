package keystore

import (
	"os"
	"path/filepath"
	"testing"
)

func newCommentStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	s, err := New(filepath.Join(dir, "keys.json"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func addRecipientForComment(t *testing.T, s *Store, name string) {
	t.Helper()
	pub := generatePublicKey(t)
	if err := s.AddRecipient(name, pub); err != nil {
		t.Fatalf("AddRecipient: %v", err)
	}
}

func TestSetAndGetComment(t *testing.T) {
	s := newCommentStore(t)
	addRecipientForComment(t, s, "alice")

	if err := s.SetComment("alice", "primary dev key"); err != nil {
		t.Fatalf("SetComment: %v", err)
	}

	comment, err := s.GetComment("alice")
	if err != nil {
		t.Fatalf("GetComment: %v", err)
	}
	if comment != "primary dev key" {
		t.Errorf("expected 'primary dev key', got %q", comment)
	}
}

func TestGetCommentNotSet(t *testing.T) {
	s := newCommentStore(t)
	addRecipientForComment(t, s, "bob")

	comment, err := s.GetComment("bob")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment != "" {
		t.Errorf("expected empty comment, got %q", comment)
	}
}

func TestClearComment(t *testing.T) {
	s := newCommentStore(t)
	addRecipientForComment(t, s, "carol")

	_ = s.SetComment("carol", "to be removed")
	if err := s.ClearComment("carol"); err != nil {
		t.Fatalf("ClearComment: %v", err)
	}

	comment, err := s.GetComment("carol")
	if err != nil {
		t.Fatalf("GetComment after clear: %v", err)
	}
	if comment != "" {
		t.Errorf("expected empty after clear, got %q", comment)
	}
}

func TestSetCommentNotFound(t *testing.T) {
	s := newCommentStore(t)
	if err := s.SetComment("ghost", "hello"); err == nil {
		t.Error("expected error for missing recipient")
	}
}

func TestSetCommentEmptyName(t *testing.T) {
	s := newCommentStore(t)
	if err := s.SetComment("", "some note"); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSetCommentEmptyText(t *testing.T) {
	s := newCommentStore(t)
	addRecipientForComment(t, s, "dave")
	if err := s.SetComment("dave", ""); err == nil {
		t.Error("expected error for empty comment text")
	}
}

func TestCommentPersists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.json")

	s1, _ := New(path)
	addRecipientForComment(t, s1, "eve")
	_ = s1.SetComment("eve", "persistent note")

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	comment, err := s2.GetComment("eve")
	if err != nil {
		t.Fatalf("GetComment after reload: %v", err)
	}
	if comment != "persistent note" {
		t.Errorf("expected 'persistent note', got %q", comment)
	}
	_ = os.RemoveAll(dir)
}
