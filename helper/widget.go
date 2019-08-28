package helper

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
