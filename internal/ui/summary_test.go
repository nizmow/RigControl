package ui

import (
	"strings"
	"testing"

	"rigcontrol/internal/machine"
)

func TestProfileSummaryIncludesExpectedLines(t *testing.T) {
	profile := machine.Profile{
		Name:         "Summary Machine",
		CPUCore:      "dynamic",
		CPUType:      "pentium",
		Cycles:       "50000",
		Machine:      "svga_s3",
		MemoryMB:     32,
		SoundBlaster: "sb16",
		GUS:          true,
		JoystickType: "auto",
		XMS:          true,
		EMS:          false,
		UMB:          true,
	}

	summary := profileSummary(profile)

	for _, want := range []string{
		"CPU: dynamic / pentium",
		"Cycles: 50000",
		"Video: svga_s3",
		"Memory: 32 MB",
		"Sound Blaster: sb16",
		"Joystick: auto",
		"GUS: on",
		"XMS / EMS / UMB: yes / no / yes",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("profileSummary() missing %q:\n%s", want, summary)
		}
	}
}

func TestProfileSummaryIncludesHardDiskDetails(t *testing.T) {
	profile := machine.Profile{
		Name:          "Disk Machine",
		CPUCore:       "auto",
		CPUType:       "486",
		Cycles:        "25000",
		Machine:       "svga_s3",
		MemoryMB:      16,
		SoundBlaster:  "sb16",
		HardDiskImage: "/tmp/dos.img",
		HardDiskCHS:   "512,63,16,142",
		JoystickType:  "auto",
		XMS:           true,
		EMS:           true,
		UMB:           true,
	}

	summary := profileSummary(profile)

	for _, want := range []string{
		"HDD Image: /tmp/dos.img",
		"HDD CHS: 512,63,16,142",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("profileSummary() missing %q:\n%s", want, summary)
		}
	}
}
