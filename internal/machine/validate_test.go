package machine

import "testing"

func TestValidateProfile(t *testing.T) {
	valid := Profile{
		Name:         "Valid Machine",
		Description:  "Test profile",
		CPUCore:      "auto",
		CPUType:      "486",
		Cycles:       "25000",
		Machine:      "svga_s3",
		MemoryMB:     16,
		SoundBlaster: "sb16",
		JoystickType: "auto",
		GUS:          false,
		XMS:          true,
		EMS:          true,
		UMB:          true,
	}

	tests := []struct {
		name    string
		profile Profile
		wantErr string
	}{
		{
			name:    "valid",
			profile: valid,
		},
		{
			name:    "missing name",
			profile: withProfile(valid, func(p *Profile) { p.Name = "   " }),
			wantErr: "profile name is required",
		},
		{
			name:    "missing cpu settings",
			profile: withProfile(valid, func(p *Profile) { p.CPUCore = "" }),
			wantErr: "cpu settings are required",
		},
		{
			name:    "missing machine",
			profile: withProfile(valid, func(p *Profile) { p.Machine = "" }),
			wantErr: "video machine is required",
		},
		{
			name:    "missing sound blaster",
			profile: withProfile(valid, func(p *Profile) { p.SoundBlaster = "" }),
			wantErr: "sound blaster type is required",
		},
		{
			name:    "missing joystick type",
			profile: withProfile(valid, func(p *Profile) { p.JoystickType = "" }),
			wantErr: "joystick type is required",
		},
		{
			name: "hard disk image without chs",
			profile: withProfile(valid, func(p *Profile) {
				p.HardDiskImage = "/tmp/disk.img"
			}),
			wantErr: "hard disk image and CHS geometry must both be set",
		},
		{
			name:    "non-positive memory",
			profile: withProfile(valid, func(p *Profile) { p.MemoryMB = 0 }),
			wantErr: "memory must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProfile(tt.profile)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("ValidateProfile() error = %v, want nil", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("ValidateProfile() error = nil, want %q", tt.wantErr)
			}
			if got := err.Error(); got != tt.wantErr {
				t.Fatalf("ValidateProfile() error = %q, want %q", got, tt.wantErr)
			}
		})
	}
}

func withProfile(profile Profile, mutate func(*Profile)) Profile {
	mutated := profile
	mutate(&mutated)
	return mutated
}
