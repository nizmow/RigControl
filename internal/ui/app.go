package ui

import (
	"errors"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"rigcontrol/internal/dosbox"
	"rigcontrol/internal/machine"
)

const (
	defaultMainWindowWidth  = 1280
	defaultMainWindowHeight = 800
	minMainWindowWidth      = 960
	minMainWindowHeight     = 640
)

type appUI struct {
	window           fyne.Window
	saveProfiles     func([]machine.Profile) error
	profiles         []machine.Profile
	selectedIndex    int
	list             *widget.List
	nameLabel        *widget.Label
	descriptionLabel *widget.Label
	summaryBox       *fyne.Container
}

func RunWithProfiles(profiles []machine.Profile, saveProfiles func([]machine.Profile) error) error {
	if len(profiles) == 0 {
		return errors.New("no machine profiles available")
	}

	fyneApp := app.NewWithID("com.nizmow.rigcontrol")
	window := fyneApp.NewWindow("RigControl")
	window.Resize(initialMainWindowSize(fyneApp))

	ui := newAppUI(window, profiles, saveProfiles)
	window.SetContent(ui.content())
	ui.list.Select(0)
	window.ShowAndRun()

	return nil
}

func newAppUI(window fyne.Window, profiles []machine.Profile, saveProfiles func([]machine.Profile) error) *appUI {
	ui := &appUI{
		window:           window,
		saveProfiles:     saveProfiles,
		profiles:         profiles,
		selectedIndex:    0,
		nameLabel:        widget.NewLabel(""),
		descriptionLabel: widget.NewLabel(""),
		summaryBox:       container.NewVBox(),
	}

	ui.nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.descriptionLabel.Wrapping = fyne.TextWrapWord

	ui.list = widget.NewList(
		func() int { return len(ui.profiles) },
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(ui.profiles[id].Name)
		},
	)
	ui.list.OnSelected = ui.selectProfile

	return ui
}

func (ui *appUI) content() fyne.CanvasObject {
	addButton := widget.NewButton("Add Machine", func() {
		showProfileEditor(ui.window, machine.NewProfile(), ui.addProfile)
	})

	editButton := widget.NewButton("Edit Machine", func() {
		showProfileEditor(ui.window, ui.selectedProfile(), ui.updateSelectedProfile)
	})

	previewButton := widget.NewButton("Preview Config", ui.showConfigPreview)
	launchButton := widget.NewButton("Launch DOSBox", ui.launchSelectedProfile)

	leftPane := ui.buildLeftPane()
	rightPane := ui.buildRightPane(addButton, editButton, previewButton, launchButton)

	split := container.NewHSplit(leftPane, rightPane)
	split.SetOffset(0.28)

	return split
}

func (ui *appUI) selectedProfile() machine.Profile {
	return ui.profiles[ui.selectedIndex]
}

func (ui *appUI) addProfile(profile machine.Profile) error {
	nextProfiles := append(append([]machine.Profile(nil), ui.profiles...), profile)
	if err := ui.persistProfiles(nextProfiles); err != nil {
		return err
	}

	ui.profiles = nextProfiles
	ui.list.Refresh()
	ui.list.Select(len(ui.profiles) - 1)
	return nil
}

func (ui *appUI) updateSelectedProfile(profile machine.Profile) error {
	nextProfiles := append([]machine.Profile(nil), ui.profiles...)
	nextProfiles[ui.selectedIndex] = profile
	if err := ui.persistProfiles(nextProfiles); err != nil {
		return err
	}

	ui.profiles = nextProfiles
	ui.list.Refresh()
	ui.refreshSelectedProfile()
	return nil
}

func (ui *appUI) persistProfiles(profiles []machine.Profile) error {
	if ui.saveProfiles == nil {
		return nil
	}
	return ui.saveProfiles(profiles)
}

func (ui *appUI) launchSelectedProfile() {
	profile := ui.selectedProfile()
	if err := machine.ValidateProfile(profile); err != nil {
		dialog.ShowError(err, ui.window)
		return
	}

	configPath, err := writeTempConfig(profile)
	if err != nil {
		dialog.ShowError(err, ui.window)
		return
	}

	if err := launchProfile(profile, configPath); err != nil {
		dialog.ShowError(err, ui.window)
		return
	}

	dialog.ShowInformation("DOSBox Started", "Launched DOSBox Staging with "+filepath.Base(configPath), ui.window)
}

func (ui *appUI) buildLeftPane() fyne.CanvasObject {
	return container.NewPadded(container.NewBorder(
		widget.NewLabel("Machine Types"),
		nil,
		nil,
		nil,
		ui.list,
	))
}

func (ui *appUI) buildRightPane(objects ...fyne.CanvasObject) fyne.CanvasObject {
	actions := container.NewHBox(objects...)
	summaryScroll := container.NewVScroll(ui.buildSummaryPane())
	return container.NewPadded(container.NewBorder(nil, actions, nil, nil, summaryScroll))
}

func (ui *appUI) buildSummaryPane() fyne.CanvasObject {
	return container.NewVBox(
		ui.nameLabel,
		ui.descriptionLabel,
		widget.NewSeparator(),
		widget.NewLabel("Configuration Summary"),
		ui.summaryBox,
	)
}

func (ui *appUI) refreshSelectedProfile() {
	profile := ui.selectedProfile()
	ui.nameLabel.SetText(profile.Name)
	ui.descriptionLabel.SetText(profile.Description)
	ui.summaryBox.Objects = buildSummaryObjects(profileSummaryLines(profile))
	ui.summaryBox.Refresh()
}

func (ui *appUI) selectProfile(id widget.ListItemID) {
	if id < 0 || id >= len(ui.profiles) {
		return
	}

	ui.selectedIndex = id
	ui.refreshSelectedProfile()
}

func (ui *appUI) showConfigPreview() {
	previewWindow := fyne.CurrentApp().NewWindow("DOSBox Config Preview")
	previewWindow.Resize(fyne.NewSize(760, 620))

	preview := widget.NewMultiLineEntry()
	preview.Wrapping = fyne.TextWrapOff
	preview.SetText(dosbox.Render(ui.selectedProfile()))
	preview.Disable()

	closeButton := widget.NewButton("Close", func() {
		previewWindow.Close()
	})

	previewWindow.SetContent(container.NewBorder(
		widget.NewLabel("Generated DOSBox Config"),
		container.NewHBox(closeButton),
		nil,
		nil,
		container.NewPadded(preview),
	))
	previewWindow.CenterOnScreen()
	previewWindow.Show()
}

func buildSummaryObjects(lines []summaryLine) []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, len(lines))
	for _, line := range lines {
		if line.IsPath {
			objects = append(objects, newSummaryPathRow(line.Label, line.Value))
			continue
		}
		label := widget.NewLabel(line.Label + ": " + line.Value)
		label.Wrapping = fyne.TextWrapBreak
		objects = append(objects, label)
	}
	return objects
}

func newSummaryPathRow(label, value string) fyne.CanvasObject {
	prefix := widget.NewLabel(label + ": ")
	prefix.Wrapping = fyne.TextWrapOff
	return container.NewBorder(nil, nil, prefix, nil, newPathValueLabel(value))
}

func initialMainWindowSize(app fyne.App) fyne.Size {
	return initialWindowSize(
		app,
		defaultMainWindowWidth,
		defaultMainWindowHeight,
		minMainWindowWidth,
		minMainWindowHeight,
	)
}

func initialWindowSize(app fyne.App, defaultWidth, defaultHeight, minWidth, minHeight float32) fyne.Size {
	if app != nil && app.Driver() != nil {
		if screenDriver, ok := app.Driver().(interface{ ScreenSize() fyne.Size }); ok {
			screen := screenDriver.ScreenSize()
			if screen.Width > 0 && screen.Height > 0 {
				return fyne.NewSize(
					maxFloat32(minWidth, screen.Width*2/3),
					maxFloat32(minHeight, screen.Height*2/3),
				)
			}
		}
	}
	return fyne.NewSize(defaultWidth, defaultHeight)
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}
