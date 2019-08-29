package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func Text(text string) *widget.Label {
	return widget.NewLabel(text)
}

func BoldText(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
}

func LoadingText(text string) *widget.Label {
	if text == "" {
		text = "Loading"
	}
	text += "..."
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Italic: true})
}
