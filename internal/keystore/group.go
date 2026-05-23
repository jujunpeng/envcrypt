package keystore

import (
	"errors"
	"fmt"
	"sort"
)

// Group represents a named collection of recipient names.
type Group struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

// AddGroup creates a new recipient group in the store.
func (s *Store) AddGroup(name string, members []string) error {
	if name == "" {
		return errors.New("group name must not be empty")
	}
	if len(members) == 0 {
		return errors.New("group must have at least one member")
	}
	for _, g := range s.config.Groups {
		if g.Name == name {
			return fmt.Errorf("group %q already exists", name)
		}
	}
	for _, m := range members {
		if _, err := s.GetRecipient(m); err != nil {
			return fmt.Errorf("recipient %q not found", m)
		}
	}
	s.config.Groups = append(s.config.Groups, Group{Name: name, Members: members})
	return s.save()
}

// GetGroup returns the group with the given name.
func (s *Store) GetGroup(name string) (*Group, error) {
	for i, g := range s.config.Groups {
		if g.Name == name {
			return &s.config.Groups[i], nil
		}
	}
	return nil, fmt.Errorf("group %q not found", name)
}

// RemoveGroup deletes a group by name.
func (s *Store) RemoveGroup(name string) error {
	for i, g := range s.config.Groups {
		if g.Name == name {
			s.config.Groups = append(s.config.Groups[:i], s.config.Groups[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("group %q not found", name)
}

// ListGroups returns all groups sorted by name.
func (s *Store) ListGroups() []Group {
	groups := make([]Group, len(s.config.Groups))
	copy(groups, s.config.Groups)
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Name < groups[j].Name
	})
	return groups
}
