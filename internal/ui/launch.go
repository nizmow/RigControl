package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"rigcontrol/internal/dosbox"
	"rigcontrol/internal/machine"
)

const dosboxExecutable = "/Applications/DOSBox Staging.app/Contents/MacOS/dosbox"

func writeTempConfig(profile machine.Profile) (string, error) {
	slug := strings.ToLower(profile.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	if slug == "" {
		slug = "rigcontrol"
	}

	file, err := os.CreateTemp("", slug+"-*.conf")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(dosbox.Render(profile)); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func launchDOSBox(configPath string) error {
	if _, err := os.Stat(dosboxExecutable); err != nil {
		return fmt.Errorf("dosbox executable not found at %s: %w", dosboxExecutable, err)
	}

	cmd := exec.Command(dosboxExecutable, "-conf", configPath)
	return cmd.Start()
}

func launchProfile(profile machine.Profile, configPath string) error {
	if _, err := os.Stat(dosboxExecutable); err != nil {
		return fmt.Errorf("dosbox executable not found at %s: %w", dosboxExecutable, err)
	}
	if profile.HardDiskImage != "" {
		if _, err := os.Stat(profile.HardDiskImage); err != nil {
			return fmt.Errorf("hard disk image not found at %s: %w", profile.HardDiskImage, err)
		}
	}

	cmd := exec.Command(dosboxExecutable, "-conf", configPath)
	return cmd.Start()
}
