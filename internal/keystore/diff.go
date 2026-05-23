package keystore

import "sort"

// DiffResult holds the result of comparing two sets of recipients.
type DiffResult struct {
	Added   []Recipient
	Removed []Recipient
	Unchanged []Recipient
}

// Recipient represents a named age public key.
type Recipient struct {
	Name      string
	PublicKey string
}

// DiffRecipients compares two slices of recipients and returns what was added,
// removed, or remained unchanged between the old and new sets.
func DiffRecipients(oldRecipients, newRecipients []Recipient) DiffResult {
	oldMap := make(map[string]Recipient, len(oldRecipients))
	for _, r := range oldRecipients {
		oldMap[r.Name] = r
	}

	newMap := make(map[string]Recipient, len(newRecipients))
	for _, r := range newRecipients {
		newMap[r.Name] = r
	}

	var result DiffResult

	for _, r := range newRecipients {
		if old, exists := oldMap[r.Name]; exists {
			if old.PublicKey == r.PublicKey {
				result.Unchanged = append(result.Unchanged, r)
			} else {
				// Key changed — treat old as removed, new as added
				result.Removed = append(result.Removed, old)
				result.Added = append(result.Added, r)
			}
		} else {
			result.Added = append(result.Added, r)
		}
	}

	for _, r := range oldRecipients {
		if _, exists := newMap[r.Name]; !exists {
			result.Removed = append(result.Removed, r)
		}
	}

	sort.Slice(result.Added, func(i, j int) bool { return result.Added[i].Name < result.Added[j].Name })
	sort.Slice(result.Removed, func(i, j int) bool { return result.Removed[i].Name < result.Removed[j].Name })
	sort.Slice(result.Unchanged, func(i, j int) bool { return result.Unchanged[i].Name < result.Unchanged[j].Name })

	return result
}
