package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envcrypt/internal/config"
	"envcrypt/internal/keystore"
)

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks := keystore.New()
	for _, r := range cfg.Recipients {
		if addErr := ks.AddRecipient(r.Name, r.PublicKey); addErr != nil {
			return fmt.Errorf("load recipient %s: %w", r.Name, addErr)
		}
	}

	results := ks.Search(query)
	if len(results) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No recipients matched %q\n", query)
		return nil
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tPUBLIC KEY\tMATCH")
	for _, r := range results {
		short := r.PublicKey
		if len(short) > 20 {
			short = short[:20] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", r.Name, short, r.MatchedBy)
	}
	w.Flush()
	return nil
}

func registerSearchCommand(root *cobra.Command) {
	searchCmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search recipients by name or public key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runSearch(cmd, args); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", err)
				return err
			}
			return nil
		},
	}
	root.AddCommand(searchCmd)
}
