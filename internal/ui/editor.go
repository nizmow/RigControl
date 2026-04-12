package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"rigcontrol/internal/machine"
)

type profileEditor struct {
	window     fyne.Window
	form       *profileForm
	errorLabel *widget.Label
	onSave     func(machine.Profile) error
}

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
	editor := newProfileEditor(profile, onSave)
	editor.show()
}

func newProfileEditor(profile machine.Profile, onSave func(machine.Profile) error) *profileEditor {
	errorLabel := widget.NewLabel("")
	errorLabel.Wrapping = fyne.TextWrapWord

	editor := &profileEditor{
		window:     fyne.CurrentApp().NewWindow("Edit Machine"),
		form:       newProfileForm(profile),
		errorLabel: errorLabel,
		onSave:     onSave,
	}

	editor.window.Resize(fyne.NewSize(620, 760))
	editor.window.SetContent(container.NewPadded(editor.content()))
	editor.window.CenterOnScreen()

	return editor
}

func (e *profileEditor) content() fyne.CanvasObject {
	saveButton := widget.NewButton("Save", e.handleSave)
	cancelButton := widget.NewButton("Cancel", func() {
		e.window.Close()
	})

	return container.NewBorder(
		nil,
		container.NewHBox(cancelButton, saveButton),
		nil,
		nil,
		container.NewVBox(
			e.errorLabel,
			e.form.widget(),
			e.form.gus,
			e.form.xms,
			e.form.ems,
			e.form.umb,
		),
	)
}

func (e *profileEditor) handleSave() {
	profile, err := e.form.profile()
	if err != nil {
		e.errorLabel.SetText(err.Error())
		return
	}
	if err := e.onSave(profile); err != nil {
		e.errorLabel.SetText(err.Error())
		return
	}

	e.window.Close()
}

func (e *profileEditor) show() {
	e.window.Show()
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

	form.setProfile(profile)
	return form
}

func (f *profileForm) widget() fyne.CanvasObject {
	return widget.NewForm(
		widget.NewFormItem("Name", f.name),
		widget.NewFormItem("Description", f.description),
		widget.NewFormItem("CPU Core", f.cpuCore),
		widget.NewFormItem("CPU Type", f.cpuType),
		widget.NewFormItem("Cycles", f.cycles),
		widget.NewFormItem("Video", f.videoMachine),
		widget.NewFormItem("Memory (MB)", f.memoryMB),
		widget.NewFormItem("Sound", f.soundBlaster),
		widget.NewFormItem("Joystick", f.joystickType),
	)
}

func (f *profileForm) setProfile(profile machine.Profile) {
	f.name.SetText(profile.Name)
	f.description.SetText(profile.Description)
	f.cpuCore.SetSelected(profile.CPUCore)
	f.cpuType.SetSelected(profile.CPUType)
	f.cycles.SetText(profile.Cycles)
	f.videoMachine.SetSelected(profile.Machine)
	f.memoryMB.SetText(fmt.Sprintf("%d", profile.MemoryMB))
	f.soundBlaster.SetSelected(profile.SoundBlaster)
	f.joystickType.SetSelected(profile.JoystickType)
	f.gus.SetChecked(profile.GUS)
	f.xms.SetChecked(profile.XMS)
	f.ems.SetChecked(profile.EMS)
	f.umb.SetChecked(profile.UMB)
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
