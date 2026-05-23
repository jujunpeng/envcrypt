package keystore

import "fmt"

// MergeResult describes the outcome of merging two keystores.
type MergeResult struct {
	Added    []string
	Skipped  []string
	Conflict []string
}

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// MergeSkipConflicts leaves existing recipients unchanged when a name collision occurs.
	MergeSkipConflicts MergeStrategy = iota
	// MergeOverwriteConflicts replaces existing recipients with the incoming ones.
	MergeOverwriteConflicts
)

// MergeRecipients merges recipients from src into dst using the given strategy.
// It returns a MergeResult summarising what happened.
func MergeRecipients(dst, src *Store, strategy MergeStrategy) (*MergeResult, error) {
	if dst == nil || src == nil {
		return nil, fmt.Errorf("merge: dst and src stores must not be nil")
	}

	result := &MergeResult{}

	incoming, err := src.AllRecipients()
	if err != nil {
		return nil, fmt.Errorf("merge: reading source recipients: %w", err)
	}

	for _, r := range incoming {
		existing, _ := dst.GetRecipient(r.Name)
		if existing != nil {
			if strategy == MergeOverwriteConflicts {
				if err := dst.AddRecipient(r.Name, r.PublicKey); err != nil {
					return nil, fmt.Errorf("merge: overwriting recipient %q: %w", r.Name, err)
				}
				result.Conflict = append(result.Conflict, r.Name)
			} else {
				result.Skipped = append(result.Skipped, r.Name)
			}
			continue
		}
		if err := dst.AddRecipient(r.Name, r.PublicKey); err != nil {
			return nil, fmt.Errorf("merge: adding recipient %q: %w", r.Name, err)
		}
		result.Added = append(result.Added, r.Name)
	}

	return result, nil
}
