package machine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadProfilesErrorsWhenConfigMissing(t *testing.T) {
	store := Store{Path: filepath.Join(t.TempDir(), "missing.json")}

	_, err := store.LoadProfiles()
	if err == nil {
		t.Fatal("LoadProfiles() error = nil, want read error")
	}
	if !strings.Contains(err.Error(), "read machine config") {
		t.Fatalf("LoadProfiles() error = %v, want read error", err)
	}
}

func TestLoadProfilesReadsConfigFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "machines.json")
	content := `{
  "profiles": [
    {
      "name": "Custom 386",
      "description": "Loaded from config",
      "cpu_core": "normal",
      "cpu_type": "386",
      "cycles": "6000",
      "machine": "ega",
      "memory_mb": 4,
      "sound_blaster": "sb2",
      "gus": false,
      "joystick_type": "disabled",
      "xms": true,
      "ems": true,
      "umb": true
    }
  ]
}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := Store{Path: path}
	profiles, err := store.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() returned error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("LoadProfiles() returned %d profiles, want 1", len(profiles))
	}
	if got := profiles[0].Name; got != "Custom 386" {
		t.Fatalf("profiles[0].Name = %q, want %q", got, "Custom 386")
	}
}

func TestLoadProfilesRejectsInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "machines.json")
	content := `{"profiles":[{"name":"","cpu_core":"normal"}]}`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := Store{Path: path}
	_, err := store.LoadProfiles()
	if err == nil {
		t.Fatal("LoadProfiles() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "profile name is required") {
		t.Fatalf("LoadProfiles() error = %v, want profile validation message", err)
	}
}
