package keystore

import (
	"os"
	"testing"
)

func newGroupStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	path := dir + "/keys.json"
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func addRecipientForGroup(t *testing.T, s *Store, name, pub string) {
	t.Helper()
	if err := s.AddRecipient(name, pub); err != nil {
		t.Fatalf("AddRecipient(%q): %v", name, err)
	}
}

func TestAddAndGetGroup(t *testing.T) {
	s := newGroupStore(t)
	pub1, _ := generatePublicKey(t)
	pub2, _ := generatePublicKey(t)
	addRecipientForGroup(t, s, "alice", pub1)
	addRecipientForGroup(t, s, "bob", pub2)

	if err := s.AddGroup("team", []string{"alice", "bob"}); err != nil {
		t.Fatalf("AddGroup: %v", err)
	}
	g, err := s.GetGroup("team")
	if err != nil {
		t.Fatalf("GetGroup: %v", err)
	}
	if len(g.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(g.Members))
	}
}

func TestAddGroupDuplicate(t *testing.T) {
	s := newGroupStore(t)
	pub, _ := generatePublicKey(t)
	addRecipientForGroup(t, s, "alice", pub)
	_ = s.AddGroup("team", []string{"alice"})
	if err := s.AddGroup("team", []string{"alice"}); err == nil {
		t.Error("expected error for duplicate group")
	}
}

func TestAddGroupMissingRecipient(t *testing.T) {
	s := newGroupStore(t)
	if err := s.AddGroup("team", []string{"ghost"}); err == nil {
		t.Error("expected error for unknown recipient")
	}
}

func TestRemoveGroup(t *testing.T) {
	s := newGroupStore(t)
	pub, _ := generatePublicKey(t)
	addRecipientForGroup(t, s, "alice", pub)
	_ = s.AddGroup("team", []string{"alice"})
	if err := s.RemoveGroup("team"); err != nil {
		t.Fatalf("RemoveGroup: %v", err)
	}
	if _, err := s.GetGroup("team"); err == nil {
		t.Error("expected error after removal")
	}
}

func TestListGroupsSorted(t *testing.T) {
	s := newGroupStore(t)
	pub, _ := generatePublicKey(t)
	addRecipientForGroup(t, s, "alice", pub)
	_ = s.AddGroup("zeta", []string{"alice"})
	_ = s.AddGroup("alpha", []string{"alice"})
	groups := s.ListGroups()
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	if groups[0].Name != "alpha" || groups[1].Name != "zeta" {
		t.Errorf("unexpected order: %v", groups)
	}
}

func TestAddGroupEmptyName(t *testing.T) {
	s := newGroupStore(t)
	if err := s.AddGroup("", []string{"x"}); err == nil {
		t.Error("expected error for empty group name")
	}
}

func init() {
	_ = os.Setenv
}
