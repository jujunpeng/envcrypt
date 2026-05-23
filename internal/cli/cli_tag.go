package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"envcrypt/internal/config"
	"envcrypt/internal/keystore"
)

func runTag(cmd *cobra.Command, args []string) error {
	name := args[0]
	tag := args[1]

	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks, err := keystore.New(cfg.KeystorePath)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	if err := ks.TagRecipient(name, tag); err != nil {
		return fmt.Errorf("tag recipient: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Tagged %q with %q\n", name, tag)
	return nil
}

func runUntag(cmd *cobra.Command, args []string) error {
	name := args[0]
	tag := args[1]

	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks, err := keystore.New(cfg.KeystorePath)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	if err := ks.UntagRecipient(name, tag); err != nil {
		return fmt.Errorf("untag recipient: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Removed tag %q from %q\n", name, tag)
	return nil
}

func runListByTag(cmd *cobra.Command, args []string) error {
	tag := args[0]

	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks, err := keystore.New(cfg.KeystorePath)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	results, err := ks.RecipientsByTag(tag)
	if err != nil {
		return fmt.Errorf("list by tag: %w", err)
	}

	if len(results) == 0 {
		fmt.Fprintf(cmd.OutOrStdout(), "No recipients found with tag %q\n", tag)
		return nil
	}

	names := make([]string, 0, len(results))
	for _, r := range results {
		names = append(names, r.Name)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Recipients tagged %q: %s\n", tag, strings.Join(names, ", "))
	return nil
}

func registerTagCommands(root *cobra.Command) {
	tagCmd := &cobra.Command{
		Use:   "tag <name> <tag>",
		Short: "Add a tag to a recipient",
		Args:  cobra.ExactArgs(2),
		RunE:  runTag,
	}

	untagCmd := &cobra.Command{
		Use:   "untag <name> <tag>",
		Short: "Remove a tag from a recipient",
		Args:  cobra.ExactArgs(2),
		RunE:  runUntag,
	}

	listByTagCmd := &cobra.Command{
		Use:   "list-by-tag <tag>",
		Short: "List recipients with a given tag",
		Args:  cobra.ExactArgs(1),
		RunE:  runListByTag,
	}

	root.AddCommand(tagCmd, untagCmd, listByTagCmd)
}
