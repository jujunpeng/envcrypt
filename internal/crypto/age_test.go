package crypto_test

import (
	"strings"
	"testing"

	"filippo.io/age"

	"github.com/your-org/envcrypt/internal/crypto"
)

func generateTestKeypair(t *testing.T) (string, string) {
	t.Helper()
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generating keypair: %v", err)
	}
	return identity.Recipient().String(), identity.String()
}

func TestEncryptDecryptRoundtrip(t *testing.T) {
	pubKey, privKey := generateTestKeypair(t)

	enc, err := crypto.NewEncryptor([]string{pubKey})
	if err != nil {
		t.Fatalf("NewEncryptor: %v", err)
	}
	if err := enc.AddIdentity(privKey); err != nil {
		t.Fatalf("AddIdentity: %v", err)
	}

	plaintext := []byte("DB_PASSWORD=supersecret\nAPI_KEY=abc123\n")
	ciphertext, err := enc.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if !strings.Contains(ciphertext, "-----BEGIN AGE ENCRYPTED FILE-----") {
		t.Error("ciphertext missing age armor header")
	}

	decrypted, err := enc.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if string(decrypted) != string(plaintext) {
		t.Errorf("roundtrip mismatch: got %q, want %q", decrypted, plaintext)
	}
}

func TestNewEncryptorNoKeys(t *testing.T) {
	_, err := crypto.NewEncryptor([]string{})
	if err == nil {
		t.Error("expected error with no recipient keys, got nil")
	}
}

func TestDecryptWithoutIdentity(t *testing.T) {
	pubKey, _ := generateTestKeypair(t)
	enc, _ := crypto.NewEncryptor([]string{pubKey})

	ciphertext, err := enc.Encrypt([]byte("SECRET=value"))
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	_, err = enc.Decrypt(ciphertext)
	if err == nil {
		t.Error("expected error decrypting without identity, got nil")
	}
}

func TestInvalidRecipientKey(t *testing.T) {
	_, err := crypto.NewEncryptor([]string{"not-a-valid-age-key"})
	if err == nil {
		t.Error("expected error for invalid recipient key, got nil")
	}
}
