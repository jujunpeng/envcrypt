package keystore_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"

	"github.com/yourorg/envcrypt/internal/keystore"
)

func generatePublicKey(t *testing.T) string {
	t.Helper()
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generating identity: %v", err)
	}
	return identity.Recipient().String()
}

func TestAddAndGetRecipient(t *testing.T) {
	ks := keystore.New()
	pub := generatePublicKey(t)

	if err := ks.AddRecipient("alice", pub); err != nil {
		t.Fatalf("AddRecipient: %v", err)
	}
	r, err := ks.GetRecipient("alice")
	if err != nil {
		t.Fatalf("GetRecipient: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil recipient")
	}
}

func TestGetRecipientNotFound(t *testing.T) {
	ks := keystore.New()
	_, err := ks.GetRecipient("nobody")
	if err == nil {
		t.Fatal("expected error for missing recipient")
	}
}

func TestInvalidPublicKey(t *testing.T) {
	ks := keystore.New()
	if err := ks.AddRecipient("bad", "not-a-real-key"); err == nil {
		t.Fatal("expected error for invalid public key")
	}
}

func TestAllRecipients(t *testing.T) {
	ks := keystore.New()
	for _, alias := range []string{"alice", "bob", "carol"} {
		if err := ks.AddRecipient(alias, generatePublicKey(t)); err != nil {
			t.Fatalf("AddRecipient %s: %v", alias, err)
		}
	}
	if got := len(ks.AllRecipients()); got != 3 {
		t.Fatalf("expected 3 recipients, got %d", got)
	}
}

func TestLoadFromFile(t *testing.T) {
	pub1 := generatePublicKey(t)
	pub2 := generatePublicKey(t)

	content := "# team keys\nalice=" + pub1 + "\nbob=" + pub2 + "\n"
	tmp := filepath.Join(t.TempDir(), "keys.txt")
	if err := os.WriteFile(tmp, []byte(content), 0600); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}

	ks := keystore.New()
	if err := ks.LoadFromFile(tmp); err != nil {
		t.Fatalf("LoadFromFile: %v", err)
	}
	if len(ks.AllRecipients()) != 2 {
		t.Fatalf("expected 2 recipients after load")
	}
}

func TestLoadFromFileMalformed(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.txt")
	if err := os.WriteFile(tmp, []byte("no-equals-sign\n"), 0600); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	ks := keystore.New()
	if err := ks.LoadFromFile(tmp); err == nil {
		t.Fatal("expected error for malformed key file")
	}
}
