package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestBuildSummaryObjectsCreatesPathValueRows(t *testing.T) {
	test.NewApp()

	objects := buildSummaryObjects([]summaryLine{
		{Label: "CPU", Value: "dynamic / pentium"},
		{Label: "HDD Image", Value: "/tmp/disks/game.img", IsPath: true},
	})

	if len(objects) != 2 {
		t.Fatalf("len(objects) = %d, want 2", len(objects))
	}

	if objects[0].Text != "CPU" {
		t.Fatalf("objects[0].Text = %q, want %q", objects[0].Text, "CPU")
	}
	if _, ok := objects[0].Widget.(*widget.Label); !ok {
		t.Fatalf("objects[0].Widget = %T, want *widget.Label", objects[0].Widget)
	}

	if objects[1].Text != "HDD Image" {
		t.Fatalf("objects[1].Text = %q, want %q", objects[1].Text, "HDD Image")
	}
	pathValue, ok := objects[1].Widget.(*pathValueLabel)
	if !ok {
		t.Fatalf("objects[1].Widget = %T, want *pathValueLabel", objects[1].Widget)
	}
	if got := pathValue.fullText; got != "/tmp/disks/game.img" {
		t.Fatalf("pathValue.fullText = %q, want %q", got, "/tmp/disks/game.img")
	}
}

func TestInitialMainWindowSizeUsesFallbackWithoutScreenSize(t *testing.T) {
	size := initialMainWindowSize(test.NewApp())

	want := fyne.NewSize(defaultMainWindowWidth, defaultMainWindowHeight)
	if size != want {
		t.Fatalf("initialMainWindowSize() = %#v, want %#v", size, want)
	}
}

func TestInitialMainWindowSizeUsesTwoThirdsOfScreen(t *testing.T) {
	app := &screenSizeApp{App: test.NewApp(), size: fyne.NewSize(1800, 1200)}

	size := initialMainWindowSize(app)
	want := fyne.NewSize(1200, 800)
	if size != want {
		t.Fatalf("initialMainWindowSize() = %#v, want %#v", size, want)
	}
}

func TestInitialMainWindowSizeHonorsMinimums(t *testing.T) {
	app := &screenSizeApp{App: test.NewApp(), size: fyne.NewSize(600, 400)}

	size := initialMainWindowSize(app)
	want := fyne.NewSize(minMainWindowWidth, minMainWindowHeight)
	if size != want {
		t.Fatalf("initialMainWindowSize() = %#v, want %#v", size, want)
	}
}

func TestInitialEditorWindowSizeUsesFallbackWithoutScreenSize(t *testing.T) {
	size := initialEditorWindowSize(test.NewApp())

	want := fyne.NewSize(defaultEditorWindowWidth, defaultEditorWindowHeight)
	if size != want {
		t.Fatalf("initialEditorWindowSize() = %#v, want %#v", size, want)
	}
}

func TestInitialEditorWindowSizeUsesTwoThirdsOfScreen(t *testing.T) {
	app := &screenSizeApp{App: test.NewApp(), size: fyne.NewSize(1800, 1200)}

	size := initialEditorWindowSize(app)
	want := fyne.NewSize(1200, 800)
	if size != want {
		t.Fatalf("initialEditorWindowSize() = %#v, want %#v", size, want)
	}
}

func TestInitialEditorWindowSizeHonorsMinimums(t *testing.T) {
	app := &screenSizeApp{App: test.NewApp(), size: fyne.NewSize(600, 400)}

	size := initialEditorWindowSize(app)
	want := fyne.NewSize(minEditorWindowWidth, minEditorWindowHeight)
	if size != want {
		t.Fatalf("initialEditorWindowSize() = %#v, want %#v", size, want)
	}
}

type screenSizeApp struct {
	fyne.App
	size fyne.Size
}

func (a *screenSizeApp) Driver() fyne.Driver {
	return &screenSizeDriver{Driver: a.App.Driver(), size: a.size}
}

type screenSizeDriver struct {
	fyne.Driver
	size fyne.Size
}

func (d *screenSizeDriver) ScreenSize() fyne.Size {
	return d.size
}
