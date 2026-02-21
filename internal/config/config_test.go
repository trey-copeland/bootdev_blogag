package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead_InvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".gatorconfig.json")

	if err := os.WriteFile(cfgPath, []byte(`{db_url: postgres://bad }`), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := readConfigAt(cfgPath)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSetUser_PersistsConfig(t *testing.T) {
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, ".gatorconfig.json")

	want := Config{
		DbUrl:           "postgres://example",
		CurrentUserName: "trey",
	}

	if err := writeConfigAt(cfgPath, want); err != nil {
		t.Fatalf("SetUser err: %v", err)
	}

	got, err := readConfigAt(cfgPath)
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if got.CurrentUserName != want.CurrentUserName {
		t.Fatalf("Expected CurrentUserName=%q, got %q", want.CurrentUserName, got.CurrentUserName)
	}
	if got.DbUrl != "postgres://example" {
		t.Fatalf("Expected DbUrl=%q, got %q", want.DbUrl, got.DbUrl)
	}

}
