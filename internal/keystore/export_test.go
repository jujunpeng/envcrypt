package keystore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExportAndImport(t *testing.T) {
	dir := t.TempDir()
	ksPath := filepath.Join(dir, "keys.json")
	exportPath := filepath.Join(dir, "export.json")

	ks, err := New(ksPath)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	pub1 := generatePublicKey(t)
	pub2 := generatePublicKey(t)

	_ = ks.AddRecipient("alice", pub1)
	_ = ks.AddRecipient("bob", pub2)

	if err := ks.Export(exportPath); err != nil {
		t.Fatalf("Export: %v", err)
	}

	raw, err := os.ReadFile(exportPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var ef ExportFormat
	if err := json.Unmarshal(raw, &ef); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if ef.Version != 1 {
		t.Errorf("expected version 1, got %d", ef.Version)
	}
	if len(ef.Recipients) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(ef.Recipients))
	}

	ks2Path := filepath.Join(dir, "keys2.json")
	ks2, err := New(ks2Path)
	if err != nil {
		t.Fatalf("New ks2: %v", err)
	}

	count, err := ks2.Import(exportPath)
	if err != nil {
		t.Fatalf("Import: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 imported, got %d", count)
	}

	got, err := ks2.GetRecipient("alice")
	if err != nil {
		t.Fatalf("GetRecipient alice: %v", err)
	}
	if got != pub1 {
		t.Errorf("alice key mismatch")
	}
}

func TestImportUnsupportedVersion(t *testing.T) {
	dir := t.TempDir()
	exportPath := filepath.Join(dir, "bad.json")

	data := `{"version":99,"recipients":{}}`
	_ = os.WriteFile(exportPath, []byte(data), 0600)

	ks, _ := New(filepath.Join(dir, "keys.json"))
	_, err := ks.Import(exportPath)
	if err == nil {
		t.Fatal("expected error for unsupported version")
	}
}

func TestImportMissingFile(t *testing.T) {
	dir := t.TempDir()
	ks, _ := New(filepath.Join(dir, "keys.json"))
	_, err := ks.Import(filepath.Join(dir, "nonexistent.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestExportCreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	exportPath := filepath.Join(dir, "subdir", "nested", "export.json")

	ks, _ := New(filepath.Join(dir, "keys.json"))
	if err := ks.Export(exportPath); err != nil {
		t.Fatalf("Export to nested path: %v", err)
	}

	if _, err := os.Stat(exportPath); err != nil {
		t.Errorf("exported file not found: %v", err)
	}
}
