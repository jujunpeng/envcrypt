package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/yourorg/envcrypt/internal/config"
	"github.com/yourorg/envcrypt/internal/sync"
)

const usage = `envcrypt - secure .env encryption using age

Usage:
  envcrypt <command> [arguments]

Commands:
  init                  Initialize envcrypt in the current directory
  encrypt               Encrypt the .env file
  decrypt               Decrypt the encrypted .env file
  add-recipient <key>   Add a recipient public key
`

// Execute parses os.Args and dispatches to the appropriate command.
func Execute() error {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		return nil
	}

	switch os.Args[1] {
	case "init":
		return runInit()
	case "encrypt":
		return runEncrypt()
	case "decrypt":
		return runDecrypt()
	case "add-recipient":
		if len(os.Args) < 3 {
			return errors.New("add-recipient requires a public key argument")
		}
		return runAddRecipient(os.Args[2])
	case "help", "--help", "-h":
		fmt.Print(usage)
		return nil
	default:
		return fmt.Errorf("unknown command: %s", os.Args[1])
	}
}

func runInit() error {
	cfg := config.DefaultConfig()
	if err := config.Save(cfg, ".envcrypt.json"); err != nil {
		return fmt.Errorf("init: %w", err)
	}
	fmt.Println("Initialized envcrypt — created .envcrypt.json")
	return nil
}

func runEncrypt() error {
	cfg, err := config.Load(".envcrypt.json")
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	s, err := sync.New(cfg)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	if err := s.Encrypt(); err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	fmt.Printf("Encrypted %s -> %s\n", cfg.EnvFile, cfg.EncryptedFile)
	return nil
}

func runDecrypt() error {
	cfg, err := config.Load(".envcrypt.json")
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	s, err := sync.New(cfg)
	if err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	if err := s.Decrypt(); err != nil {
		return fmt.Errorf("decrypt: %w", err)
	}
	fmt.Printf("Decrypted %s -> %s\n", cfg.EncryptedFile, cfg.EnvFile)
	return nil
}

func runAddRecipient(pubKey string) error {
	cfg, err := config.Load(".envcrypt.json")
	if err != nil {
		return fmt.Errorf("add-recipient: %w", err)
	}
	if err := cfg.AddRecipient(pubKey); err != nil {
		return fmt.Errorf("add-recipient: %w", err)
	}
	if err := config.Save(cfg, ".envcrypt.json"); err != nil {
		return fmt.Errorf("add-recipient: %w", err)
	}
	fmt.Printf("Added recipient: %s\n", pubKey)
	return nil
}
