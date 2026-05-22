package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key   string
	Value string
}

// EnvFile holds the parsed contents of a .env file.
type EnvFile struct {
	Entries []Entry
}

// Parse reads and parses an .env file from the given path.
// Lines starting with '#' and empty lines are ignored.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envfile: open %q: %w", path, err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("envfile: line %d: invalid format %q", lineNum, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("envfile: line %d: empty key", lineNum)
		}
		entries = append(entries, Entry{Key: key, Value: value})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envfile: scan: %w", err)
	}
	return &EnvFile{Entries: entries}, nil
}

// Serialize converts the EnvFile back to its string representation.
func (e *EnvFile) Serialize() string {
	var sb strings.Builder
	for _, entry := range e.Entries {
		sb.WriteString(entry.Key)
		sb.WriteString("=")
		sb.WriteString(entry.Value)
		sb.WriteString("\n")
	}
	return sb.String()
}

// Get returns the value for the given key, and whether it was found.
func (e *EnvFile) Get(key string) (string, bool) {
	for _, entry := range e.Entries {
		if entry.Key == key {
			return entry.Value, true
		}
	}
	return "", false
}

// Set adds or updates the value for the given key.
func (e *EnvFile) Set(key, value string) {
	for i, entry := range e.Entries {
		if entry.Key == key {
			e.Entries[i].Value = value
			return
		}
	}
	e.Entries = append(e.Entries, Entry{Key: key, Value: value})
}
