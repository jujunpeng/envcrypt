package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// EventType represents the type of audit event.
type EventType string

const (
	EventEncrypt      EventType = "encrypt"
	EventDecrypt      EventType = "decrypt"
	EventAddRecipient EventType = "add_recipient"
	EventRotate       EventType = "rotate"
	EventRemove       EventType = "remove"
	EventExport       EventType = "export"
	EventImport       EventType = "import"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Actor     string    `json:"actor,omitempty"`
	Details   string    `json:"details"`
	Success   bool      `json:"success"`
}

// Logger writes audit entries to a log file.
type Logger struct {
	path string
}

// New creates a new audit Logger writing to the given file path.
func New(path string) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("audit: create log directory: %w", err)
	}
	return &Logger{path: path}, nil
}

// Log appends an audit entry to the log file.
func (l *Logger) Log(event EventType, details string, success bool) error {
	return l.LogAs("", event, details, success)
}

// LogAs appends an audit entry with an explicit actor.
func (l *Logger) LogAs(actor string, event EventType, details string, success bool) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Actor:     actor,
		Details:   details,
		Success:   success,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// ReadAll reads all audit entries from the log file.
func (l *Logger) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read log: %w", err)
	}
	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
