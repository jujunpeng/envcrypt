package keystore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffNoChanges(t *testing.T) {
	recipients := []Recipient{
		{Name: "alice", PublicKey: "age1alice"},
		{Name: "bob", PublicKey: "age1bob"},
	}
	result := DiffRecipients(recipients, recipients)
	assert.Empty(t, result.Added)
	assert.Empty(t, result.Removed)
	assert.Len(t, result.Unchanged, 2)
}

func TestDiffAddedRecipient(t *testing.T) {
	old := []Recipient{{Name: "alice", PublicKey: "age1alice"}}
	new := []Recipient{
		{Name: "alice", PublicKey: "age1alice"},
		{Name: "bob", PublicKey: "age1bob"},
	}
	result := DiffRecipients(old, new)
	assert.Len(t, result.Added, 1)
	assert.Equal(t, "bob", result.Added[0].Name)
	assert.Empty(t, result.Removed)
	assert.Len(t, result.Unchanged, 1)
}

func TestDiffRemovedRecipient(t *testing.T) {
	old := []Recipient{
		{Name: "alice", PublicKey: "age1alice"},
		{Name: "bob", PublicKey: "age1bob"},
	}
	new := []Recipient{{Name: "alice", PublicKey: "age1alice"}}
	result := DiffRecipients(old, new)
	assert.Empty(t, result.Added)
	assert.Len(t, result.Removed, 1)
	assert.Equal(t, "bob", result.Removed[0].Name)
	assert.Len(t, result.Unchanged, 1)
}

func TestDiffChangedKey(t *testing.T) {
	old := []Recipient{{Name: "alice", PublicKey: "age1alice_old"}}
	new := []Recipient{{Name: "alice", PublicKey: "age1alice_new"}}
	result := DiffRecipients(old, new)
	assert.Len(t, result.Added, 1)
	assert.Equal(t, "age1alice_new", result.Added[0].PublicKey)
	assert.Len(t, result.Removed, 1)
	assert.Equal(t, "age1alice_old", result.Removed[0].PublicKey)
	assert.Empty(t, result.Unchanged)
}

func TestDiffEmptyOld(t *testing.T) {
	new := []Recipient{
		{Name: "alice", PublicKey: "age1alice"},
		{Name: "bob", PublicKey: "age1bob"},
	}
	result := DiffRecipients(nil, new)
	assert.Len(t, result.Added, 2)
	assert.Empty(t, result.Removed)
	assert.Empty(t, result.Unchanged)
}

func TestDiffEmptyNew(t *testing.T) {
	old := []Recipient{
		{Name: "alice", PublicKey: "age1alice"},
		{Name: "bob", PublicKey: "age1bob"},
	}
	result := DiffRecipients(old, nil)
	assert.Empty(t, result.Added)
	assert.Len(t, result.Removed, 2)
	assert.Empty(t, result.Unchanged)
}

func TestDiffSortedOutput(t *testing.T) {
	old := []Recipient{}
	new := []Recipient{
		{Name: "zara", PublicKey: "age1zara"},
		{Name: "alice", PublicKey: "age1alice"},
		{Name: "mike", PublicKey: "age1mike"},
	}
	result := DiffRecipients(old, new)
	assert.Equal(t, "alice", result.Added[0].Name)
	assert.Equal(t, "mike", result.Added[1].Name)
	assert.Equal(t, "zara", result.Added[2].Name)
}
