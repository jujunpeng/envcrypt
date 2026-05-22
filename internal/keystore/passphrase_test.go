package keystore

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
)

func generateIdentity(t *testing.T) *age.X25519Identity {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}
	return id
}

func TestEncryptDecryptIdentityRoundtrip(t *testing.T) {
	id := generateIdentity(t)
	passphrase := "super-secret-passphrase"

	ciphertext, err := EncryptIdentityWithPassphrase(id, passphrase)
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	recovered, err := DecryptIdentityWithPassphrase(ciphertext, passphrase)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}

	if recovered.String() != id.String() {
		t.Errorf("identity mismatch: got %q, want %q", recovered.String(), id.String())
	}
}

func TestEncryptIdentityEmptyPassphrase(t *testing.T) {
	id := generateIdentity(t)
	_, err := EncryptIdentityWithPassphrase(id, "")
	if err == nil {
		t.Fatal("expected error for empty passphrase, got nil")
	}
}

func TestDecryptIdentityWrongPassphrase(t *testing.T) {
	id := generateIdentity(t)

	ciphertext, err := EncryptIdentityWithPassphrase(id, "correct-passphrase")
	if err != nil {
		t.Fatalf("encrypt: %v", err)
	}

	_, err = DecryptIdentityWithPassphrase(ciphertext, "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error for wrong passphrase, got nil")
	}
}

func TestDecryptIdentityEmptyPassphrase(t *testing.T) {
	_, err := DecryptIdentityWithPassphrase([]byte("some data"), "")
	if err == nil {
		t.Fatal("expected error for empty passphrase, got nil")
	}
}

func TestSaveAndLoadIdentityEncrypted(t *testing.T) {
	id := generateIdentity(t)
	passphrase := "file-passphrase-123"

	dir := t.TempDir()
	path := filepath.Join(dir, "identity.age")

	if err := SaveIdentityEncrypted(path, id, passphrase); err != nil {
		t.Fatalf("save: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}

	recovered, err := LoadIdentityEncrypted(path, passphrase)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if recovered.String() != id.String() {
		t.Errorf("identity mismatch after save/load")
	}
}

func TestLoadIdentityEncryptedMissingFile(t *testing.T) {
	_, err := LoadIdentityEncrypted("/nonexistent/path/identity.age", "passphrase")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
