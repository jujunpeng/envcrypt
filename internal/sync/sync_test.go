package sync_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/envcrypt/internal/crypto"
	"github.com/envcrypt/internal/keystore"
	envsync "github.com/envcrypt/internal/sync"
)

func setupKeypair(t *testing.T) (*age.X25519Identity, *age.X25519Recipient) {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	return id, id.Recipient()
}

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp env: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	id, recipient := setupKeypair(t)

	enc, err := crypto.NewEncryptor(id)
	if err != nil {
		t.Fatalf("new encryptor: %v", err)
	}

	ks := keystore.New()
	if err := ks.AddRecipient("alice", recipient.String()); err != nil {
		t.Fatalf("add recipient: %v", err)
	}

	mgr := envsync.New(enc, ks)

	envContent := "APP_ENV=production\nDB_URL=postgres://localhost/mydb\n"
	envPath := writeTempEnv(t, envContent)

	tmpDir := t.TempDir()
	encPath := filepath.Join(tmpDir, "env.age")
	outPath := filepath.Join(tmpDir, ".env.decrypted")

	if err := mgr.Encrypt(envPath, encPath); err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	if _, err := os.Stat(encPath); err != nil {
		t.Fatalf("encrypted file not created: %v", err)
	}

	if err := mgr.Decrypt(encPath, outPath); err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read decrypted file: %v", err)
	}

	if string(got) != envContent {
		t.Errorf("roundtrip mismatch:\ngot:  %q\nwant: %q", string(got), envContent)
	}
}

func TestEncryptNoRecipients(t *testing.T) {
	id, _ := setupKeypair(t)
	enc, _ := crypto.NewEncryptor(id)
	ks := keystore.New()
	mgr := envsync.New(enc, ks)

	envPath := writeTempEnv(t, "KEY=value\n")
	if err := mgr.Encrypt(envPath, filepath.Join(t.TempDir(), "out.age")); err == nil {
		t.Error("expected error with no recipients, got nil")
	}
}

func TestDecryptMissingFile(t *testing.T) {
	id, _ := setupKeypair(t)
	enc, _ := crypto.NewEncryptor(id)
	mgr := envsync.New(enc, keystore.New())

	if err := mgr.Decrypt("/nonexistent/path.age", "/tmp/out.env"); err == nil {
		t.Error("expected error for missing encrypted file, got nil")
	}
}
