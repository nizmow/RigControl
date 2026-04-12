package ui

import (
	"fmt"
	"strings"

	"rigcontrol/internal/machine"
)

func profileSummary(profile machine.Profile) string {
	lines := []string{
		fmt.Sprintf("CPU: %s / %s", profile.CPUCore, profile.CPUType),
		fmt.Sprintf("Cycles: %s", profile.Cycles),
		fmt.Sprintf("Video: %s", profile.Machine),
		fmt.Sprintf("Memory: %d MB", profile.MemoryMB),
		fmt.Sprintf("Sound Blaster: %s", profile.SoundBlaster),
		fmt.Sprintf("Joystick: %s", profile.JoystickType),
		fmt.Sprintf("GUS: %s", onOff(profile.GUS)),
		fmt.Sprintf("XMS / EMS / UMB: %s / %s / %s", yesNo(profile.XMS), yesNo(profile.EMS), yesNo(profile.UMB)),
	}
	if profile.HardDiskImage != "" {
		lines = append(lines,
			fmt.Sprintf("HDD Image: %s", profile.HardDiskImage),
			fmt.Sprintf("HDD CHS: %s", profile.HardDiskCHS),
		)
	}

	return strings.Join(lines, "\n")
}

func yesNo(enabled bool) string {
	if enabled {
		return "yes"
	}
	return "no"
}

func onOff(enabled bool) string {
	if enabled {
		return "on"
	}
	return "off"
}
