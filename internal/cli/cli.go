package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/config"
	"envcrypt/internal/keystore"
	"envcrypt/internal/sync"
)

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "envcrypt",
		Short: "Securely encrypt and sync .env files using age encryption",
	}
	registerInitCommand(root)
	registerEncryptCommand(root)
	registerDecryptCommand(root)
	registerAddRecipientCommand(root)
	registerExportImportCommands(root)
	registerRotateCommands(root)
	registerAuditCommands(root)
	registerPassphraseCommands(root)
	registerSearchCommand(root)
	return root
}

// Execute is the CLI entry point.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func registerInitCommand(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialise envcrypt in the current directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cmd, args)
		},
	})
}

func registerEncryptCommand(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt the .env file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runEncrypt(cmd, args)
		},
	})
}

func registerDecryptCommand(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt the .env.enc file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDecrypt(cmd, args)
		},
	})
}

func registerAddRecipientCommand(root *cobra.Command) {
	root.AddCommand(&cobra.Command{
		Use:   "add-recipient <name> <public-key>",
		Short: "Add a recipient to the keystore",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddRecipient(cmd, args)
		},
	})
}

func runInit(_ *cobra.Command, _ []string) error {
	cfg := config.DefaultConfig()
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	fmt.Println("Initialised envcrypt")
	return nil
}

func runEncrypt(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	ks := keystore.New()
	for _, r := range cfg.Recipients {
		if addErr := ks.AddRecipient(r.Name, r.PublicKey); addErr != nil {
			return addErr
		}
	}
	s := sync.New(cfg, ks)
	return s.Encrypt()
}

func runDecrypt(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	ks := keystore.New()
	s := sync.New(cfg, ks)
	return s.Decrypt()
}

func runAddRecipient(_ *cobra.Command, args []string) error {
	name, pubKey := args[0], args[1]
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	if addErr := config.AddRecipient(cfg, name, pubKey); addErr != nil {
		return addErr
	}
	return config.Save(cfg)
}
