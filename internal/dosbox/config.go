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
		"core":    profile.CPUCore,
		"cputype": profile.CPUType,
		"cycles":  profile.Cycles,
	})
	writeSection(&builder, "mixer", map[string]string{
		"rate":      "44100",
		"blocksize": "1024",
		"prebuffer": "20",
	})
	writeSection(&builder, "sblaster", map[string]string{
		"sbtype": profile.SoundBlaster,
		"sbbase": "220",
		"irq":    "7",
		"dma":    "1",
		"hdma":   "5",
		"mixer":  "true",
	})
	writeSection(&builder, "gus", map[string]string{
		"gus": boolValue(profile.GUS),
	})
	writeSection(&builder, "joystick", map[string]string{
		"joysticktype": joystickValue(profile.Joystick),
	})
	writeSection(&builder, "dos", map[string]string{
		"xms": boolValue(profile.XMS),
		"ems": boolValue(profile.EMS),
		"umb": boolValue(profile.UMB),
	})

	return strings.TrimSpace(builder.String()) + "\n"
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

func boolValue(enabled bool) string {
	if enabled {
		return "true"
	}
	return "false"
}

func joystickValue(enabled bool) string {
	if enabled {
		return "auto"
	}
	return "none"
}
