package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
)

func ErrorDialog(err error, masterWindow fyne.Window) {
	dialog.ShowError(err, masterWindow)
}

func InfoDialog(message string, masterWindow fyne.Window) {
	dialog.ShowInformation("Success", message, masterWindow)
}
