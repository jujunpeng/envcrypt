package keystore

import "fmt"

// RenameRecipient changes the display name of a recipient identified by oldName.
// Returns an error if oldName does not exist or if newName is already taken.
func (s *Store) RenameRecipient(oldName, newName string) error {
	if oldName == "" {
		return fmt.Errorf("old name must not be empty")
	}
	if newName == "" {
		return fmt.Errorf("new name must not be empty")
	}
	if oldName == newName {
		return fmt.Errorf("old and new names are identical")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	recipient, ok := s.recipients[oldName]
	if !ok {
		return fmt.Errorf("recipient %q not found", oldName)
	}
	if _, exists := s.recipients[newName]; exists {
		return fmt.Errorf("recipient %q already exists", newName)
	}

	recipient.Name = newName
	s.recipients[newName] = recipient
	delete(s.recipients, oldName)

	// Update tags index if present
	for tag, names := range s.tags {
		for i, n := range names {
			if n == oldName {
				s.tags[tag][i] = newName
				break
			}
		}
	}

	// Update groups index if present
	for group, members := range s.groups {
		for i, n := range members {
			if n == oldName {
				s.groups[group][i] = newName
				break
			}
		}
	}

	return nil
}
