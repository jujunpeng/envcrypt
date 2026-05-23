package keystore

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
)

func generateVerifyKey(t *testing.T) string {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	return id.Recipient().String()
}

func newVerifyStore(t *testing.T) *KeyStore {
	t.Helper()
	dir := t.TempDir()
	ks, err := New(filepath.Join(dir, "keys.json"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return ks
}

func TestVerifyRecipientValid(t *testing.T) {
	ks := newVerifyStore(t)
	pub := generateVerifyKey(t)
	if err := ks.AddRecipient("alice", pub); err != nil {
		t.Fatalf("AddRecipient: %v", err)
	}

	res := ks.VerifyRecipient("alice")
	if !res.Valid {
		t.Errorf("expected valid, got error: %v", res.Error)
	}
	if res.PublicKey != pub {
		t.Errorf("public key mismatch")
	}
}

func TestVerifyRecipientNotFound(t *testing.T) {
	ks := newVerifyStore(t)
	res := ks.VerifyRecipient("nobody")
	if res.Valid {
		t.Error("expected invalid for missing recipient")
	}
	if res.Error == nil {
		t.Error("expected non-nil error")
	}
}

func TestVerifyRecipientBadKey(t *testing.T) {
	ks := newVerifyStore(t)
	// Manually inject a malformed key by writing directly to the store file.
	path := filepath.Join(t.TempDir(), "keys.json")
	content := `{"version":1,"recipients":[{"name":"bob","public_key":"notavalidkey"}]}`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	ks2, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	res := ks2.VerifyRecipient("bob")
	if res.Valid {
		t.Error("expected invalid for bad key")
	}
}

func TestVerifyAll(t *testing.T) {
	ks := newVerifyStore(t)
	for _, name := range []string{"alice", "bob", "carol"} {
		if err := ks.AddRecipient(name, generateVerifyKey(t)); err != nil {
			t.Fatalf("AddRecipient %s: %v", name, err)
		}
	}

	results := ks.VerifyAll()
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	for _, r := range results {
		if !r.Valid {
			t.Errorf("recipient %s should be valid: %v", r.Name, r.Error)
		}
	}
}

func TestVerifyAllEmpty(t *testing.T) {
	ks := newVerifyStore(t)
	results := ks.VerifyAll()
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty store, got %d", len(results))
	}
}
