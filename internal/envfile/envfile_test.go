package envfile

import (
	"os"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseValidEnvFile(t *testing.T) {
	path := writeTempEnv(t, "# comment\nDB_HOST=localhost\nDB_PORT=5432\n\nSECRET=abc123\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(env.Entries))
	}
	if env.Entries[0].Key != "DB_HOST" || env.Entries[0].Value != "localhost" {
		t.Errorf("unexpected first entry: %+v", env.Entries[0])
	}
}

func TestParseInvalidLine(t *testing.T) {
	path := writeTempEnv(t, "INVALID_LINE_WITHOUT_EQUALS\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParseFileNotFound(t *testing.T) {
	_, err := Parse("/nonexistent/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestSerialize(t *testing.T) {
	env := &EnvFile{
		Entries: []Entry{
			{Key: "FOO", Value: "bar"},
			{Key: "BAZ", Value: "qux"},
		},
	}
	got := env.Serialize()
	want := "FOO=bar\nBAZ=qux\n"
	if got != want {
		t.Errorf("Serialize() = %q, want %q", got, want)
	}
}

func TestGetAndSet(t *testing.T) {
	env := &EnvFile{}
	env.Set("KEY", "value1")
	v, ok := env.Get("KEY")
	if !ok || v != "value1" {
		t.Errorf("Get(KEY) = %q, %v; want %q, true", v, ok, "value1")
	}
	env.Set("KEY", "value2")
	v, ok = env.Get("KEY")
	if !ok || v != "value2" {
		t.Errorf("Get(KEY) after update = %q, %v; want %q, true", v, ok, "value2")
	}
	if len(env.Entries) != 1 {
		t.Errorf("expected 1 entry after update, got %d", len(env.Entries))
	}
}

func TestGetMissingKey(t *testing.T) {
	env := &EnvFile{}
	_, ok := env.Get("MISSING")
	if ok {
		t.Error("expected Get to return false for missing key")
	}
}
