package keystore

import (
	"errors"
	"os"

	"filippo.io/age"
	"filippo.io/age/armor"
)

// EncryptIdentityWithPassphrase encrypts an age identity (private key) using a
// passphrase-based recipient and returns the armored ciphertext.
func EncryptIdentityWithPassphrase(identity *age.X25519Identity, passphrase string) ([]byte, error) {
	if passphrase == "" {
		return nil, errors.New("passphrase must not be empty")
	}

	recipient, err := age.NewScryptRecipient(passphrase)
	if err != nil {
		return nil, err
	}

	var buf []byte
	w := newBytesWriter(&buf)
	armorWriter := armor.NewWriter(w)

	enc, err := age.Encrypt(armorWriter, recipient)
	if err != nil {
		return nil, err
	}

	if _, err := enc.Write([]byte(identity.String())); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	if err := armorWriter.Close(); err != nil {
		return nil, err
	}

	return buf, nil
}

// DecryptIdentityWithPassphrase decrypts an armored age identity using the
// given passphrase and returns the parsed X25519Identity.
func DecryptIdentityWithPassphrase(ciphertext []byte, passphrase string) (*age.X25519Identity, error) {
	if passphrase == "" {
		return nil, errors.New("passphrase must not be empty")
	}

	id, err := age.NewScryptIdentity(passphrase)
	if err != nil {
		return nil, err
	}

	armorReader := armor.NewReader(newBytesReader(ciphertext))
	r, err := age.Decrypt(armorReader, id)
	if err != nil {
		return nil, err
	}

	var raw []byte
	buf := make([]byte, 256)
	for {
		n, readErr := r.Read(buf)
		raw = append(raw, buf[:n]...)
		if readErr != nil {
			break
		}
	}

	identity, err := age.ParseX25519Identity(string(raw))
	if err != nil {
		return nil, err
	}
	return identity, nil
}

// SaveIdentityEncrypted writes an age identity encrypted with passphrase to
// the given file path.
func SaveIdentityEncrypted(path string, identity *age.X25519Identity, passphrase string) error {
	data, err := EncryptIdentityWithPassphrase(identity, passphrase)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LoadIdentityEncrypted reads and decrypts an age identity from a file.
func LoadIdentityEncrypted(path string, passphrase string) (*age.X25519Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecryptIdentityWithPassphrase(data, passphrase)
}

// --- helpers ----------------------------------------------------------------

type bytesWriter struct{ buf *[]byte }

func newBytesWriter(b *[]byte) *bytesWriter { return &bytesWriter{buf: b} }
func (w *bytesWriter) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}

type bytesReader struct {
	data []byte
	pos  int
}

func newBytesReader(data []byte) *bytesReader { return &bytesReader{data: data} }
func (r *bytesReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, errors.New("EOF")
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	if r.pos >= len(r.data) {
		return n, errors.New("EOF")
	}
	return n, nil
}
