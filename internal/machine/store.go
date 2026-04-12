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

	profiles := cloneProfiles(data.Profiles)
	for i, profile := range profiles {
		if err := ValidateProfile(profile); err != nil {
			return nil, fmt.Errorf("machine config %s profile %d (%q): %w", s.Path, i, profile.Name, err)
		}
	}

	return profiles, nil
}

func cloneProfiles(profiles []Profile) []Profile {
	return append([]Profile(nil), profiles...)
}
