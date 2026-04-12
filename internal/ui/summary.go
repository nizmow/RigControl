package ui

import (
	"fmt"
	"strings"

	"rigcontrol/internal/machine"
)

type summaryLine struct {
	Label  string
	Value  string
	IsPath bool
}

func profileSummaryLines(profile machine.Profile) []summaryLine {
	lines := []summaryLine{
		{Label: "CPU", Value: fmt.Sprintf("%s / %s", profile.CPUCore, profile.CPUType)},
		{Label: "Cycles", Value: profile.Cycles},
		{Label: "Video", Value: profile.Machine},
		{Label: "Memory", Value: fmt.Sprintf("%d MB", profile.MemoryMB)},
		{Label: "Sound Blaster", Value: profile.SoundBlaster},
		{Label: "Mouse Capture", Value: profile.MouseCapture},
		{Label: "Mouse Raw Input", Value: yesNo(profile.MouseRawInput)},
		{Label: "Mouse Immediate", Value: yesNo(profile.DOSMouseImmediate)},
		{Label: "Joystick", Value: profile.JoystickType},
		{Label: "GUS", Value: onOff(profile.GUS)},
		{Label: "XMS / EMS / UMB", Value: fmt.Sprintf("%s / %s / %s", yesNo(profile.XMS), yesNo(profile.EMS), yesNo(profile.UMB))},
	}
	if len(profile.FloppyDiskImages) > 0 {
		for i, path := range profile.FloppyDiskImages {
			lines = append(lines, summaryLine{
				Label:  fmt.Sprintf("Floppy %d", i+1),
				Value:  path,
				IsPath: true,
			})
		}
	}
	if profile.HardDiskImage != "" {
		lines = append(lines,
			summaryLine{Label: "HDD Image", Value: profile.HardDiskImage, IsPath: true},
			summaryLine{Label: "HDD CHS", Value: profile.HardDiskCHS},
		)
	}

	return lines
}

func profileSummary(profile machine.Profile) string {
	lines := profileSummaryLines(profile)
	parts := make([]string, 0, len(lines))
	for _, line := range lines {
		parts = append(parts, fmt.Sprintf("%s: %s", line.Label, line.Value))
	}
	return strings.Join(parts, "\n")
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
