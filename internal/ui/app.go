package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"rigcontrol/internal/dosbox"
	"rigcontrol/internal/machine"
)

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
	cpuCoreEntry := widget.NewEntry()
	cpuTypeEntry := widget.NewEntry()
	cyclesEntry := widget.NewEntry()
	machineEntry := widget.NewEntry()
	memoryEntry := widget.NewEntry()
	sbEntry := widget.NewEntry()
	gusCheck := widget.NewCheck("Gravis Ultrasound", nil)
	joystickCheck := widget.NewCheck("Joystick enabled", nil)
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
		cpuCoreEntry.SetText(profile.CPUCore)
		cpuTypeEntry.SetText(profile.CPUType)
		cyclesEntry.SetText(profile.Cycles)
		machineEntry.SetText(profile.Machine)
		memoryEntry.SetText(fmt.Sprintf("%d", profile.MemoryMB))
		sbEntry.SetText(profile.SoundBlaster)
		gusCheck.SetChecked(profile.GUS)
		joystickCheck.SetChecked(profile.Joystick)
		xmsCheck.SetChecked(profile.XMS)
		emsCheck.SetChecked(profile.EMS)
		umbCheck.SetChecked(profile.UMB)
		preview.SetText(dosbox.Render(profile))
	}

	buildProfile := func() (machine.Profile, error) {
		profile := selected
		profile.Name = strings.TrimSpace(nameEntry.Text)
		profile.Description = strings.TrimSpace(descriptionEntry.Text)
		profile.CPUCore = strings.TrimSpace(cpuCoreEntry.Text)
		profile.CPUType = strings.TrimSpace(cpuTypeEntry.Text)
		profile.Cycles = strings.TrimSpace(cyclesEntry.Text)
		profile.Machine = strings.TrimSpace(machineEntry.Text)
		profile.SoundBlaster = strings.TrimSpace(sbEntry.Text)
		profile.GUS = gusCheck.Checked
		profile.Joystick = joystickCheck.Checked
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

	for _, field := range []*widget.Entry{nameEntry, descriptionEntry, cpuCoreEntry, cpuTypeEntry, cyclesEntry, machineEntry, memoryEntry, sbEntry} {
		field.OnChanged = func(string) {
			refreshPreview()
		}
	}
	for _, check := range []*widget.Check{gusCheck, joystickCheck, xmsCheck, emsCheck, umbCheck} {
		check.OnChanged = func(bool) {
			refreshPreview()
		}
	}

	saveJSON := widget.NewButton("Save Profile JSON", func() {
		profile, err := buildProfile()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()

			payload, err := json.MarshalIndent(profile, "", "  ")
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if _, err := writer.Write(append(payload, '\n')); err != nil {
				dialog.ShowError(err, window)
			}
		}, window)
	})

	saveConfig := widget.NewButton("Export DOSBox Config", func() {
		profile, err := buildProfile()
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if writer == nil {
				return
			}
			defer writer.Close()

			if _, err := io.WriteString(writer, dosbox.Render(profile)); err != nil {
				dialog.ShowError(err, window)
			}
		}, window)
	})

	form := widget.NewForm(
		widget.NewFormItem("Preset", presetSelect),
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Description", descriptionEntry),
		widget.NewFormItem("CPU Core", cpuCoreEntry),
		widget.NewFormItem("CPU Type", cpuTypeEntry),
		widget.NewFormItem("Cycles", cyclesEntry),
		widget.NewFormItem("Video", machineEntry),
		widget.NewFormItem("Memory (MB)", memoryEntry),
		widget.NewFormItem("Sound", sbEntry),
	)

	formPanel := container.NewBorder(
		nil,
		container.NewHBox(saveJSON, saveConfig),
		nil,
		nil,
		container.NewVBox(
			form,
			gusCheck,
			joystickCheck,
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
