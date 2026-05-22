package keystore

import (
	"fmt"
	"sort"
	"strings"
)

// RecipientInfo holds display information about a recipient.
type RecipientInfo struct {
	Name      string
	PublicKey string
	ShortKey  string
}

// ListRecipients returns a sorted slice of RecipientInfo for all recipients
// in the keystore, suitable for display in CLI output.
func (ks *Keystore) ListRecipients() []RecipientInfo {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	infos := make([]RecipientInfo, 0, len(ks.recipients))
	for name, pubKey := range ks.recipients {
		infos = append(infos, RecipientInfo{
			Name:      name,
			PublicKey: pubKey,
			ShortKey:  shortKey(pubKey),
		})
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return infos
}

// FormatRecipientTable returns a human-readable table string of all recipients.
func (ks *Keystore) FormatRecipientTable() string {
	infos := ks.ListRecipients()
	if len(infos) == 0 {
		return "No recipients configured.\n"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-20s  %s\n", "NAME", "PUBLIC KEY (short)"))
	sb.WriteString(strings.Repeat("-", 52) + "\n")
	for _, info := range infos {
		sb.WriteString(fmt.Sprintf("%-20s  %s\n", info.Name, info.ShortKey))
	}
	return sb.String()
}

// shortKey returns a truncated version of a public key for display purposes.
func shortKey(pubKey string) string {
	const maxLen = 24
	if len(pubKey) <= maxLen {
		return pubKey
	}
	return pubKey[:8] + "..." + pubKey[len(pubKey)-8:]
}
