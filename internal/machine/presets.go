package machine

func Presets() []Profile {
	return []Profile{
		{
			Name:         "286 EGA",
			Description:  "Low-end late 80s setup for EGA adventure and strategy titles.",
			CPUCore:      "normal",
			CPUType:      "286",
			Cycles:       "fixed 1200",
			Machine:      "ega",
			MemoryMB:     2,
			SoundBlaster: "sb1",
			GUS:          false,
			Joystick:     false,
			XMS:          false,
			EMS:          true,
			UMB:          true,
		},
		{
			Name:         "486 SVGA",
			Description:  "Balanced early 90s machine for VGA DOS games with Sound Blaster 16.",
			CPUCore:      "auto",
			CPUType:      "486_slow",
			Cycles:       "max 70%",
			Machine:      "svga_s3",
			MemoryMB:     16,
			SoundBlaster: "sb16",
			GUS:          false,
			Joystick:     true,
			XMS:          true,
			EMS:          true,
			UMB:          true,
		},
		{
			Name:         "Fast Pentium",
			Description:  "High-end DOS setup for demanding SVGA and Build engine games.",
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
			UMB:          true,
		},
	}
}

func ByName(name string) (Profile, bool) {
	for _, preset := range Presets() {
		if preset.Name == name {
			return preset, true
		}
	}

	return Profile{}, false
}
