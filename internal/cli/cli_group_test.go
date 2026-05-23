package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envcrypt/internal/keystore"
)

func writeGroupConfig(t *testing.T) (string, *keystore.Store) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.json")
	ks, err := keystore.New(path)
	if err != nil {
		t.Fatalf("keystore.New: %v", err)
	}
	return path, ks
}

func TestGroupAddAndList(t *testing.T) {
	path, ks := writeGroupConfig(t)
	pub1 := addTestRecipient(t, ks, "alice")
	pub2 := addTestRecipient(t, ks, "bob")
	_ = pub1
	_ = pub2

	out, err := withArgs("group", "add", "devs", "alice,bob", "--config", path)
	if err != nil {
		t.Fatalf("group add: %v", err)
	}
	if !strings.Contains(out, "devs") {
		t.Errorf("expected group name in output, got: %s", out)
	}

	out, err = withArgs("group", "list", "--config", path)
	if err != nil {
		t.Fatalf("group list: %v", err)
	}
	if !strings.Contains(out, "devs") {
		t.Errorf("expected devs in list output, got: %s", out)
	}
}

func TestGroupListEmpty(t *testing.T) {
	path, _ := writeGroupConfig(t)
	out, err := withArgs("group", "list", "--config", path)
	if err != nil {
		t.Fatalf("group list: %v", err)
	}
	if !strings.Contains(out, "No groups") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestGroupRemove(t *testing.T) {
	path, ks := writeGroupConfig(t)
	addTestRecipient(t, ks, "alice")
	_, _ = withArgs("group", "add", "devs", "alice", "--config", path)

	_, err := withArgs("group", "remove", "devs", "--config", path)
	if err != nil {
		t.Fatalf("group remove: %v", err)
	}
	out, _ := withArgs("group", "list", "--config", path)
	if strings.Contains(out, "devs") {
		t.Errorf("group should be removed, got: %s", out)
	}
}

func TestGroupAddMissingRecipient(t *testing.T) {
	path, _ := writeGroupConfig(t)
	_, err := withArgs("group", "add", "devs", "ghost", "--config", path)
	if err == nil {
		t.Error("expected error for unknown recipient")
	}
}

func addTestRecipient(t *testing.T, ks *keystore.Store, name string) string {
	t.Helper()
	pub, _ := generatePublicKey(t)
	if err := ks.AddRecipient(name, pub); err != nil {
		t.Fatalf("AddRecipient(%q): %v", name, err)
	}
	return pub
}

func init() {
	_ = os.Getenv
}
