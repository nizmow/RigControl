package dosbox

import (
	"fmt"
	"slices"
	"strings"

	"rigcontrol/internal/machine"
)

func Render(profile machine.Profile) string {
	var builder strings.Builder

	writeSection(&builder, "dosbox", map[string]string{
		"machine": profile.Machine,
		"memsize": fmt.Sprintf("%d", profile.MemoryMB),
	})
	writeSection(&builder, "cpu", map[string]string{
		"core":       profile.CPUCore,
		"cputype":    profile.CPUType,
		"cpu_cycles": renderCycles(profile),
	})
	writeSection(&builder, "mixer", map[string]string{
		"rate":      "44100",
		"blocksize": "1024",
		"prebuffer": "20",
	})
	writeSection(&builder, "sblaster", map[string]string{
		"sbtype":  profile.SoundBlaster,
		"sbbase":  "220",
		"irq":     "7",
		"dma":     "1",
		"hdma":    "5",
		"sbmixer": "true",
	})
	writeSection(&builder, "gus", map[string]string{
		"gus": boolValue(profile.GUS),
	})
	writeSection(&builder, "mouse", map[string]string{
		"dos_mouse_immediate": boolValue(profile.DOSMouseImmediate),
		"mouse_capture":       profile.MouseCapture,
		"mouse_raw_input":     boolValue(profile.MouseRawInput),
	})
	writeSection(&builder, "joystick", map[string]string{
		"joysticktype": profile.JoystickType,
	})
	writeSection(&builder, "dos", map[string]string{
		"xms": boolValue(profile.XMS),
		"ems": boolValue(profile.EMS),
		"umb": boolValue(profile.UMB),
	})
	if autoexec := autoexecCommands(profile); len(autoexec) > 0 {
		writeListSection(&builder, "autoexec", autoexec)
	}

	return strings.TrimSpace(builder.String()) + "\n"
}

func renderCycles(profile machine.Profile) string {
	cycles := strings.TrimSpace(profile.Cycles)
	if profile.FixedCycles && cycles != "auto" && cycles != "max" && !strings.HasPrefix(cycles, "fixed ") {
		return "fixed " + cycles
	}
	return cycles
}

func writeSection(builder *strings.Builder, name string, values map[string]string) {
	builder.WriteString("[" + name + "]\n")
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.Sort(keys)

	for _, key := range keys {
		value := values[key]
		builder.WriteString(fmt.Sprintf("%s=%s\n", key, value))
	}
	builder.WriteString("\n")
}

func writeListSection(builder *strings.Builder, name string, values []string) {
	builder.WriteString("[" + name + "]\n")
	for _, value := range values {
		builder.WriteString(value + "\n")
	}
	builder.WriteString("\n")
}

func autoexecCommands(profile machine.Profile) []string {
	commands := make([]string, 0, 2)
	if len(profile.FloppyDiskImages) > 0 {
		var builder strings.Builder
		builder.WriteString("imgmount a")
		for _, path := range profile.FloppyDiskImages {
			builder.WriteString(fmt.Sprintf(` "%s"`, path))
		}
		builder.WriteString(" -t floppy")
		commands = append(commands, builder.String())
	}
	if profile.HardDiskImage != "" {
		commands = append(commands,
			fmt.Sprintf(`imgmount 2 "%s" -t hdd -fs none -size %s`, profile.HardDiskImage, profile.HardDiskCHS),
			"boot -l c",
		)
	}
	return commands
}

func boolValue(enabled bool) string {
	if enabled {
		return "true"
	}
	return "false"
}
