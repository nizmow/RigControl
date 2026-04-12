package machine

type Profile struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	CPUCore      string `json:"cpu_core"`
	CPUType      string `json:"cpu_type"`
	Cycles       string `json:"cycles"`
	Machine      string `json:"machine"`
	MemoryMB     int    `json:"memory_mb"`
	SoundBlaster string `json:"sound_blaster"`
	GUS          bool   `json:"gus"`
	JoystickType string `json:"joystick_type"`
	XMS          bool   `json:"xms"`
	EMS          bool   `json:"ems"`
	UMB          bool   `json:"umb"`
}
