package machine

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateProfile(profile Profile) error {
	if strings.TrimSpace(profile.Name) == "" {
		return errors.New("profile name is required")
	}
	if strings.TrimSpace(profile.CPUCore) == "" || strings.TrimSpace(profile.CPUType) == "" || strings.TrimSpace(profile.Cycles) == "" {
		return errors.New("cpu settings are required")
	}
	if strings.TrimSpace(profile.Machine) == "" {
		return errors.New("video machine is required")
	}
	if strings.TrimSpace(profile.SoundBlaster) == "" {
		return errors.New("sound blaster type is required")
	}
	if strings.TrimSpace(profile.JoystickType) == "" {
		return errors.New("joystick type is required")
	}
	if profile.MemoryMB <= 0 {
		return fmt.Errorf("memory must be a positive integer")
	}

	return nil
}
