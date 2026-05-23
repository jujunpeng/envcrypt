package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Snapshot represents a point-in-time copy of the keystore recipients.
type Snapshot struct {
	Version    int                  `json:"version"`
	CreatedAt  time.Time            `json:"created_at"`
	Label      string               `json:"label,omitempty"`
	Recipients map[string]Recipient `json:"recipients"`
}

// TakeSnapshot saves the current state of the keystore to a snapshot file.
func (s *Store) TakeSnapshot(dir, label string) (string, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("create snapshot dir: %w", err)
	}

	recipients := make(map[string]Recipient)
	for k, v := range s.recipients {
		recipients[k] = v
	}

	snap := Snapshot{
		Version:    1,
		CreatedAt:  time.Now().UTC(),
		Label:      label,
		Recipients: recipients,
	}

	filename := fmt.Sprintf("snapshot-%d.json", snap.CreatedAt.UnixNano())
	path := filepath.Join(dir, filename)

	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal snapshot: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return "", fmt.Errorf("write snapshot: %w", err)
	}

	return path, nil
}

// LoadSnapshot reads a snapshot from disk.
func LoadSnapshot(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read snapshot: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("parse snapshot: %w", err)
	}

	return &snap, nil
}

// RestoreSnapshot replaces the store's recipients with those from the snapshot.
func (s *Store) RestoreSnapshot(snap *Snapshot) {
	s.recipients = make(map[string]Recipient)
	for k, v := range snap.Recipients {
		s.recipients[k] = v
	}
}
