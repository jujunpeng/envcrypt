package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envcrypt/internal/config"
	"envcrypt/internal/keystore"
)

func runExport(cmd *cobra.Command, args []string) error {
	outPath, _ := cmd.Flags().GetString("output")
	if outPath == "" {
		return fmt.Errorf("--output flag is required")
	}

	cfg, err := config.Load(".envcrypt.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks, err := keystore.New(cfg.KeyStorePath)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	if err := ks.Export(outPath); err != nil {
		return fmt.Errorf("export: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Recipients exported to %s\n", outPath)
	return nil
}

func runImport(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("import requires a file path argument")
	}
	inPath := args[0]

	cfg, err := config.Load(".envcrypt.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ks, err := keystore.New(cfg.KeyStorePath)
	if err != nil {
		return fmt.Errorf("open keystore: %w", err)
	}

	count, err := ks.Import(inPath)
	if err != nil {
		return fmt.Errorf("import: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Imported %d recipient(s) from %s\n", count, inPath)
	return nil
}

func registerExportImportCommands(root *cobra.Command) {
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export all recipients to a JSON file",
		RunE:  runExport,
	}
	exportCmd.Flags().StringP("output", "o", "", "Output file path (required)")

	importCmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import recipients from a JSON export file",
		Args:  cobra.ExactArgs(1),
		RunE:  runImport,
	}

	root.AddCommand(exportCmd)
	root.AddCommand(importCmd)
}
