package keystore

import (
	"errors"
	"time"
)

// RecipientComment holds a human-readable note attached to a recipient.
type RecipientComment struct {
	Text      string    `json:"text"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SetComment sets or updates the comment for the named recipient.
func (s *Store) SetComment(name, text string) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}
	if text == "" {
		return errors.New("comment text must not be empty")
	}

	recipient, ok := s.recipients[name]
	if !ok {
		return errors.New("recipient not found: " + name)
	}

	if recipient.Metadata == nil {
		recipient.Metadata = make(map[string]string)
	}
	recipient.Metadata["comment"] = text
	recipient.Metadata["comment_updated_at"] = time.Now().UTC().Format(time.RFC3339)
	s.recipients[name] = recipient
	return s.save()
}

// GetComment returns the comment for the named recipient, or an empty string if none is set.
func (s *Store) GetComment(name string) (string, error) {
	if name == "" {
		return "", errors.New("recipient name must not be empty")
	}

	recipient, ok := s.recipients[name]
	if !ok {
		return "", errors.New("recipient not found: " + name)
	}

	if recipient.Metadata == nil {
		return "", nil
	}
	return recipient.Metadata["comment"], nil
}

// ClearComment removes the comment from the named recipient.
func (s *Store) ClearComment(name string) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}

	recipient, ok := s.recipients[name]
	if !ok {
		return errors.New("recipient not found: " + name)
	}

	if recipient.Metadata != nil {
		delete(recipient.Metadata, "comment")
		delete(recipient.Metadata, "comment_updated_at")
		s.recipients[name] = recipient
	}
	return s.save()
}
