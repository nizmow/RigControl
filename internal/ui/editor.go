package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	fynetooltip "github.com/dweymouth/fyne-tooltip"
	ttwidget "github.com/dweymouth/fyne-tooltip/widget"

	"rigcontrol/internal/machine"
)

const (
	defaultEditorWindowWidth  = 900
	defaultEditorWindowHeight = 760
	minEditorWindowWidth      = 720
	minEditorWindowHeight     = 560
)

type profileEditor struct {
	window     fyne.Window
	form       *profileForm
	errorLabel *widget.Label
	onSave     func(machine.Profile) error
}

type profileForm struct {
	name              *widget.Entry
	description       *widget.Entry
	cpuCore           *widget.Select
	cpuType           *widget.Select
	cycles            *widget.Entry
	fixedCycles       *widget.Check
	videoMachine      *widget.Select
	memoryMB          *widget.Entry
	soundBlaster      *widget.Select
	mouseCapture      *widget.Select
	mouseRawInput     *widget.Check
	dosMouseImmediate *widget.Check
	floppyDiskImages  []string
	floppyList        *widget.List
	floppySelection   int
	moveFloppyUp      *widget.Button
	moveFloppyDown    *widget.Button
	removeFloppy      *widget.Button
	hardDiskImage     *widget.Entry
	hardDiskCHS       *widget.Entry
	joystickType      *widget.Select
	gus               *widget.Check
	xms               *widget.Check
	ems               *widget.Check
	umb               *widget.Check
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

	editor.window.Resize(initialEditorWindowSize(fyne.CurrentApp()))
	fynetooltip.SetToolTipTextSizeName(theme.SizeNameText)
	content := container.NewPadded(editor.content())
	editor.window.SetContent(fynetooltip.AddWindowToolTipLayer(content, editor.window.Canvas()))
	editor.window.CenterOnScreen()

	return editor
}

func (e *profileEditor) content() fyne.CanvasObject {
	saveButton := widget.NewButton("Save", e.handleSave)
	cancelButton := widget.NewButton("Cancel", func() {
		e.window.Close()
	})

	tabs := container.NewAppTabs(
		container.NewTabItem("General", e.form.generalTab()),
		container.NewTabItem("Hardware", e.form.hardwareTab()),
		container.NewTabItem("Storage", e.form.storageTab(e.window)),
		container.NewTabItem("Audio & Input", e.form.audioInputTab()),
	)

	return container.NewBorder(
		e.errorLabel,
		container.NewHBox(cancelButton, saveButton),
		nil,
		nil,
		tabs,
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
		name:              widget.NewEntry(),
		description:       description,
		cpuCore:           widget.NewSelect(machine.CPUCoreOptions, nil),
		cpuType:           widget.NewSelect(machine.CPUTypeOptions, nil),
		cycles:            widget.NewEntry(),
		fixedCycles:       widget.NewCheck("Fixed", nil),
		videoMachine:      widget.NewSelect(machine.MachineOptions, nil),
		memoryMB:          widget.NewEntry(),
		soundBlaster:      widget.NewSelect(machine.SoundBlasterOptions, nil),
		mouseCapture:      widget.NewSelect(machine.MouseCaptureOptions, nil),
		mouseRawInput:     widget.NewCheck("Use raw host mouse input", nil),
		dosMouseImmediate: widget.NewCheck("Update DOS mouse immediately", nil),
		floppyDiskImages:  nil,
		floppySelection:   -1,
		hardDiskImage:     widget.NewEntry(),
		hardDiskCHS:       widget.NewEntry(),
		joystickType:      widget.NewSelect(machine.JoystickTypeOptions, nil),
		gus:               widget.NewCheck("Gravis Ultrasound", nil),
		xms:               widget.NewCheck("XMS", nil),
		ems:               widget.NewCheck("EMS", nil),
		umb:               widget.NewCheck("UMB", nil),
	}
	form.floppyList = widget.NewList(
		func() int { return len(form.floppyDiskImages) },
		func() fyne.CanvasObject {
			return newFloppyListRow()
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			updateFloppyListRow(obj, id+1, form.floppyDiskImages[id])
		},
	)
	form.floppyList.OnSelected = func(id widget.ListItemID) {
		form.floppySelection = id
		form.refreshFloppyButtons()
	}
	form.floppyList.OnUnselected = func(widget.ListItemID) {
		form.floppySelection = -1
		form.refreshFloppyButtons()
	}
	form.moveFloppyUp = widget.NewButton("Move Up", form.moveSelectedFloppyUp)
	form.moveFloppyDown = widget.NewButton("Move Down", form.moveSelectedFloppyDown)
	form.removeFloppy = widget.NewButton("Remove", form.removeSelectedFloppy)
	form.refreshFloppyButtons()

	form.setProfile(profile)
	return form
}

func (f *profileForm) generalTab() fyne.CanvasObject {
	return container.NewVScroll(container.NewPadded(widget.NewForm(
		widget.NewFormItem("Name", f.name),
		widget.NewFormItem("Description", f.description),
	)))
}

func (f *profileForm) hardwareTab() fyne.CanvasObject {
	return container.NewVScroll(container.NewPadded(widget.NewForm(
		widget.NewFormItem("CPU Core", f.cpuCore),
		widget.NewFormItem("CPU Type", f.cpuType),
		widget.NewFormItem("Cycles", container.NewBorder(nil, nil, nil, f.fixedCycles, f.cycles)),
		widget.NewFormItem("Video", f.videoMachine),
		widget.NewFormItem("Memory (MB)", f.memoryMB),
		widget.NewFormItem("DOS Memory", container.NewVBox(f.xms, f.ems, f.umb)),
	)))
}

func (f *profileForm) storageTab(window fyne.Window) fyne.CanvasObject {
	browseButton := widget.NewButton("Browse...", func() {
		f.pickHardDiskImage(window)
	})
	autoFillCHSButton := widget.NewButton("Auto-fill", func() {
		chs, err := inferCHSFromImage(f.hardDiskImage.Text)
		if err != nil {
			f.hardDiskCHS.SetText("")
			return
		}
		f.hardDiskCHS.SetText(chs)
	})

	return container.NewVScroll(container.NewPadded(widget.NewForm(
		widget.NewFormItem("Floppy Disks", f.buildFloppyEditor(window)),
		widget.NewFormItem("HDD Image", container.NewBorder(nil, nil, nil, browseButton, f.hardDiskImage)),
		widget.NewFormItem("HDD CHS", fieldWithHelp(container.NewBorder(nil, nil, nil, autoFillCHSButton, f.hardDiskCHS), "Cylinder-head-sector geometry used for HDD image booting. Auto-fill works for images sized as 512-byte sectors, 63 sectors/track, 16 heads.")),
	)))
}

func (f *profileForm) audioInputTab() fyne.CanvasObject {
	return container.NewVScroll(container.NewPadded(widget.NewForm(
		widget.NewFormItem("Sound", f.soundBlaster),
		widget.NewFormItem("GUS", f.gus),
		widget.NewFormItem("Joystick", f.joystickType),
		widget.NewFormItem("Mouse Capture", fieldWithHelp(f.mouseCapture, "How DOSBox grabs or releases the mouse in the host window.")),
		widget.NewFormItem("Mouse Options", fieldWithHelp(container.NewVBox(f.mouseRawInput, f.dosMouseImmediate), "Raw input bypasses host acceleration while captured. Immediate updates can help some action games but may over-accelerate others.")),
	)))
}

func (f *profileForm) buildFloppyEditor(window fyne.Window) fyne.CanvasObject {
	addButton := widget.NewButton("Add Disk...", func() {
		f.pickFloppyDisk(window)
	})

	return container.NewBorder(
		nil,
		nil,
		nil,
		container.NewVBox(addButton, f.moveFloppyUp, f.moveFloppyDown, f.removeFloppy),
		f.floppyList,
	)
}

type helpIcon struct {
	ttwidget.ToolTipWidget
	icon *widget.Icon
}

func newHelpIcon(helpText string) *helpIcon {
	h := &helpIcon{
		icon: widget.NewIcon(theme.QuestionIcon()),
	}
	h.SetToolTip(helpText)
	h.ExtendBaseWidget(h)
	return h
}

func (h *helpIcon) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(h.icon)
}

func (h *helpIcon) MinSize() fyne.Size {
	return h.icon.MinSize()
}

func (h *helpIcon) Resize(s fyne.Size) {
	h.BaseWidget.Resize(s)
	h.icon.Resize(s)
}

func fieldWithHelp(field fyne.CanvasObject, help string) fyne.CanvasObject {
	if strings.TrimSpace(help) == "" {
		return field
	}

	// Use Border layout to keep the icon on the right, and wrap it in a VBox
	// so it stays at the top of the container height instead of centering.
	return container.NewBorder(nil, nil, nil, container.NewVBox(newHelpIcon(help)), field)
}

func (f *profileForm) setProfile(profile machine.Profile) {
	f.name.SetText(profile.Name)
	f.description.SetText(profile.Description)
	f.cpuCore.SetSelected(profile.CPUCore)
	f.cpuType.SetSelected(profile.CPUType)
	f.cycles.SetText(profile.Cycles)
	f.fixedCycles.SetChecked(profile.FixedCycles)
	f.videoMachine.SetSelected(profile.Machine)
	f.memoryMB.SetText(fmt.Sprintf("%d", profile.MemoryMB))
	f.soundBlaster.SetSelected(profile.SoundBlaster)
	f.mouseCapture.SetSelected(profile.MouseCapture)
	f.mouseRawInput.SetChecked(profile.MouseRawInput)
	f.dosMouseImmediate.SetChecked(profile.DOSMouseImmediate)
	f.floppyDiskImages = append([]string(nil), profile.FloppyDiskImages...)
	f.floppySelection = -1
	f.floppyList.UnselectAll()
	f.floppyList.Refresh()
	f.refreshFloppyButtons()
	f.hardDiskImage.SetText(profile.HardDiskImage)
	f.hardDiskCHS.SetText(profile.HardDiskCHS)
	f.joystickType.SetSelected(profile.JoystickType)
	f.gus.SetChecked(profile.GUS)
	f.xms.SetChecked(profile.XMS)
	f.ems.SetChecked(profile.EMS)
	f.umb.SetChecked(profile.UMB)
}

func (f *profileForm) profile() (machine.Profile, error) {
	profile := machine.Profile{
		Name:              strings.TrimSpace(f.name.Text),
		Description:       strings.TrimSpace(f.description.Text),
		CPUCore:           strings.TrimSpace(f.cpuCore.Selected),
		CPUType:           strings.TrimSpace(f.cpuType.Selected),
		Cycles:            strings.TrimSpace(f.cycles.Text),
		FixedCycles:       f.fixedCycles.Checked,
		Machine:           strings.TrimSpace(f.videoMachine.Selected),
		SoundBlaster:      strings.TrimSpace(f.soundBlaster.Selected),
		MouseCapture:      strings.TrimSpace(f.mouseCapture.Selected),
		MouseRawInput:     f.mouseRawInput.Checked,
		DOSMouseImmediate: f.dosMouseImmediate.Checked,
		FloppyDiskImages:  append([]string(nil), f.floppyDiskImages...),
		HardDiskImage:     strings.TrimSpace(f.hardDiskImage.Text),
		HardDiskCHS:       strings.TrimSpace(f.hardDiskCHS.Text),
		JoystickType:      strings.TrimSpace(f.joystickType.Selected),
		GUS:               f.gus.Checked,
		XMS:               f.xms.Checked,
		EMS:               f.ems.Checked,
		UMB:               f.umb.Checked,
	}

	if _, err := fmt.Sscanf(strings.TrimSpace(f.memoryMB.Text), "%d", &profile.MemoryMB); err != nil {
		return machine.Profile{}, fmt.Errorf("memory must be a positive integer")
	}
	if err := machine.ValidateProfile(profile); err != nil {
		return machine.Profile{}, err
	}

	return profile, nil
}

func (f *profileForm) pickFloppyDisk(window fyne.Window) {
	open := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		f.addFloppyDiskPath(reader.URI().Path())
	}, window)
	if dir := preferredImageDir(f.selectedFloppyPath(), f.hardDiskImage.Text); dir != nil {
		open.SetLocation(dir)
	}
	open.Show()
}

func (f *profileForm) addFloppyDiskPath(path string) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}
	f.floppyDiskImages = append(f.floppyDiskImages, path)
	f.floppyList.Refresh()
	f.floppyList.Select(len(f.floppyDiskImages) - 1)
	f.refreshFloppyButtons()
}

func (f *profileForm) removeSelectedFloppy() {
	if !f.hasSelectedFloppy() {
		return
	}
	index := f.floppySelection
	f.floppyDiskImages = append(f.floppyDiskImages[:index], f.floppyDiskImages[index+1:]...)
	f.floppySelection = -1
	f.floppyList.Refresh()
	f.floppyList.UnselectAll()
	if len(f.floppyDiskImages) > 0 {
		next := index
		if next >= len(f.floppyDiskImages) {
			next = len(f.floppyDiskImages) - 1
		}
		f.floppyList.Select(next)
	} else {
		f.refreshFloppyButtons()
	}
}

func (f *profileForm) moveSelectedFloppyUp() {
	if !f.hasSelectedFloppy() || f.floppySelection == 0 {
		return
	}
	index := f.floppySelection
	f.floppyDiskImages[index-1], f.floppyDiskImages[index] = f.floppyDiskImages[index], f.floppyDiskImages[index-1]
	f.floppyList.Refresh()
	f.floppyList.Select(index - 1)
}

func (f *profileForm) moveSelectedFloppyDown() {
	if !f.hasSelectedFloppy() || f.floppySelection >= len(f.floppyDiskImages)-1 {
		return
	}
	index := f.floppySelection
	f.floppyDiskImages[index], f.floppyDiskImages[index+1] = f.floppyDiskImages[index+1], f.floppyDiskImages[index]
	f.floppyList.Refresh()
	f.floppyList.Select(index + 1)
}

func (f *profileForm) selectedFloppyPath() string {
	if !f.hasSelectedFloppy() {
		return ""
	}
	return f.floppyDiskImages[f.floppySelection]
}

func (f *profileForm) hasSelectedFloppy() bool {
	return f.floppySelection >= 0 && f.floppySelection < len(f.floppyDiskImages)
}

func (f *profileForm) refreshFloppyButtons() {
	if f.hasSelectedFloppy() {
		f.removeFloppy.Enable()
	} else {
		f.removeFloppy.Disable()
	}
	if f.hasSelectedFloppy() && f.floppySelection > 0 {
		f.moveFloppyUp.Enable()
	} else {
		f.moveFloppyUp.Disable()
	}
	if f.hasSelectedFloppy() && f.floppySelection < len(f.floppyDiskImages)-1 {
		f.moveFloppyDown.Enable()
	} else {
		f.moveFloppyDown.Disable()
	}
}

func (f *profileForm) pickHardDiskImage(window fyne.Window) {
	open := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		f.setHardDiskImagePath(reader.URI().Path())
	}, window)

	if dir := preferredImageDir(f.hardDiskImage.Text, f.selectedFloppyPath()); dir != nil {
		open.SetLocation(dir)
	}

	open.Show()
}

func (f *profileForm) setHardDiskImagePath(path string) {
	f.hardDiskImage.SetText(strings.TrimSpace(path))
	if chs, err := inferCHSFromImage(f.hardDiskImage.Text); err == nil {
		f.hardDiskCHS.SetText(chs)
	}
}

func preferredImageDir(paths ...string) fyne.ListableURI {
	for _, candidate := range paths {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}
		dirURI := storage.NewFileURI(filepath.Dir(candidate))
		if dir, err := storage.ListerForURI(dirURI); err == nil {
			return dir
		}
	}
	return nil
}

func newFloppyListRow() fyne.CanvasObject {
	prefix := widget.NewLabel("1. ")
	prefix.Wrapping = fyne.TextWrapOff
	prefixContainer := container.NewWithoutLayout(prefix)
	prefixContainer.Resize(fyne.NewSize(28, prefix.MinSize().Height))
	return container.NewBorder(nil, nil, prefixContainer, nil, newPathValueLabel("/full/path/to/disk.img"))
}

func updateFloppyListRow(obj fyne.CanvasObject, index int, path string) {
	row := obj.(*fyne.Container)
	prefixContainer, pathValue := findFloppyRowParts(row)
	if prefixContainer == nil || pathValue == nil {
		panic("unexpected floppy row shape")
	}
	prefixContainer.Objects[0].(*widget.Label).SetText(fmt.Sprintf("%d. ", index))
	pathValue.SetText(path)
}

func findFloppyRowParts(row *fyne.Container) (*fyne.Container, *pathValueLabel) {
	var prefixContainer *fyne.Container
	var pathValue *pathValueLabel
	for _, obj := range row.Objects {
		switch typed := obj.(type) {
		case *fyne.Container:
			if len(typed.Objects) == 1 {
				if _, ok := typed.Objects[0].(*widget.Label); ok {
					prefixContainer = typed
				}
			}
		case *pathValueLabel:
			pathValue = typed
		}
	}
	return prefixContainer, pathValue
}

func inferCHSFromImage(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("hard disk image path is required")
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("hard disk image path must be a file")
	}

	const (
		bytesPerSector  = int64(512)
		sectorsPerTrack = int64(63)
		heads           = int64(16)
	)

	unit := bytesPerSector * sectorsPerTrack * heads
	if info.Size() <= 0 || info.Size()%unit != 0 {
		return "", fmt.Errorf("could not infer CHS from image size")
	}

	cylinders := info.Size() / unit
	if cylinders <= 0 {
		return "", fmt.Errorf("could not infer CHS from image size")
	}

	return fmt.Sprintf("%d,%d,%d,%d", bytesPerSector, sectorsPerTrack, heads, cylinders), nil
}

func initialEditorWindowSize(app fyne.App) fyne.Size {
	return initialWindowSize(
		app,
		defaultEditorWindowWidth,
		defaultEditorWindowHeight,
		minEditorWindowWidth,
		minEditorWindowHeight,
	)
}
