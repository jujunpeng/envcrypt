package keystore

import (
	"strings"
	"testing"
)

func TestListRecipientsEmpty(t *testing.T) {
	dir := t.TempDir()
	ks, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	infos := ks.ListRecipients()
	if len(infos) != 0 {
		t.Errorf("expected 0 recipients, got %d", len(infos))
	}
}

func TestListRecipientsSorted(t *testing.T) {
	dir := t.TempDir()
	ks, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	pub1 := generatePublicKey(t)
	pub2 := generatePublicKey(t)
	pub3 := generatePublicKey(t)

	_ = ks.AddRecipient("charlie", pub1)
	_ = ks.AddRecipient("alice", pub2)
	_ = ks.AddRecipient("bob", pub3)

	infos := ks.ListRecipients()
	if len(infos) != 3 {
		t.Fatalf("expected 3 recipients, got %d", len(infos))
	}

	if infos[0].Name != "alice" || infos[1].Name != "bob" || infos[2].Name != "charlie" {
		t.Errorf("recipients not sorted: %v", infos)
	}
}

func TestListRecipientsShortKey(t *testing.T) {
	dir := t.TempDir()
	ks, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	pub := generatePublicKey(t)
	_ = ks.AddRecipient("testuser", pub)

	infos := ks.ListRecipients()
	if len(infos) != 1 {
		t.Fatalf("expected 1 recipient, got %d", len(infos))
	}

	if infos[0].PublicKey != pub {
		t.Errorf("full public key mismatch")
	}
	if len(infos[0].ShortKey) > 24 {
		t.Errorf("short key too long: %s", infos[0].ShortKey)
	}
}

func TestFormatRecipientTableEmpty(t *testing.T) {
	dir := t.TempDir()
	ks, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	out := ks.FormatRecipientTable()
	if !strings.Contains(out, "No recipients") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestFormatRecipientTableWithEntries(t *testing.T) {
	dir := t.TempDir()
	ks, err := New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	pub := generatePublicKey(t)
	_ = ks.AddRecipient("alice", pub)

	out := ks.FormatRecipientTable()
	if !strings.Contains(out, "alice") {
		t.Errorf("expected 'alice' in table output, got: %s", out)
	}
	if !strings.Contains(out, "NAME") {
		t.Errorf("expected header in table output, got: %s", out)
	}
}
