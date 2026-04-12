package dosbox

import (
	"strings"
	"testing"

	"rigcontrol/internal/machine"
)

func TestRenderIncludesExpectedSectionsAndValues(t *testing.T) {
	profile := machine.Profile{
		Name:          "Test Machine",
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
	}

	rendered := Render(profile)

	for _, want := range []string{
		"[dosbox]",
		"machine=svga_s3",
		"memsize=32",
		"[cpu]",
		"core=dynamic",
		"cputype=pentium",
		"cpu_cycles=50000",
		"[sblaster]",
		"sbtype=sb16",
		"[gus]",
		"gus=true",
		"[mouse]",
		"mouse_capture=onclick",
		"mouse_raw_input=true",
		"dos_mouse_immediate=false",
		"[joystick]",
		"joysticktype=auto",
		"[dos]",
		"xms=true",
		"ems=true",
		"umb=false",
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered config missing %q:\n%s", want, rendered)
		}
	}
}

func TestRenderIncludesDisabledAndFalseValues(t *testing.T) {
	profile := machine.Profile{
		Name:              "Minimal Machine",
		CPUCore:           "normal",
		CPUType:           "386",
		Cycles:            "1200",
		Machine:           "ega",
		MemoryMB:          2,
		SoundBlaster:      "sb1",
		MouseCapture:      "seamless",
		MouseRawInput:     false,
		DOSMouseImmediate: true,
		GUS:               false,
		JoystickType:      "disabled",
		XMS:               false,
		EMS:               true,
		UMB:               false,
	}

	rendered := Render(profile)

	for _, want := range []string{
		"core=normal",
		"cputype=386",
		"cpu_cycles=1200",
		"gus=false",
		"mouse_capture=seamless",
		"mouse_raw_input=false",
		"dos_mouse_immediate=true",
		"joysticktype=disabled",
		"xms=false",
		"ems=true",
		"umb=false",
		"sbmixer=true",
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered config missing %q:\n%s", want, rendered)
		}
	}
}

func TestRenderIncludesHardDiskAutoexecCommands(t *testing.T) {
	profile := machine.Profile{
		Name:          "Boot HDD",
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

	rendered := Render(profile)

	for _, want := range []string{
		"[autoexec]",
		`imgmount 2 "/tmp/dos.img" -t hdd -fs none -size 512,63,16,142`,
		"boot -l c",
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered config missing %q:\n%s", want, rendered)
		}
	}
}

func TestRenderIncludesFloppySwapMountCommand(t *testing.T) {
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

	rendered := Render(profile)

	for _, want := range []string{
		"[autoexec]",
		`imgmount a "/tmp/disk1.img" "/tmp/disk2.img" -t floppy`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered config missing %q:\n%s", want, rendered)
		}
	}
}
