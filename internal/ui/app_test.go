package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestNewSummaryPathRowUsesExpectedWidgetShape(t *testing.T) {
	test.NewApp()

	row := newSummaryPathRow("HDD Image", "/tmp/disks/game.img")

	containerRow := row.(*fyne.Container)
	prefix, pathValue := findSummaryRowParts(t, containerRow)

	if got := prefix.Text; got != "HDD Image: " {
		t.Fatalf("prefix.Text = %q, want %q", got, "HDD Image: ")
	}
	if got := pathValue.fullText; got != "/tmp/disks/game.img" {
		t.Fatalf("pathValue.fullText = %q, want %q", got, "/tmp/disks/game.img")
	}
}

func TestBuildSummaryObjectsCreatesPathValueRows(t *testing.T) {
	test.NewApp()

	objects := buildSummaryObjects([]summaryLine{
		{Label: "CPU", Value: "dynamic / pentium"},
		{Label: "HDD Image", Value: "/tmp/disks/game.img", IsPath: true},
	})

	if len(objects) != 2 {
		t.Fatalf("len(objects) = %d, want 2", len(objects))
	}
	if _, ok := objects[0].(*widget.Label); !ok {
		t.Fatalf("objects[0] = %T, want *widget.Label", objects[0])
	}

	row := objects[1].(*fyne.Container)
	_, pathValue := findSummaryRowParts(t, row)
	if pathValue == nil {
		t.Fatal("pathValue = nil, want *pathValueLabel")
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

func findSummaryRowParts(t *testing.T, row *fyne.Container) (*widget.Label, *pathValueLabel) {
	t.Helper()

	var prefix *widget.Label
	var pathValue *pathValueLabel
	for _, obj := range row.Objects {
		switch typed := obj.(type) {
		case *widget.Label:
			prefix = typed
		case *pathValueLabel:
			pathValue = typed
		}
	}
	if prefix == nil || pathValue == nil {
		t.Fatalf("unexpected summary row shape: %#v", row.Objects)
	}
	return prefix, pathValue
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
