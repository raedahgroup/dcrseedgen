package handlers

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type (
	AddressHandler struct {
	}
)

func (h *AddressHandler) BeforeRender() {

}

func (h *AddressHandler) Render() fyne.CanvasObject {
	return widget.NewHBox()
}
