package ui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"rigcontrol/internal/dosbox"
	"rigcontrol/internal/machine"
)

const dosboxExecutable = "/Applications/DOSBox Staging.app/Contents/MacOS/dosbox"

func Run() error {
	profiles := machine.Presets()
	if len(profiles) == 0 {
		return errors.New("no machine presets available")
	}

	fyneApp := app.NewWithID("com.nizmow.rigcontrol")
	window := fyneApp.NewWindow("RigControl")
	window.Resize(fyne.NewSize(980, 640))

	selected := profiles[0]

	nameEntry := widget.NewEntry()
	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.Wrapping = fyne.TextWrapWord
	cpuCoreSelect := widget.NewSelect(machine.CPUCoreOptions, nil)
	cpuTypeSelect := widget.NewSelect(machine.CPUTypeOptions, nil)
	cyclesEntry := widget.NewEntry()
	machineSelect := widget.NewSelect(machine.MachineOptions, nil)
	memoryEntry := widget.NewEntry()
	sbSelect := widget.NewSelect(machine.SoundBlasterOptions, nil)
	gusCheck := widget.NewCheck("Gravis Ultrasound", nil)
	joystickSelect := widget.NewSelect(machine.JoystickTypeOptions, nil)
	xmsCheck := widget.NewCheck("XMS", nil)
	emsCheck := widget.NewCheck("EMS", nil)
	umbCheck := widget.NewCheck("UMB", nil)

	preview := widget.NewMultiLineEntry()
	preview.Wrapping = fyne.TextWrapOff
	preview.Disable()

	applyProfile := func(profile machine.Profile) {
		selected = profile
		nameEntry.SetText(profile.Name)
		descriptionEntry.SetText(profile.Description)
		cpuCoreSelect.SetSelected(profile.CPUCore)
		cpuTypeSelect.SetSelected(profile.CPUType)
		cyclesEntry.SetText(profile.Cycles)
		machineSelect.SetSelected(profile.Machine)
		memoryEntry.SetText(fmt.Sprintf("%d", profile.MemoryMB))
		sbSelect.SetSelected(profile.SoundBlaster)
		gusCheck.SetChecked(profile.GUS)
		joystickSelect.SetSelected(profile.JoystickType)
		xmsCheck.SetChecked(profile.XMS)
		emsCheck.SetChecked(profile.EMS)
		umbCheck.SetChecked(profile.UMB)
		preview.SetText(dosbox.Render(profile))
	}

	buildProfile := func() (machine.Profile, error) {
		profile := selected
		profile.Name = strings.TrimSpace(nameEntry.Text)
		profile.Description = strings.TrimSpace(descriptionEntry.Text)
		profile.CPUCore = strings.TrimSpace(cpuCoreSelect.Selected)
		profile.CPUType = strings.TrimSpace(cpuTypeSelect.Selected)
		profile.Cycles = strings.TrimSpace(cyclesEntry.Text)
		profile.Machine = strings.TrimSpace(machineSelect.Selected)
		profile.SoundBlaster = strings.TrimSpace(sbSelect.Selected)
		profile.GUS = gusCheck.Checked
		profile.JoystickType = strings.TrimSpace(joystickSelect.Selected)
		profile.XMS = xmsCheck.Checked
		profile.EMS = emsCheck.Checked
		profile.UMB = umbCheck.Checked

		if profile.Name == "" {
			return machine.Profile{}, errors.New("profile name is required")
		}
		if profile.CPUCore == "" || profile.CPUType == "" || profile.Cycles == "" {
			return machine.Profile{}, errors.New("cpu settings are required")
		}
		if profile.Machine == "" {
			return machine.Profile{}, errors.New("video machine is required")
		}
		if profile.SoundBlaster == "" {
			return machine.Profile{}, errors.New("sound blaster type is required")
		}
		if profile.JoystickType == "" {
			return machine.Profile{}, errors.New("joystick type is required")
		}

		var memoryMB int
		if _, err := fmt.Sscanf(strings.TrimSpace(memoryEntry.Text), "%d", &memoryMB); err != nil || memoryMB <= 0 {
			return machine.Profile{}, errors.New("memory must be a positive integer")
		}
		profile.MemoryMB = memoryMB

		return profile, nil
	}

	refreshPreview := func() {
		profile, err := buildProfile()
		if err != nil {
			preview.SetText("# " + err.Error())
			return
		}
		selected = profile
		preview.SetText(dosbox.Render(profile))
	}

	presetNames := make([]string, 0, len(profiles))
	for _, profile := range profiles {
		presetNames = append(presetNames, profile.Name)
	}

	presetSelect := widget.NewSelect(presetNames, func(name string) {
		profile, ok := machine.ByName(name)
		if !ok {
			return
		}
		applyProfile(profile)
	})

	for _, field := range []*widget.Entry{nameEntry, descriptionEntry, cyclesEntry, memoryEntry} {
		field.OnChanged = func(string) {
			refreshPreview()
		}
	}
	for _, selectWidget := range []*widget.Select{cpuCoreSelect, cpuTypeSelect, machineSelect, sbSelect, joystickSelect} {
		selectWidget.OnChanged = func(string) {
			refreshPreview()
		}
	}
	for _, check := range []*widget.Check{gusCheck, xmsCheck, emsCheck, umbCheck} {
		check.OnChanged = func(bool) {
			refreshPreview()
		}
	}

	launchButton := widget.NewButton("Launch DOSBox", func() {
		profile, err := buildProfile()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		configPath, err := writeTempConfig(profile)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		if err := launchDOSBox(configPath); err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowInformation("DOSBox Started", "Launched DOSBox Staging with "+filepath.Base(configPath), window)
	})

	form := widget.NewForm(
		widget.NewFormItem("Preset", presetSelect),
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Description", descriptionEntry),
		widget.NewFormItem("CPU Core", cpuCoreSelect),
		widget.NewFormItem("CPU Type", cpuTypeSelect),
		widget.NewFormItem("Cycles", cyclesEntry),
		widget.NewFormItem("Video", machineSelect),
		widget.NewFormItem("Memory (MB)", memoryEntry),
		widget.NewFormItem("Sound", sbSelect),
		widget.NewFormItem("Joystick", joystickSelect),
	)

	formPanel := container.NewBorder(
		nil,
		container.NewHBox(launchButton),
		nil,
		nil,
		container.NewVBox(
			form,
			gusCheck,
			xmsCheck,
			emsCheck,
			umbCheck,
		),
	)

	previewPanel := container.NewBorder(
		widget.NewLabel("DOSBox Staging config preview"),
		nil,
		nil,
		nil,
		container.NewStack(preview),
	)

	window.SetContent(container.NewHSplit(formPanel, previewPanel))
	applyProfile(selected)
	presetSelect.SetSelected(selected.Name)
	window.ShowAndRun()

	return nil
}

func writeTempConfig(profile machine.Profile) (string, error) {
	slug := strings.ToLower(profile.Name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	if slug == "" {
		slug = "rigcontrol"
	}

	file, err := os.CreateTemp("", slug+"-*.conf")
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := file.WriteString(dosbox.Render(profile)); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func launchDOSBox(configPath string) error {
	if _, err := os.Stat(dosboxExecutable); err != nil {
		return fmt.Errorf("dosbox executable not found at %s: %w", dosboxExecutable, err)
	}

	cmd := exec.Command(dosboxExecutable, "-conf", configPath)
	return cmd.Start()
}
