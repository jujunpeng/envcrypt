package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"envcrypt/internal/config"
	"envcrypt/internal/keystore"
)

func runDiff(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")

	current, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	compareFile, _ := cmd.Flags().GetString("compare")
	if compareFile == "" {
		return fmt.Errorf("--compare flag is required")
	}

	other, err := config.Load(compareFile)
	if err != nil {
		return fmt.Errorf("loading compare config: %w", err)
	}

	toRecipients := func(cfg *config.Config) []keystore.Recipient {
		var out []keystore.Recipient
		for _, r := range cfg.Recipients {
			out = append(out, keystore.Recipient{Name: r.Name, PublicKey: r.PublicKey})
		}
		return out
	}

	result := keystore.DiffRecipients(toRecipients(current), toRecipients(other))

	if len(result.Added) == 0 && len(result.Removed) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No differences in recipients.")
		return nil
	}

	for _, r := range result.Added {
		fmt.Fprintf(cmd.OutOrStdout(), "+ %-20s %s\n", r.Name, r.PublicKey)
	}
	for _, r := range result.Removed {
		fmt.Fprintf(cmd.OutOrStdout(), "- %-20s %s\n", r.Name, r.PublicKey)
	}

	return nil
}

func registerDiffCommand(root *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "diff",
		Short: "Show recipient differences between two config files",
		RunE:  runDiff,
	}
	cmd.Flags().String("compare", "", "Path to config file to compare against (required)")
	root.AddCommand(cmd)
}
