package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"rigcontrol/internal/machine"
)

type profileForm struct {
	name         *widget.Entry
	description  *widget.Entry
	cpuCore      *widget.Select
	cpuType      *widget.Select
	cycles       *widget.Entry
	videoMachine *widget.Select
	memoryMB     *widget.Entry
	soundBlaster *widget.Select
	joystickType *widget.Select
	gus          *widget.Check
	xms          *widget.Check
	ems          *widget.Check
	umb          *widget.Check
}

func showProfileEditor(window fyne.Window, profile machine.Profile, onSave func(machine.Profile) error) {
	form := newProfileForm(profile)
	errorLabel := widget.NewLabel("")
	errorLabel.Wrapping = fyne.TextWrapWord

	saveButton := widget.NewButton("Save", nil)
	cancelButton := widget.NewButton("Cancel", nil)

	editorWindow := fyne.CurrentApp().NewWindow("Edit Machine")
	editorWindow.Resize(fyne.NewSize(620, 760))

	content := container.NewBorder(
		nil,
		container.NewHBox(cancelButton, saveButton),
		nil,
		nil,
		container.NewVBox(
			errorLabel,
			widget.NewForm(
				widget.NewFormItem("Name", form.name),
				widget.NewFormItem("Description", form.description),
				widget.NewFormItem("CPU Core", form.cpuCore),
				widget.NewFormItem("CPU Type", form.cpuType),
				widget.NewFormItem("Cycles", form.cycles),
				widget.NewFormItem("Video", form.videoMachine),
				widget.NewFormItem("Memory (MB)", form.memoryMB),
				widget.NewFormItem("Sound", form.soundBlaster),
				widget.NewFormItem("Joystick", form.joystickType),
			),
			form.gus,
			form.xms,
			form.ems,
			form.umb,
		),
	)

	editorWindow.SetContent(container.NewPadded(content))
	editorWindow.CenterOnScreen()

	cancelButton.OnTapped = func() {
		editorWindow.Close()
	}
	saveButton.OnTapped = func() {
		updated, err := form.profile()
		if err != nil {
			errorLabel.SetText(err.Error())
			return
		}
		if err := onSave(updated); err != nil {
			errorLabel.SetText(err.Error())
			return
		}
		editorWindow.Close()
	}

	editorWindow.Show()
}

func newProfileForm(profile machine.Profile) *profileForm {
	description := widget.NewMultiLineEntry()
	description.Wrapping = fyne.TextWrapWord

	form := &profileForm{
		name:         widget.NewEntry(),
		description:  description,
		cpuCore:      widget.NewSelect(machine.CPUCoreOptions, nil),
		cpuType:      widget.NewSelect(machine.CPUTypeOptions, nil),
		cycles:       widget.NewEntry(),
		videoMachine: widget.NewSelect(machine.MachineOptions, nil),
		memoryMB:     widget.NewEntry(),
		soundBlaster: widget.NewSelect(machine.SoundBlasterOptions, nil),
		joystickType: widget.NewSelect(machine.JoystickTypeOptions, nil),
		gus:          widget.NewCheck("Gravis Ultrasound", nil),
		xms:          widget.NewCheck("XMS", nil),
		ems:          widget.NewCheck("EMS", nil),
		umb:          widget.NewCheck("UMB", nil),
	}

	form.name.SetText(profile.Name)
	form.description.SetText(profile.Description)
	form.cpuCore.SetSelected(profile.CPUCore)
	form.cpuType.SetSelected(profile.CPUType)
	form.cycles.SetText(profile.Cycles)
	form.videoMachine.SetSelected(profile.Machine)
	form.memoryMB.SetText(fmt.Sprintf("%d", profile.MemoryMB))
	form.soundBlaster.SetSelected(profile.SoundBlaster)
	form.joystickType.SetSelected(profile.JoystickType)
	form.gus.SetChecked(profile.GUS)
	form.xms.SetChecked(profile.XMS)
	form.ems.SetChecked(profile.EMS)
	form.umb.SetChecked(profile.UMB)

	return form
}

func (f *profileForm) profile() (machine.Profile, error) {
	profile := machine.Profile{
		Name:         strings.TrimSpace(f.name.Text),
		Description:  strings.TrimSpace(f.description.Text),
		CPUCore:      strings.TrimSpace(f.cpuCore.Selected),
		CPUType:      strings.TrimSpace(f.cpuType.Selected),
		Cycles:       strings.TrimSpace(f.cycles.Text),
		Machine:      strings.TrimSpace(f.videoMachine.Selected),
		SoundBlaster: strings.TrimSpace(f.soundBlaster.Selected),
		JoystickType: strings.TrimSpace(f.joystickType.Selected),
		GUS:          f.gus.Checked,
		XMS:          f.xms.Checked,
		EMS:          f.ems.Checked,
		UMB:          f.umb.Checked,
	}

	if _, err := fmt.Sscanf(strings.TrimSpace(f.memoryMB.Text), "%d", &profile.MemoryMB); err != nil {
		return machine.Profile{}, fmt.Errorf("memory must be a positive integer")
	}
	if err := machine.ValidateProfile(profile); err != nil {
		return machine.Profile{}, err
	}

	return profile, nil
}
