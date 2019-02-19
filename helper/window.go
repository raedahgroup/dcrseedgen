package helper

import (
	"image"

	"github.com/aarzilli/nucular"
)

type Window struct {
	*nucular.Window
}

func NewWindow(title string, w *nucular.Window, flags nucular.WindowFlags) *Window {
	if window := w.GroupBegin(title, flags); window != nil {
		return &Window{
			window,
		}
	}
	return nil
}

func (w *Window) NewWindow(title string, flags nucular.WindowFlags) *Window {
	if window := w.GroupBegin(title, flags); window != nil {
		return &Window{
			window,
		}
	}
	return nil
}

func (w *Window) End() {
	w.GroupEnd()
}

func (w *Window) Style() {
	style := w.Master().Style()
	style.GroupWindow.Padding = image.Point{20, 20}

	w.Master().SetStyle(style)
}

func DrawPageHeader(window *nucular.Window) {
	window.Row(70).Dynamic(1)

	style := window.Master().Style()
	style.GroupWindow.FixedBackground.Data.Color = whiteColor
	style.NormalWindow.Padding = noPadding
	style.GroupWindow.Padding = noPadding
	window.Master().SetStyle(style)

	if group := window.GroupBegin("header", 0); window != nil {
		// style window
		group.Row(60).Dynamic(1)
		group.Image(logo)
		group.GroupEnd()
	}
	window.Row(10).Dynamic(1)
	window.Label("", "LC")

	style = window.Master().Style()
	style.GroupWindow.FixedBackground.Data.Color = colorTable.ColorWindow
	style.GroupWindow.Padding = groupWindowPadding
	window.Master().SetStyle(style)
}
