package ui

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/test"

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

func newTestProfile() machine.Profile {
	return machine.NewProfile()
}
