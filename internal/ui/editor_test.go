package ui

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"rigcontrol/internal/machine"
)

func TestInferCHSFromImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "disk.img")
	size := int64(512 * 63 * 16 * 142)
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := inferCHSFromImage(path)
	if err != nil {
		t.Fatalf("inferCHSFromImage() error = %v", err)
	}
	want := "512,63,16,142"
	if got != want {
		t.Fatalf("inferCHSFromImage() = %q, want %q", got, want)
	}
}

func TestInferCHSFromImageRejectsNonDivisibleSize(t *testing.T) {
	path := filepath.Join(t.TempDir(), "disk.img")
	if err := os.WriteFile(path, make([]byte, 12345), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := inferCHSFromImage(path)
	if err == nil {
		t.Fatal("inferCHSFromImage() error = nil, want size error")
	}
}

func TestSetHardDiskImagePathAutoFillsCHS(t *testing.T) {
	test.NewApp()

	path := filepath.Join(t.TempDir(), "disk.img")
	size := int64(512 * 63 * 16 * 142)
	if err := os.WriteFile(path, make([]byte, size), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	form := newProfileForm(newTestProfile())
	form.setHardDiskImagePath(path)

	if got := form.hardDiskImage.Text; got != path {
		t.Fatalf("hardDiskImage.Text = %q, want %q", got, path)
	}
	if got := form.hardDiskCHS.Text; got != "512,63,16,142" {
		t.Fatalf("hardDiskCHS.Text = %q, want %q", got, "512,63,16,142")
	}
}

func TestSetHardDiskImagePathKeepsExistingCHSWhenInferenceFails(t *testing.T) {
	test.NewApp()

	path := filepath.Join(t.TempDir(), "disk.img")
	if err := os.WriteFile(path, make([]byte, 12345), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	form := newProfileForm(newTestProfile())
	form.hardDiskCHS.SetText("manual-value")
	form.setHardDiskImagePath(path)

	if got := form.hardDiskImage.Text; got != path {
		t.Fatalf("hardDiskImage.Text = %q, want %q", got, path)
	}
	if got := form.hardDiskCHS.Text; got != "manual-value" {
		t.Fatalf("hardDiskCHS.Text = %q, want %q", got, "manual-value")
	}
}

func TestProfileFormRoundTripsMouseSettings(t *testing.T) {
	test.NewApp()

	profile := newTestProfile()
	profile.MouseCapture = "seamless"
	profile.MouseRawInput = false
	profile.DOSMouseImmediate = true

	form := newProfileForm(profile)
	got, err := form.profile()
	if err != nil {
		t.Fatalf("form.profile() error = %v", err)
	}

	if got.MouseCapture != "seamless" {
		t.Fatalf("MouseCapture = %q, want %q", got.MouseCapture, "seamless")
	}
	if got.MouseRawInput != false {
		t.Fatalf("MouseRawInput = %v, want false", got.MouseRawInput)
	}
	if got.DOSMouseImmediate != true {
		t.Fatalf("DOSMouseImmediate = %v, want true", got.DOSMouseImmediate)
	}
}

func TestAddFloppyDiskPathAppendsAndSelects(t *testing.T) {
	test.NewApp()

	form := newProfileForm(newTestProfile())
	form.addFloppyDiskPath("/tmp/disk1.img")
	form.addFloppyDiskPath("/tmp/disk2.img")

	if got := form.floppyDiskImages; len(got) != 2 || got[0] != "/tmp/disk1.img" || got[1] != "/tmp/disk2.img" {
		t.Fatalf("floppyDiskImages = %#v, want appended images", got)
	}
	if got := form.floppySelection; got != 1 {
		t.Fatalf("floppySelection = %d, want 1", got)
	}
}

func TestMoveSelectedFloppyPreservesOrder(t *testing.T) {
	test.NewApp()

	form := newProfileForm(newTestProfile())
	form.floppyDiskImages = []string{"/tmp/disk1.img", "/tmp/disk2.img", "/tmp/disk3.img"}
	form.floppyList.Refresh()
	form.floppyList.Select(1)

	form.moveSelectedFloppyDown()
	if got := form.floppyDiskImages; got[0] != "/tmp/disk1.img" || got[1] != "/tmp/disk3.img" || got[2] != "/tmp/disk2.img" {
		t.Fatalf("after move down floppyDiskImages = %#v", got)
	}
	if got := form.floppySelection; got != 2 {
		t.Fatalf("after move down floppySelection = %d, want 2", got)
	}

	form.moveSelectedFloppyUp()
	if got := form.floppyDiskImages; got[0] != "/tmp/disk1.img" || got[1] != "/tmp/disk2.img" || got[2] != "/tmp/disk3.img" {
		t.Fatalf("after move up floppyDiskImages = %#v", got)
	}
	if got := form.floppySelection; got != 1 {
		t.Fatalf("after move up floppySelection = %d, want 1", got)
	}
}

func TestRemoveSelectedFloppyKeepsNearestSelection(t *testing.T) {
	test.NewApp()

	form := newProfileForm(newTestProfile())
	form.floppyDiskImages = []string{"/tmp/disk1.img", "/tmp/disk2.img", "/tmp/disk3.img"}
	form.floppyList.Refresh()
	form.floppyList.Select(1)

	form.removeSelectedFloppy()

	if got := form.floppyDiskImages; len(got) != 2 || got[0] != "/tmp/disk1.img" || got[1] != "/tmp/disk3.img" {
		t.Fatalf("floppyDiskImages = %#v, want disk2 removed", got)
	}
	if got := form.floppySelection; got != 1 {
		t.Fatalf("floppySelection = %d, want 1", got)
	}
}

func TestUpdateFloppyListRowUsesExpectedWidgetShape(t *testing.T) {
	test.NewApp()

	row := newFloppyListRow()
	updateFloppyListRow(row, 7, "/tmp/disks/install-disk-07.img")

	containerRow := row.(*fyne.Container)
	prefixContainer, pathValue := findFloppyRowParts(containerRow)
	if prefixContainer == nil || pathValue == nil {
		t.Fatalf("unexpected floppy row shape: %#v", containerRow.Objects)
	}
	prefix := prefixContainer.Objects[0].(*widget.Label)

	if got := prefix.Text; got != "7. " {
		t.Fatalf("prefix.Text = %q, want %q", got, "7. ")
	}
	if got := pathValue.fullText; got != "/tmp/disks/install-disk-07.img" {
		t.Fatalf("pathValue.fullText = %q, want %q", got, "/tmp/disks/install-disk-07.img")
	}
}

func TestProfileFormRoundTripsFixedCycles(t *testing.T) {
	test.NewApp()

	profile := newTestProfile()
	profile.Cycles = "3000"
	profile.FixedCycles = true

	form := newProfileForm(profile)
	got, err := form.profile()
	if err != nil {
		t.Fatalf("form.profile() error = %v", err)
	}

	if got.Cycles != "3000" {
		t.Fatalf("Cycles = %q, want %q", got.Cycles, "3000")
	}
	if got.FixedCycles != true {
		t.Fatalf("FixedCycles = %v, want true", got.FixedCycles)
	}
}

func newTestProfile() machine.Profile {
	return machine.NewProfile()
}
