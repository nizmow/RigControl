package machine

import (
	"fmt"
	"strings"
)

var FriendlyNames = map[string]string{
	// CPU Cores
	"auto":    "Automatic",
	"dynamic": "Dynamic (JIT)",
	"normal":  "Normal (Interpreter)",
	"simple":  "Simple",

	// CPU Types
	"386":          "Intel 386",
	"386_fast":     "Intel 386 (Fast)",
	"386_prefetch": "Intel 386 (with Prefetch)",
	"486":          "Intel 486",
	"pentium":      "Intel Pentium",
	"pentium_mmx":  "Intel Pentium MMX",

	// Video Machines
	"hercules":      "Hercules (Monochrome)",
	"cga_mono":      "CGA (Monochrome)",
	"cga":           "CGA (Color)",
	"pcjr":          "IBM PCjr",
	"tandy":         "Tandy 1000",
	"ega":           "EGA (Enhanced Graphics)",
	"svga_s3":       "S3 Trio64 (SVGA)",
	"svga_et3000":   "Tseng ET3000 (SVGA)",
	"svga_et4000":   "Tseng ET4000 (SVGA)",
	"svga_paradise": "Paradise PVGA1A (SVGA)",
	"vesa_nolfb":    "VESA (No LFB)",
	"vesa_oldvbe":   "VESA (Old VBE)",

	// Sound Blaster
	"gb":     "Game Blaster (CMS)",
	"sb1":    "Sound Blaster 1.0",
	"sb2":    "Sound Blaster 2.0",
	"sbpro1": "Sound Blaster Pro 1",
	"sbpro2": "Sound Blaster Pro 2",
	"sb16":   "Sound Blaster 16",
	"ess":    "ESS AudioDrive",
	"none":   "Disabled",

	// Mouse Capture
	"seamless": "Seamless (Host integration)",
	"onclick":  "On Click",
	"onstart":  "On Start",
	"nomouse":  "No Mouse",

	// Joystick
	"2axis":    "2-Axis Joystick",
	"4axis":    "4-Axis Joystick",
	"4axis_2":  "Dual 4-Axis Joysticks",
	"fcs":      "Thrustmaster FCS",
	"ch":       "CH Flightstick",
	"hidden":   "Hidden",
	"disabled": "Disabled",
}

func DisplayLabel(key string) string {
	friendly, ok := FriendlyNames[key]
	if !ok || key == "auto" {
		if key == "auto" {
			return "Automatic (auto)"
		}
		return key
	}
	return fmt.Sprintf("%s (%s)", friendly, key)
}

func InternalKey(label string) string {
	// Label format: "Friendly Name (key)"
	if !strings.Contains(label, "(") || !strings.HasSuffix(label, ")") {
		return label
	}
	start := strings.LastIndex(label, "(")
	return label[start+1 : len(label)-1]
}

func DisplayLabels(keys []string) []string {
	labels := make([]string, len(keys))
	for i, k := range keys {
		labels[i] = DisplayLabel(k)
	}
	return labels
}
