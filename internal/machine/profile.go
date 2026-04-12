package machine

type Profile struct {
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	CPUCore           string   `json:"cpu_core"`
	CPUType           string   `json:"cpu_type"`
	Cycles            string   `json:"cycles"`
	FixedCycles       bool     `json:"fixed_cycles"`
	Machine           string   `json:"machine"`
	MemoryMB          int      `json:"memory_mb"`
	SoundBlaster      string   `json:"sound_blaster"`
	MouseCapture      string   `json:"mouse_capture"`
	MouseRawInput     bool     `json:"mouse_raw_input"`
	DOSMouseImmediate bool     `json:"dos_mouse_immediate"`
	FloppyDiskImages  []string `json:"floppy_disk_images,omitempty"`
	HardDiskImage     string   `json:"hard_disk_image,omitempty"`
	HardDiskCHS       string   `json:"hard_disk_chs,omitempty"`
	GUS               bool     `json:"gus"`
	JoystickType      string   `json:"joystick_type"`
	XMS               bool     `json:"xms"`
	EMS               bool     `json:"ems"`
	UMB               bool     `json:"umb"`
}
