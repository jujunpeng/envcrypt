package keystore

import (
	"errors"
	"fmt"
	"strings"

	"filippo.io/age"
)

// VerifyResult holds the outcome of verifying a single recipient.
type VerifyResult struct {
	Name    string
	PublicKey string
	Valid   bool
	Error   error
}

// VerifyRecipient checks whether a stored recipient's public key is a valid age
// X25519 public key and returns a VerifyResult.
func (ks *KeyStore) VerifyRecipient(name string) VerifyResult {
	recipient, err := ks.GetRecipient(name)
	if err != nil {
		return VerifyResult{Name: name, Valid: false, Error: err}
	}

	if err := validateAgePublicKey(recipient.PublicKey); err != nil {
		return VerifyResult{
			Name:      name,
			PublicKey: recipient.PublicKey,
			Valid:     false,
			Error:     err,
		}
	}

	return VerifyResult{
		Name:      name,
		PublicKey: recipient.PublicKey,
		Valid:     true,
	}
}

// VerifyAll verifies every recipient in the keystore and returns a slice of
// VerifyResult. It never returns an error itself; individual failures are
// captured inside each VerifyResult.
func (ks *KeyStore) VerifyAll() []VerifyResult {
	recipients := ks.AllRecipients()
	results := make([]VerifyResult, 0, len(recipients))
	for _, r := range recipients {
		results = append(results, ks.VerifyRecipient(r.Name))
	}
	return results
}

// validateAgePublicKey parses the given string as an age X25519 recipient to
// confirm it is well-formed.
func validateAgePublicKey(pub string) error {
	if !strings.HasPrefix(pub, "age1") {
		return errors.New("public key must start with \"age1\"")
	}
	if _, err := age.ParseX25519Recipient(pub); err != nil {
		return fmt.Errorf("invalid age public key: %w", err)
	}
	return nil
}
