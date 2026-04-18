package machine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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
	if got := profiles[0].MouseCapture; got != "onclick" {
		t.Fatalf("profiles[0].MouseCapture = %q, want %q", got, "onclick")
	}
	if got := profiles[0].MouseRawInput; got != true {
		t.Fatalf("profiles[0].MouseRawInput = %v, want true", got)
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

func TestSaveProfilesWritesConfigFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "machines.json")
	store := Store{Path: path}
	profiles := []Profile{
		{
			Name:          "Saved 486",
			Description:   "Persisted machine",
			CPUCore:       "auto",
			CPUType:       "486",
			Cycles:        "25000",
			Machine:       "svga_s3",
			MemoryMB:      16,
			SoundBlaster:  "sb16",
			MouseCapture:  "onclick",
			MouseRawInput: true,
			GUS:           false,
			JoystickType:  "auto",
			XMS:           true,
			EMS:           true,
			UMB:           true,
		},
	}

	if err := store.SaveProfiles(profiles); err != nil {
		t.Fatalf("SaveProfiles() error = %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(content), `"name": "Saved 486"`) {
		t.Fatalf("saved content missing expected profile: %s", string(content))
	}
}

func TestSaveProfilesRejectsInvalidProfiles(t *testing.T) {
	path := filepath.Join(t.TempDir(), "machines.json")
	store := Store{Path: path}

	err := store.SaveProfiles([]Profile{{}})
	if err == nil {
		t.Fatal("SaveProfiles() error = nil, want validation error")
	}
	if !strings.Contains(err.Error(), "profile name is required") {
		t.Fatalf("SaveProfiles() error = %v, want validation message", err)
	}
}

func TestSaveProfilesThenLoadProfilesRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "machines.json")
	store := Store{Path: path}
	want := []Profile{
		{
			Name:          "Round Trip Machine",
			Description:   "Stored and loaded",
			CPUCore:       "dynamic",
			CPUType:       "pentium",
			Cycles:        "50000",
			Machine:       "svga_s3",
			MemoryMB:      32,
			SoundBlaster:  "sb16",
			MouseCapture:  "onclick",
			MouseRawInput: true,
			GUS:           true,
			JoystickType:  "auto",
			XMS:           true,
			EMS:           true,
			UMB:           false,
		},
	}

	if err := store.SaveProfiles(want); err != nil {
		t.Fatalf("SaveProfiles() error = %v", err)
	}

	got, err := store.LoadProfiles()
	if err != nil {
		t.Fatalf("LoadProfiles() error = %v", err)
	}

	if len(got) != len(want) {
		t.Fatalf("LoadProfiles() returned %d profiles, want %d", len(got), len(want))
	}
	if !reflect.DeepEqual(got[0], want[0]) {
		t.Fatalf("LoadProfiles() profile = %#v, want %#v", got[0], want[0])
	}
}

func TestSaveProfilesCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "config", "machines.json")
	store := Store{Path: path}

	err := store.SaveProfiles([]Profile{
		{
			Name:          "Nested Save",
			CPUCore:       "auto",
			CPUType:       "486",
			Cycles:        "25000",
			Machine:       "svga_s3",
			MemoryMB:      16,
			SoundBlaster:  "sb16",
			MouseCapture:  "onclick",
			MouseRawInput: true,
			JoystickType:  "auto",
		},
	})
	if err != nil {
		t.Fatalf("SaveProfiles() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("saved file missing at %s: %v", path, err)
	}
}

func TestLoadProfilesRejectsEmptyProfilesList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "machines.json")
	payload, err := json.Marshal(fileData{Profiles: []Profile{}})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := Store{Path: path}
	_, err = store.LoadProfiles()
	if err == nil {
		t.Fatal("LoadProfiles() error = nil, want empty profile error")
	}
	if !strings.Contains(err.Error(), "contains no profiles") {
		t.Fatalf("LoadProfiles() error = %v, want empty profile message", err)
	}
}

func TestLoadProfilesOrCreateSeedsMissingConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "machines.json")
	store := Store{Path: path}
	initialProfiles := StarterProfiles()

	got, err := store.LoadProfilesOrCreate(initialProfiles)
	if err != nil {
		t.Fatalf("LoadProfilesOrCreate() error = %v", err)
	}
	if len(got) != len(initialProfiles) {
		t.Fatalf("LoadProfilesOrCreate() returned %d profiles, want %d", len(got), len(initialProfiles))
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if !strings.Contains(string(content), `"name": "286 EGA"`) {
		t.Fatalf("saved content missing expected starter profile: %s", string(content))
	}
}
