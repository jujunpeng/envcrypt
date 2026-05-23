package keystore

import (
	"strings"
)

// SearchResult holds a matched recipient with match metadata.
type SearchResult struct {
	Name      string
	PublicKey string
	MatchedBy string // "name" or "key"
}

// Search finds recipients whose name or public key contains the given query
// (case-insensitive). Returns all matches.
func (s *Store) Search(query string) []SearchResult {
	if query == "" {
		return nil
	}
	q := strings.ToLower(query)
	var results []SearchResult

	s.mu.RLock()
	defer s.mu.RUnlock()

	for name, key := range s.recipients {
		if strings.Contains(strings.ToLower(name), q) {
			results = append(results, SearchResult{
				Name:      name,
				PublicKey: key,
				MatchedBy: "name",
			})
			continue
		}
		if strings.Contains(strings.ToLower(key), q) {
			results = append(results, SearchResult{
				Name:      name,
				PublicKey: key,
				MatchedBy: "key",
			})
		}
	}
	return results
}

// FindByPrefix returns the first recipient whose name starts with prefix
// (case-insensitive). Returns empty string if not found.
func (s *Store) FindByPrefix(prefix string) (name, key string, found bool) {
	if prefix == "" {
		return "", "", false
	}
	p := strings.ToLower(prefix)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for n, k := range s.recipients {
		if strings.HasPrefix(strings.ToLower(n), p) {
			return n, k, true
		}
	}
	return "", "", false
}
