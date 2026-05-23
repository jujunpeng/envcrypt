package keystore

import (
	"errors"
	"sort"
	"strings"
)

// TagRecipient adds a tag to a recipient's tag list.
func (s *Store) TagRecipient(name, tag string) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}
	if tag == "" {
		return errors.New("tag must not be empty")
	}
	tag = strings.ToLower(strings.TrimSpace(tag))

	r, err := s.Get(name)
	if err != nil {
		return err
	}

	for _, t := range r.Tags {
		if t == tag {
			return nil // already tagged
		}
	}

	r.Tags = append(r.Tags, tag)
	sort.Strings(r.Tags)
	return s.update(r)
}

// UntagRecipient removes a tag from a recipient's tag list.
func (s *Store) UntagRecipient(name, tag string) error {
	if name == "" {
		return errors.New("recipient name must not be empty")
	}
	tag = strings.ToLower(strings.TrimSpace(tag))

	r, err := s.Get(name)
	if err != nil {
		return err
	}

	updated := r.Tags[:0]
	for _, t := range r.Tags {
		if t != tag {
			updated = append(updated, t)
		}
	}
	r.Tags = updated
	return s.update(r)
}

// RecipientsByTag returns all recipients that have the given tag.
func (s *Store) RecipientsByTag(tag string) ([]Recipient, error) {
	if tag == "" {
		return nil, errors.New("tag must not be empty")
	}
	tag = strings.ToLower(strings.TrimSpace(tag))

	all, err := s.All()
	if err != nil {
		return nil, err
	}

	var matched []Recipient
	for _, r := range all {
		for _, t := range r.Tags {
			if t == tag {
				matched = append(matched, r)
				break
			}
		}
	}
	return matched, nil
}
