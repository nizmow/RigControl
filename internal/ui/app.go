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

type appUI struct {
	window           fyne.Window
	saveProfiles     func([]machine.Profile) error
	profiles         []machine.Profile
	selectedIndex    int
	list             *widget.List
	nameLabel        *widget.Label
	descriptionLabel *widget.Label
	summaryLabel     *widget.Label
	preview          *widget.Entry
}

func RunWithProfiles(profiles []machine.Profile, saveProfiles func([]machine.Profile) error) error {
	if len(profiles) == 0 {
		return errors.New("no machine profiles available")
	}

	fyneApp := app.NewWithID("com.nizmow.rigcontrol")
	window := fyneApp.NewWindow("RigControl")
	window.Resize(fyne.NewSize(1080, 720))

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
		summaryLabel:     widget.NewLabel(""),
		preview:          widget.NewMultiLineEntry(),
	}

	ui.nameLabel.TextStyle = fyne.TextStyle{Bold: true}
	ui.descriptionLabel.Wrapping = fyne.TextWrapWord
	ui.summaryLabel.Wrapping = fyne.TextWrapWord
	ui.preview.Wrapping = fyne.TextWrapOff
	ui.preview.Disable()

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
		showProfileEditor(ui.window, newMachineProfile(), func(created machine.Profile) error {
			nextProfiles := append(append([]machine.Profile(nil), ui.profiles...), created)

			if ui.saveProfiles != nil {
				if err := ui.saveProfiles(nextProfiles); err != nil {
					return err
				}
			}

			ui.profiles = nextProfiles
			ui.list.Refresh()
			ui.list.Select(len(ui.profiles) - 1)
			return nil
		})
	})

	editButton := widget.NewButton("Edit Machine", func() {
		showProfileEditor(ui.window, ui.selectedProfile(), func(updated machine.Profile) error {
			nextProfiles := append([]machine.Profile(nil), ui.profiles...)
			nextProfiles[ui.selectedIndex] = updated

			if ui.saveProfiles != nil {
				if err := ui.saveProfiles(nextProfiles); err != nil {
					return err
				}
			}

			ui.profiles = nextProfiles
			ui.list.Refresh()
			ui.selectProfile(ui.selectedIndex)
			return nil
		})
	})

	launchButton := widget.NewButton("Launch DOSBox", func() {
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

		if err := launchDOSBox(configPath); err != nil {
			dialog.ShowError(err, ui.window)
			return
		}

		dialog.ShowInformation("DOSBox Started", "Launched DOSBox Staging with "+filepath.Base(configPath), ui.window)
	})

	leftPane := container.NewBorder(
		widget.NewLabel("Machine Types"),
		nil,
		nil,
		nil,
		ui.list,
	)

	summaryPane := container.NewVBox(
		ui.nameLabel,
		ui.descriptionLabel,
		widget.NewSeparator(),
		widget.NewLabel("Configuration Summary"),
		ui.summaryLabel,
	)

	actions := container.NewHBox(addButton, editButton, launchButton)

	rightPane := container.NewBorder(
		nil,
		actions,
		nil,
		nil,
		func() *container.Split {
			split := container.NewVSplit(
				container.NewPadded(summaryPane),
				container.NewBorder(
					widget.NewLabel("DOSBox Config Preview"),
					nil,
					nil,
					nil,
					ui.preview,
				),
			)
			split.SetOffset(0.42)
			return split
		}(),
	)

	split := container.NewHSplit(container.NewPadded(leftPane), container.NewPadded(rightPane))
	split.SetOffset(0.28)

	return split
}

func (ui *appUI) selectedProfile() machine.Profile {
	return ui.profiles[ui.selectedIndex]
}

func (ui *appUI) selectProfile(id widget.ListItemID) {
	if id < 0 || id >= len(ui.profiles) {
		return
	}

	ui.selectedIndex = id
	profile := ui.selectedProfile()

	ui.nameLabel.SetText(profile.Name)
	ui.descriptionLabel.SetText(profile.Description)
	ui.summaryLabel.SetText(profileSummary(profile))
	ui.preview.SetText(dosbox.Render(profile))
}

func newMachineProfile() machine.Profile {
	return machine.Profile{
		Name:         "New Machine",
		Description:  "",
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
}
