package cli

import (
	"os"
	"testing"
)

func withArgs(args []string, fn func()) {
	orig := os.Args
	os.Args = args
	defer func() { os.Args = orig }()
	fn()
}

func TestExecuteHelp(t *testing.T) {
	withArgs([]string{"envcrypt", "help"}, func() {
		if err := Execute(); err != nil {
			t.Fatalf("expected no error for help, got: %v", err)
		}
	})
}

func TestExecuteNoArgs(t *testing.T) {
	withArgs([]string{"envcrypt"}, func() {
		if err := Execute(); err != nil {
			t.Fatalf("expected no error for no args, got: %v", err)
		}
	})
}

func TestExecuteUnknownCommand(t *testing.T) {
	withArgs([]string{"envcrypt", "foobar"}, func() {
		err := Execute()
		if err == nil {
			t.Fatal("expected error for unknown command")
		}
	})
}

func TestExecuteAddRecipientMissingKey(t *testing.T) {
	withArgs([]string{"envcrypt", "add-recipient"}, func() {
		err := Execute()
		if err == nil {
			t.Fatal("expected error when key argument is missing")
		}
	})
}

func TestRunInitCreatesConfig(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	if err := runInit(); err != nil {
		t.Fatalf("runInit: %v", err)
	}
	if _, err := os.Stat(".envcrypt.json"); err != nil {
		t.Fatalf(".envcrypt.json not created: %v", err)
	}
}

func TestRunEncryptMissingConfig(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	err := runEncrypt()
	if err == nil {
		t.Fatal("expected error when config file is missing")
	}
}

func TestRunDecryptMissingConfig(t *testing.T) {
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	err := runDecrypt()
	if err == nil {
		t.Fatal("expected error when config file is missing")
	}
}
