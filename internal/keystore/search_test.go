package keystore

import (
	"testing"

	"filippo.io/age"
)

func newSearchStore(t *testing.T) *Store {
	t.Helper()
	s := New()
	keys := map[string]string{}
	names := []string{"alice", "bob", "charlie", "dave"}
	for _, n := range names {
		id, err := age.GenerateX25519Identity()
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		keys[n] = id.Recipient().String()
	}
	for n, k := range keys {
		if err := s.AddRecipient(n, k); err != nil {
			t.Fatalf("add recipient: %v", err)
		}
	}
	return s
}

func TestSearchByName(t *testing.T) {
	s := newSearchStore(t)
	results := s.Search("ali")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "alice" {
		t.Errorf("expected alice, got %s", results[0].Name)
	}
	if results[0].MatchedBy != "name" {
		t.Errorf("expected match by name, got %s", results[0].MatchedBy)
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	s := newSearchStore(t)
	results := s.Search("")
	if results != nil {
		t.Errorf("expected nil results for empty query")
	}
}

func TestSearchNoMatch(t *testing.T) {
	s := newSearchStore(t)
	results := s.Search("zzz_no_match")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearchCaseInsensitive(t *testing.T) {
	s := newSearchStore(t)
	results := s.Search("ALICE")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestFindByPrefix(t *testing.T) {
	s := newSearchStore(t)
	name, _, found := s.FindByPrefix("bob")
	if !found {
		t.Fatal("expected to find bob")
	}
	if name != "bob" {
		t.Errorf("expected bob, got %s", name)
	}
}

func TestFindByPrefixNotFound(t *testing.T) {
	s := newSearchStore(t)
	_, _, found := s.FindByPrefix("xyz")
	if found {
		t.Error("expected not found")
	}
}

func TestFindByPrefixEmpty(t *testing.T) {
	s := newSearchStore(t)
	_, _, found := s.FindByPrefix("")
	if found {
		t.Error("expected not found for empty prefix")
	}
}
