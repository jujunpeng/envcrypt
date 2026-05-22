package keystore

import (
	"fmt"
	"time"
)

// RotationRecord tracks when a key was rotated and by whom.
type RotationRecord struct {
	OldKey    string    `json:"old_key"`
	NewKey    string    `json:"new_key"`
	RotatedAt time.Time `json:"rotated_at"`
	RotatedBy string    `json:"rotated_by"`
}

// RotateRecipient replaces an existing recipient's public key with a new one.
// It returns a RotationRecord describing the change.
func (ks *KeyStore) RotateRecipient(name, newPublicKey string) (*RotationRecord, error) {
	if name == "" {
		return nil, fmt.Errorf("recipient name must not be empty")
	}

	old, err := ks.GetRecipient(name)
	if err != nil {
		return nil, fmt.Errorf("recipient %q not found: %w", name, err)
	}

	if err := ks.AddRecipient(name, newPublicKey); err != nil {
		return nil, fmt.Errorf("failed to set new key for %q: %w", name, err)
	}

	record := &RotationRecord{
		OldKey:    old,
		NewKey:    newPublicKey,
		RotatedAt: time.Now().UTC(),
		RotatedBy: name,
	}
	return record, nil
}

// RemoveRecipient removes a recipient from the keystore by name.
func (ks *KeyStore) RemoveRecipient(name string) error {
	if name == "" {
		return fmt.Errorf("recipient name must not be empty")
	}
	if _, err := ks.GetRecipient(name); err != nil {
		return fmt.Errorf("recipient %q not found: %w", name, err)
	}
	delete(ks.recipients, name)
	return nil
}
