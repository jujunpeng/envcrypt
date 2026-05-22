package crypto

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// Encryptor handles age-based encryption and decryption of .env file contents.
type Encryptor struct {
	recipients []age.Recipient
	identities []age.Identity
}

// NewEncryptor creates an Encryptor from PEM-encoded public keys (recipients).
func NewEncryptor(publicKeys []string) (*Encryptor, error) {
	var recipients []age.Recipient
	for _, pub := range publicKeys {
		r, err := age.ParseX25519Recipient(strings.TrimSpace(pub))
		if err != nil {
			return nil, fmt.Errorf("invalid recipient key %q: %w", pub, err)
		}
		recipients = append(recipients, r)
	}
	if len(recipients) == 0 {
		return nil, fmt.Errorf("at least one recipient key is required")
	}
	return &Encryptor{recipients: recipients}, nil
}

// AddIdentity registers a private key identity for decryption.
func (e *Encryptor) AddIdentity(privateKeyPEM string) error {
	ids, err := age.ParseIdentities(strings.NewReader(privateKeyPEM))
	if err != nil {
		return fmt.Errorf("parsing identity: %w", err)
	}
	e.identities = append(e.identities, ids...)
	return nil
}

// Encrypt encrypts plaintext and returns an ASCII-armored ciphertext string.
func (e *Encryptor) Encrypt(plaintext []byte) (string, error) {
	var buf bytes.Buffer
	armorWriter := armor.NewWriter(&buf)
	w, err := age.Encrypt(armorWriter, e.recipients...)
	if err != nil {
		return "", fmt.Errorf("initializing encryption: %w", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		return "", fmt.Errorf("writing plaintext: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("finalizing encryption: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return "", fmt.Errorf("closing armor writer: %w", err)
	}
	return buf.String(), nil
}

// Decrypt decrypts an ASCII-armored ciphertext string using registered identities.
func (e *Encryptor) Decrypt(ciphertext string) ([]byte, error) {
	if len(e.identities) == 0 {
		return nil, fmt.Errorf("no identities registered for decryption")
	}
	armorReader := armor.NewReader(strings.NewReader(ciphertext))
	r, err := age.Decrypt(armorReader, e.identities...)
	if err != nil {
		return nil, fmt.Errorf("decrypting: %w", err)
	}
	plaintext, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("reading decrypted data: %w", err)
	}
	return plaintext, nil
}
