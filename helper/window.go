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
	style := window.Master().Style()
	style.GroupWindow.Padding = image.Point{12, 0}
	style.GroupWindow.FixedBackground.Data.Color = whiteColor
	window.Master().SetStyle(style)

	window.Row(70).Dynamic(1)
	if group := window.GroupBegin("header", 0); window != nil {
		group.Row(60).Dynamic(1)
		group.Image(logo)
		group.GroupEnd()
	}
	style = window.Master().Style()
	style.GroupWindow.FixedBackground.Data.Color = colorTable.ColorWindow
	style.GroupWindow.Padding = groupWindowPadding
	window.Master().SetStyle(style)
}
