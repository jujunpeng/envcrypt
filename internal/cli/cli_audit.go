package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"envcrypt/internal/audit"
	"envcrypt/internal/config"
)

func runAuditLog(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	if cfgPath == "" {
		cfgPath = config.DefaultFileName
	}
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	logPath := cfg.AuditLog
	if logPath == "" {
		logPath = ".envcrypt-audit.log"
	}
	l, err := audit.New(logPath)
	if err != nil {
		return fmt.Errorf("open audit log: %w", err)
	}
	entries, err := l.ReadAll()
	if err != nil {
		return fmt.Errorf("read audit log: %w", err)
	}
	if len(entries) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "No audit entries found.")
		return nil
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tEVENT\tACTOR\tDETAILS\tSUCCESS")
	for _, e := range entries {
		actor := e.Actor
		if actor == "" {
			actor = "-"
		}
		success := "yes"
		if !e.Success {
			success = "no"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02 15:04:05"),
			string(e.Event),
			actor,
			e.Details,
			success,
		)
	}
	return w.Flush()
}

func registerAuditCommands(root *cobra.Command) {
	auditCmd := &cobra.Command{
		Use:   "audit",
		Short: "View the audit log of envcrypt operations",
		RunE:  runAuditLog,
	}
	auditCmd.Flags().String("config", "", "path to config file (default: .envcrypt.json)")
	root.AddCommand(auditCmd)
}
