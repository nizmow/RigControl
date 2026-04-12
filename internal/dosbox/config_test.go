package dosbox

import (
	"strings"
	"testing"

	"rigcontrol/internal/machine"
)

func TestRenderIncludesExpectedSectionsAndValues(t *testing.T) {
	profile := machine.Profile{
		Name:         "Test Machine",
		CPUCore:      "dynamic",
		CPUType:      "pentium_slow",
		Cycles:       "max",
		Machine:      "svga_s3",
		MemoryMB:     32,
		SoundBlaster: "sb16",
		GUS:          true,
		Joystick:     true,
		XMS:          true,
		EMS:          true,
		UMB:          false,
	}

	rendered := Render(profile)

	for _, want := range []string{
		"[dosbox]",
		"machine=svga_s3",
		"memsize=32",
		"[cpu]",
		"core=dynamic",
		"cputype=pentium_slow",
		"cycles=max",
		"[sblaster]",
		"sbtype=sb16",
		"[gus]",
		"gus=true",
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
