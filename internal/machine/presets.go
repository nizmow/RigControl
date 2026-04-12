package machine

func Presets() []Profile {
	return []Profile{
		{
			Name:         "286 EGA",
			Description:  "Low-end late 80s setup for EGA adventure and strategy titles.",
			CPUCore:      "normal",
			CPUType:      "386",
			Cycles:       "1200",
			Machine:      "ega",
			MemoryMB:     2,
			SoundBlaster: "sb1",
			GUS:          false,
			JoystickType: "disabled",
			XMS:          false,
			EMS:          true,
			UMB:          true,
		},
		{
			Name:         "486 SVGA",
			Description:  "Balanced early 90s machine for VGA DOS games with Sound Blaster 16.",
			CPUCore:      "auto",
			CPUType:      "486",
			Cycles:       "25000",
			Machine:      "svga_s3",
			MemoryMB:     16,
			SoundBlaster: "sb16",
			GUS:          false,
			JoystickType: "auto",
			XMS:          true,
			EMS:          true,
			UMB:          true,
		},
		{
			Name:         "Fast Pentium",
			Description:  "High-end DOS setup for demanding SVGA and Build engine games.",
			CPUCore:      "dynamic",
			CPUType:      "pentium",
			Cycles:       "50000",
			Machine:      "svga_s3",
			MemoryMB:     32,
			SoundBlaster: "sb16",
			GUS:          true,
			JoystickType: "auto",
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
