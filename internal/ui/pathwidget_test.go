package ui

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
)

func TestTruncateLeftTextKeepsPathTail(t *testing.T) {
	test.NewApp()

	full := "/very/long/path/to/game/disks/install-disk-02.img"
	width := fyne.MeasureText(".../disks/install-disk-02.img", 14, textStyle()).Width
	got := truncatePathText(full, width, 14, textStyle())

	if got != ".../disks/install-disk-02.img" {
		t.Fatalf("truncatePathText() = %q, want trailing path segments", got)
	}
}

func TestTruncateLeftTextLeavesShortTextAlone(t *testing.T) {
	test.NewApp()

	full := "/tmp/disk.img"
	got := truncatePathText(full, 1000, 14, textStyle())

	if got != full {
		t.Fatalf("truncatePathText() = %q, want %q", got, full)
	}
}

func textStyle() fyne.TextStyle {
	return fyne.TextStyle{}
}
