package keystore

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ExportFormat represents the format for exported keystore data.
type ExportFormat struct {
	Version    int               `json:"version"`
	Recipients map[string]string `json:"recipients"`
}

// Export writes all recipients to a JSON file at the given path.
func (ks *KeyStore) Export(path string) error {
	recipients, err := ks.AllRecipients()
	if err != nil {
		return fmt.Errorf("export: failed to list recipients: %w", err)
	}

	data := ExportFormat{
		Version:    1,
		Recipients: recipients,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("export: failed to marshal data: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("export: failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, bytes, 0600); err != nil {
		return fmt.Errorf("export: failed to write file: %w", err)
	}

	return nil
}

// Import reads recipients from a JSON file and merges them into the keystore.
// Existing recipients with the same name are overwritten.
func (ks *KeyStore) Import(path string) (int, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("import: failed to read file: %w", err)
	}

	var data ExportFormat
	if err := json.Unmarshal(bytes, &data); err != nil {
		return 0, fmt.Errorf("import: failed to parse file: %w", err)
	}

	if data.Version != 1 {
		return 0, fmt.Errorf("import: unsupported export version %d", data.Version)
	}

	count := 0
	for name, pubKey := range data.Recipients {
		if err := ks.AddRecipient(name, pubKey); err != nil {
			return count, fmt.Errorf("import: failed to add recipient %q: %w", name, err)
		}
		count++
	}

	return count, nil
}
