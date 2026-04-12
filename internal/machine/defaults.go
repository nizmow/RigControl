package machine

func NewProfile() Profile {
	return Profile{
		Name:              "New Machine",
		Description:       "",
		CPUCore:           "auto",
		CPUType:           "486",
		Cycles:            "25000",
		Machine:           "svga_s3",
		MemoryMB:          16,
		SoundBlaster:      "sb16",
		MouseCapture:      "onclick",
		MouseRawInput:     true,
		DOSMouseImmediate: false,
		FloppyDiskImages:  nil,
		HardDiskImage:     "",
		HardDiskCHS:       "",
		JoystickType:      "auto",
		GUS:               false,
		XMS:               true,
		EMS:               true,
		UMB:               true,
	}
}
