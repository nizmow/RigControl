package ui

import (
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	ttwidget "github.com/dweymouth/fyne-tooltip/widget"
)

type pathValueLabel struct {
	ttwidget.ToolTipWidget
	fullText string
	label    *widget.Label
}

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

func (p *pathValueLabel) refreshDisplay() {
	if p.label == nil {
		return
	}
	width := p.Size().Width
	if width <= 0 {
		p.label.SetText(p.fullText)
		p.SetToolTip("")
		return
	}
	displayText := truncatePathText(p.fullText, width, theme.TextSize(), fyne.TextStyle{})
	p.label.SetText(displayText)

	// Only show tooltip if text was truncated
	if displayText != p.fullText {
		p.SetToolTip(p.fullText)
	} else {
		p.SetToolTip("")
	}
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
