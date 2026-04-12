package machine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const appConfigDirName = "RigControl"
const machinesConfigName = "machines.json"

type Store struct {
	Path string
}

type fileData struct {
	Profiles []Profile `json:"profiles"`
}

func DefaultStore() (Store, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return Store{}, fmt.Errorf("resolve user config dir: %w", err)
	}

	return Store{
		Path: filepath.Join(configRoot, appConfigDirName, machinesConfigName),
	}, nil
}

func (s Store) LoadProfiles() ([]Profile, error) {
	content, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, fmt.Errorf("read machine config %s: %w", s.Path, err)
	}

	var data fileData
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("parse machine config %s: %w", s.Path, err)
	}
	if len(data.Profiles) == 0 {
		return nil, fmt.Errorf("machine config %s contains no profiles", s.Path)
	}

	profiles, err := validatedProfiles(data.Profiles)
	if err != nil {
		return nil, wrapProfileValidationError("machine config "+s.Path, err)
	}

	return profiles, nil
}

func (s Store) SaveProfiles(profiles []Profile) error {
	if len(profiles) == 0 {
		return fmt.Errorf("no machine profiles to save")
	}

	cloned, err := validatedProfiles(profiles)
	if err != nil {
		return wrapProfileValidationError("save machine config "+s.Path, err)
	}

	payload, err := json.MarshalIndent(fileData{Profiles: cloned}, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal machine config %s: %w", s.Path, err)
	}

	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return fmt.Errorf("create machine config dir for %s: %w", s.Path, err)
	}
	if err := os.WriteFile(s.Path, append(payload, '\n'), 0o644); err != nil {
		return fmt.Errorf("write machine config %s: %w", s.Path, err)
	}

	return nil
}

func cloneProfiles(profiles []Profile) []Profile {
	return append([]Profile(nil), profiles...)
}

func validatedProfiles(profiles []Profile) ([]Profile, error) {
	cloned := cloneProfiles(profiles)
	for i, profile := range cloned {
		if err := ValidateProfile(profile); err != nil {
			return nil, profileError{
				Index: i,
				Name:  profile.Name,
				Err:   err,
			}
		}
	}

	return cloned, nil
}

type profileError struct {
	Index int
	Name  string
	Err   error
}

func (e profileError) Error() string {
	return fmt.Sprintf("profile %d (%q): %v", e.Index, e.Name, e.Err)
}

func (e profileError) Unwrap() error {
	return e.Err
}

func wrapProfileValidationError(context string, err error) error {
	return fmt.Errorf("%s %w", context, err)
}
