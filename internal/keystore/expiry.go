package keystore

import (
	"errors"
	"time"
)

const expiryLayout = time.RFC3339

// SetExpiry sets an expiration date on a recipient key.
func (s *Store) SetExpiry(name string, expiry time.Time) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}
	r, ok := s.recipients[name]
	if !ok {
		return errors.New("recipient not found: " + name)
	}
	if r.Meta == nil {
		r.Meta = make(map[string]string)
	}
	r.Meta["expiry"] = expiry.UTC().Format(expiryLayout)
	s.recipients[name] = r
	return s.persist()
}

// ClearExpiry removes the expiration date from a recipient key.
func (s *Store) ClearExpiry(name string) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}
	r, ok := s.recipients[name]
	if !ok {
		return errors.New("recipient not found: " + name)
	}
	if r.Meta != nil {
		delete(r.Meta, "expiry")
	}
	s.recipients[name] = r
	return s.persist()
}

// GetExpiry returns the expiration time for a recipient, and whether it is set.
func (s *Store) GetExpiry(name string) (time.Time, bool, error) {
	r, ok := s.recipients[name]
	if !ok {
		return time.Time{}, false, errors.New("recipient not found: " + name)
	}
	if r.Meta == nil {
		return time.Time{}, false, nil
	}
	val, exists := r.Meta["expiry"]
	if !exists || val == "" {
		return time.Time{}, false, nil
	}
	t, err := time.Parse(expiryLayout, val)
	if err != nil {
		return time.Time{}, false, errors.New("invalid expiry format for " + name)
	}
	return t, true, nil
}

// IsExpired reports whether a recipient's key has passed its expiry date.
func (s *Store) IsExpired(name string) (bool, error) {
	t, set, err := s.GetExpiry(name)
	if err != nil {
		return false, err
	}
	if !set {
		return false, nil
	}
	return time.Now().UTC().After(t), nil
}
