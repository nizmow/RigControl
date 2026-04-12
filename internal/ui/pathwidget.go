package ui

import (
	"math"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type pathValueLabel struct {
	widget.BaseWidget
	fullText string
	label    *widget.Label
	tooltip  *widget.PopUp
}

var _ desktop.Hoverable = (*pathValueLabel)(nil)

func newPathValueLabel(text string) *pathValueLabel {
	p := &pathValueLabel{
		label: widget.NewLabel(""),
	}
	p.label.Alignment = fyne.TextAlignLeading
	p.label.Wrapping = fyne.TextWrapOff
	p.SetText(text)
	p.ExtendBaseWidget(p)
	return p
}

func (p *pathValueLabel) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.label)
}

func (p *pathValueLabel) MinSize() fyne.Size {
	textSize := theme.TextSize()
	return fyne.MeasureText("…", textSize, fyne.TextStyle{})
}

func (p *pathValueLabel) Resize(size fyne.Size) {
	p.BaseWidget.Resize(size)
	p.refreshDisplay()
}

func (p *pathValueLabel) SetText(text string) {
	p.fullText = strings.TrimSpace(text)
	p.refreshDisplay()
}

func (p *pathValueLabel) MouseIn(*desktop.MouseEvent) {
	p.showTooltip()
}

func (p *pathValueLabel) MouseMoved(*desktop.MouseEvent) {
}

func (p *pathValueLabel) MouseOut() {
	if p.tooltip != nil {
		p.tooltip.Hide()
		p.tooltip = nil
	}
}

func (p *pathValueLabel) refreshDisplay() {
	if p.label == nil {
		return
	}
	width := p.Size().Width
	if width <= 0 {
		p.label.SetText(p.fullText)
		return
	}
	p.label.SetText(truncatePathText(p.fullText, width, theme.TextSize(), fyne.TextStyle{}))
}

func (p *pathValueLabel) showTooltip() {
	if p.fullText == "" || p.label.Text == p.fullText {
		return
	}
	canvas := fyne.CurrentApp().Driver().CanvasForObject(p)
	if canvas == nil {
		return
	}
	if p.tooltip != nil {
		p.tooltip.Hide()
	}

	full := widget.NewLabel(p.fullText)
	full.Wrapping = fyne.TextWrapBreak
	content := container.NewPadded(full)
	p.tooltip = widget.NewPopUp(content, canvas)
	padding := theme.Padding() * 6
	textWidth := fyne.MeasureText(p.fullText, theme.TextSize(), fyne.TextStyle{}).Width
	maxWidth := canvas.Size().Width - theme.Padding()*8
	if maxWidth < 200 {
		maxWidth = 200
	}
	width := float32(math.Min(float64(maxWidth), math.Max(360, float64(textWidth+padding))))
	p.tooltip.Resize(fyne.NewSize(width, content.MinSize().Height))
	p.tooltip.ShowAtRelativePosition(fyne.NewPos(0, p.Size().Height), p)
}

func truncatePathText(text string, maxWidth, textSize float32, style fyne.TextStyle) string {
	text = strings.TrimSpace(text)
	if text == "" || maxWidth <= 0 {
		return text
	}
	if fyne.MeasureText(text, textSize, style).Width <= maxWidth {
		return text
	}

	if candidate := truncatePathBySegments(text, maxWidth, textSize, style); candidate != "" {
		return candidate
	}

	return truncateLeftText(text, maxWidth, textSize, style)
}

func truncatePathBySegments(text string, maxWidth, textSize float32, style fyne.TextStyle) string {
	normalized := filepath.ToSlash(text)
	parts := strings.FieldsFunc(normalized, func(r rune) bool {
		return r == '/'
	})
	if len(parts) == 0 {
		return ""
	}

	last := parts[len(parts)-1]
	candidate := ""
	if fyne.MeasureText(last, textSize, style).Width <= maxWidth {
		candidate = last
	} else {
		candidate = ".../" + last
		if fyne.MeasureText(candidate, textSize, style).Width > maxWidth {
			return ""
		}
	}
	for i := len(parts) - 2; i >= 0; i-- {
		next := ".../" + strings.Join(parts[i:], "/")
		if fyne.MeasureText(next, textSize, style).Width > maxWidth {
			break
		}
		candidate = next
	}
	return candidate
}

func truncateLeftText(text string, maxWidth, textSize float32, style fyne.TextStyle) string {
	text = strings.TrimSpace(text)
	if text == "" || maxWidth <= 0 {
		return text
	}

	const ellipsis = "..."
	ellipsisWidth := fyne.MeasureText(ellipsis, textSize, style).Width
	if ellipsisWidth >= maxWidth {
		return ellipsis
	}

	runes := []rune(text)
	start := len(runes) - 1
	for start >= 0 {
		candidate := ellipsis + string(runes[start:])
		if fyne.MeasureText(candidate, textSize, style).Width <= maxWidth {
			return candidate
		}
		start--
	}
	return ellipsis
}
