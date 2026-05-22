package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/config"
	"envcrypt/internal/keystore"
)

func runRotate(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: envcrypt rotate <name> <new-public-key>")
	}
	name, newKey := args[0], args[1]

	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks := keystore.New()
	for _, r := range cfg.Recipients {
		if err := ks.AddRecipient(r.Name, r.PublicKey); err != nil {
			return fmt.Errorf("load recipient %q: %w", r.Name, err)
		}
	}

	rec, err := ks.RotateRecipient(name, newKey)
	if err != nil {
		return fmt.Errorf("rotate: %w", err)
	}

	cfg.Recipients = nil
	for _, r := range ks.AllRecipients() {
		cfg.AddRecipient(r.Name, r.PublicKey)
	}

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	cmd.Printf("Rotated key for %q\n  old: %s\n  new: %s\n", name, rec.OldKey[:20]+"...", rec.NewKey[:20]+"...")
	return nil
}

func runRemove(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envcrypt remove-recipient <name>")
	}
	name := args[0]

	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks := keystore.New()
	for _, r := range cfg.Recipients {
		if err := ks.AddRecipient(r.Name, r.PublicKey); err != nil {
			return fmt.Errorf("load recipient %q: %w", r.Name, err)
		}
	}

	if err := ks.RemoveRecipient(name); err != nil {
		return fmt.Errorf("remove: %w", err)
	}

	cfg.Recipients = nil
	for _, r := range ks.AllRecipients() {
		cfg.AddRecipient(r.Name, r.PublicKey)
	}

	if err := config.Save(cfg, ""); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	cmd.Printf("Removed recipient %q\n", name)
	return nil
}

func registerRotateCommands(root *cobra.Command) {
	rotateCmd := &cobra.Command{
		Use:   "rotate <name> <new-public-key>",
		Short: "Rotate the public key for an existing recipient",
		RunE:  runRotate,
	}
	removeCmd := &cobra.Command{
		Use:   "remove-recipient <name>",
		Short: "Remove a recipient from the keystore",
		RunE:  runRemove,
	}
	root.AddCommand(rotateCmd, removeCmd)
}
