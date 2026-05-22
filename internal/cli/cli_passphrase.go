package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/spf13/cobra"

	"github.com/yourorg/envcrypt/internal/keystore"
)

// registerPassphraseCommands attaches passphrase-related sub-commands to root.
func registerPassphraseCommands(root *cobra.Command) {
	root.AddCommand(newProtectKeyCmd())
	root.AddCommand(newUnlockKeyCmd())
}

// newProtectKeyCmd returns the `envcrypt protect-key` command.
func newProtectKeyCmd() *cobra.Command {
	var (
		keyPath    string
		passphrase string
		outPath    string
	)

	cmd := &cobra.Command{
		Use:   "protect-key",
		Short: "Encrypt a plain-text age private key with a passphrase",
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyPath == "" {
				return errors.New("--key is required")
			}
			if passphrase == "" {
				return errors.New("--passphrase is required")
			}

			raw, err := os.ReadFile(keyPath)
			if err != nil {
				return fmt.Errorf("reading key file: %w", err)
			}

			id, err := age.ParseX25519Identity(string(raw))
			if err != nil {
				return fmt.Errorf("parsing identity: %w", err)
			}

			if outPath == "" {
				outPath = keyPath + ".enc"
			}

			if err := os.MkdirAll(filepath.Dir(outPath), 0700); err != nil {
				return fmt.Errorf("creating output directory: %w", err)
			}

			if err := keystore.SaveIdentityEncrypted(outPath, id, passphrase); err != nil {
				return fmt.Errorf("saving encrypted identity: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Protected key written to %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&keyPath, "key", "", "Path to plain-text age private key file")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "Passphrase used to encrypt the key")
	cmd.Flags().StringVar(&outPath, "out", "", "Output path (default: <key>.enc)")
	return cmd
}

// newUnlockKeyCmd returns the `envcrypt unlock-key` command.
func newUnlockKeyCmd() *cobra.Command {
	var (
		keyPath    string
		passphrase string
		outPath    string
	)

	cmd := &cobra.Command{
		Use:   "unlock-key",
		Short: "Decrypt a passphrase-protected age private key to plain text",
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyPath == "" {
				return errors.New("--key is required")
			}
			if passphrase == "" {
				return errors.New("--passphrase is required")
			}

			id, err := keystore.LoadIdentityEncrypted(keyPath, passphrase)
			if err != nil {
				return fmt.Errorf("decrypting identity: %w", err)
			}

			if outPath == "" {
				outPath = keyPath + ".plain"
			}

			if err := os.WriteFile(outPath, []byte(id.String()), 0600); err != nil {
				return fmt.Errorf("writing plain key: %w", err)
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Plain-text key written to %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().StringVar(&keyPath, "key", "", "Path to encrypted age private key file")
	cmd.Flags().StringVar(&passphrase, "passphrase", "", "Passphrase used to decrypt the key")
	cmd.Flags().StringVar(&outPath, "out", "", "Output path (default: <key>.plain)")
	return cmd
}
