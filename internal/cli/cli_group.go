package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func runAddGroup(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	ks, err := loadKeystore(cfgPath)
	if err != nil {
		return err
	}
	groupName := args[0]
	members := strings.Split(args[1], ",")
	for i, m := range members {
		members[i] = strings.TrimSpace(m)
	}
	if err := ks.AddGroup(groupName, members); err != nil {
		return fmt.Errorf("add-group: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Group %q created with %d member(s).\n", groupName, len(members))
	return nil
}

func runRemoveGroup(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	ks, err := loadKeystore(cfgPath)
	if err != nil {
		return err
	}
	if err := ks.RemoveGroup(args[0]); err != nil {
		return fmt.Errorf("remove-group: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Group %q removed.\n", args[0])
	return nil
}

func runListGroups(cmd *cobra.Command, _ []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	ks, err := loadKeystore(cfgPath)
	if err != nil {
		return err
	}
	groups := ks.ListGroups()
	if len(groups) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No groups defined.")
		return nil
	}
	for _, g := range groups {
		fmt.Fprintf(cmd.OutOrStdout(), "%-20s %s\n", g.Name, strings.Join(g.Members, ", "))
	}
	return nil
}

func registerGroupCommands(root *cobra.Command) {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage recipient groups",
	}

	addCmd := &cobra.Command{
		Use:   "add <name> <member1,member2,...>",
		Short: "Create a new group",
		Args:  cobra.ExactArgs(2),
		RunE:  runAddGroup,
	}

	removeCmd := &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a group",
		Args:  cobra.ExactArgs(1),
		RunE:  runRemoveGroup,
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		Args:  cobra.NoArgs,
		RunE:  runListGroups,
	}

	groupCmd.AddCommand(addCmd, removeCmd, listCmd)
	root.AddCommand(groupCmd)
}
