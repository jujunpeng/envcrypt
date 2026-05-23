package keystore

import (
	"errors"
	"fmt"
)

// PinRecipient marks a recipient as pinned, indicating it should always be
// included when encrypting, regardless of group or tag filters.
func (s *Store) PinRecipient(name string) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}
	r, err := s.GetRecipient(name)
	if err != nil {
		return fmt.Errorf("pin: %w", err)
	}
	if r.Pinned {
		return fmt.Errorf("recipient %q is already pinned", name)
	}
	r.Pinned = true
	return s.updateRecipient(r)
}

// UnpinRecipient removes the pinned status from a recipient.
func (s *Store) UnpinRecipient(name string) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}
	r, err := s.GetRecipient(name)
	if err != nil {
		return fmt.Errorf("unpin: %w", err)
	}
	if !r.Pinned {
		return fmt.Errorf("recipient %q is not pinned", name)
	}
	r.Pinned = false
	return s.updateRecipient(r)
}

// PinnedRecipients returns all recipients that are currently pinned.
func (s *Store) PinnedRecipients() ([]Recipient, error) {
	all, err := s.AllRecipients()
	if err != nil {
		return nil, fmt.Errorf("pinned recipients: %w", err)
	}
	var pinned []Recipient
	for _, r := range all {
		if r.Pinned {
			pinned = append(pinned, r)
		}
	}
	return pinned, nil
}
