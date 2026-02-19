package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRead_InvalidJSON(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfgPath := filepath.Join(home, ".gatorconfig.json")
	if err := os.WriteFile(cfgPath, []byte(`{db_url: postgres://bad }`), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := Read()
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestSetUser_PersistsConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg := Config{
		DbUrl: "postgres://example",
	}

	if err := cfg.SetUser("trey"); err != nil {
		t.Fatalf("SetUser err: %v", err)
	}

	got, err := Read()
	if err != nil {
		t.Fatalf("Read error: %v", err)
	}

	if got.CurrentUserName != "trey" {
		t.Fatalf("Expected CurrentUserName=trey, got %q", got.CurrentUserName)
	}
	if got.DbUrl != "postgres://example" {
		t.Fatalf("Expected DbUrl=postgres://example, got %q", got.DbUrl)
	}

}
