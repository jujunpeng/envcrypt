package cli

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"envcrypt/internal/config"
	"envcrypt/internal/keystore"
)

const expiryDateLayout = "2006-01-02"

func runSetExpiry(cmd *cobra.Command, args []string) error {
	name := args[0]
	dateStr := args[1]

	expiry, err := time.Parse(expiryDateLayout, dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}

	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks, err := keystore.New(cfg.KeystoreDir)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	if err := ks.SetExpiry(name, expiry.UTC()); err != nil {
		return fmt.Errorf("set expiry: %w", err)
	}

	cmd.Printf("Expiry for %q set to %s\n", name, expiry.Format(expiryDateLayout))
	return nil
}

func runCheckExpiry(cmd *cobra.Command, args []string) error {
	name := args[0]

	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks, err := keystore.New(cfg.KeystoreDir)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	t, set, err := ks.GetExpiry(name)
	if err != nil {
		return err
	}
	if !set {
		cmd.Printf("No expiry set for %q\n", name)
		return nil
	}

	expired, _ := ks.IsExpired(name)
	status := "valid"
	if expired {
		status = "EXPIRED"
	}
	cmd.Printf("Expiry for %q: %s [%s]\n", name, t.Format(expiryDateLayout), status)
	return nil
}

func registerExpiryCommands(root *cobra.Command) {
	expiryCmd := &cobra.Command{
		Use:   "expiry",
		Short: "Manage recipient key expiry dates",
	}

	setCmd := &cobra.Command{
		Use:   "set <name> <YYYY-MM-DD>",
		Short: "Set an expiry date for a recipient key",
		Args:  cobra.ExactArgs(2),
		RunE:  runSetExpiry,
	}

	checkCmd := &cobra.Command{
		Use:   "check <name>",
		Short: "Check the expiry status of a recipient key",
		Args:  cobra.ExactArgs(1),
		RunE:  runCheckExpiry,
	}

	expiryCmd.AddCommand(setCmd, checkCmd)
	root.AddCommand(expiryCmd)
}
