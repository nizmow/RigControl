package ui

import (
	"strings"
	"testing"

	"rigcontrol/internal/machine"
)

func TestProfileSummaryIncludesExpectedLines(t *testing.T) {
	profile := machine.Profile{
		Name:          "Summary Machine",
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
		EMS:           false,
		UMB:           true,
	}

	summary := profileSummary(profile)

	for _, want := range []string{
		"CPU: dynamic / pentium",
		"Cycles: 50000",
		"Video: svga_s3",
		"Memory: 32 MB",
		"Sound Blaster: sb16",
		"Mouse Capture: onclick",
		"Mouse Raw Input: yes",
		"Mouse Immediate: no",
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
		MouseCapture:  "onclick",
		MouseRawInput: true,
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

func TestProfileSummaryIncludesFloppyDetails(t *testing.T) {
	profile := machine.Profile{
		Name:             "Floppy Machine",
		CPUCore:          "auto",
		CPUType:          "486",
		Cycles:           "25000",
		Machine:          "svga_s3",
		MemoryMB:         16,
		SoundBlaster:     "sb16",
		MouseCapture:     "onclick",
		MouseRawInput:    true,
		FloppyDiskImages: []string{"/tmp/disk1.img", "/tmp/disk2.img"},
		JoystickType:     "auto",
		XMS:              true,
		EMS:              true,
		UMB:              true,
	}

	summary := profileSummary(profile)

	for _, want := range []string{
		"Floppy 1: /tmp/disk1.img",
		"Floppy 2: /tmp/disk2.img",
	} {
		if !strings.Contains(summary, want) {
			t.Fatalf("profileSummary() missing %q:\n%s", want, summary)
		}
	}
}
