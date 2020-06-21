// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Editor struct {
	material.EditorStyle
	lineColor color.RGBA
}

func (t *Theme) Editor(hint string, e *widget.Editor) Editor {
	return Editor{
		EditorStyle: material.Editor(t.Theme, e, hint),
		lineColor:   t.Color.Hint,
	}
}

func (e Editor) Layout(gtx layout.Context) layout.Dimensions {
	col := e.lineColor

	if e.Editor.Focused() {
		e.lineColor = color.RGBA{41, 112, 255, 255}
	} else {
		e.lineColor = col
	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return e.EditorStyle.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return e.editorLine(gtx)
		}),
	)
}

func (e Editor) editorLine(gtx layout.Context) layout.Dimensions {
	return layout.Flex{}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			rect := f32.Rectangle{
				Max: f32.Point{
					X: float32(gtx.Constraints.Max.X),
					Y: 2,
				},
			}
			op.TransformOp{}.Offset(f32.Point{
				X: 0,
				Y: 0,
			}).Add(gtx.Ops)
			paint.ColorOp{Color: e.lineColor}.Add(gtx.Ops)
			paint.PaintOp{Rect: rect}.Add(gtx.Ops)

			return layout.Dimensions{}
		}),
	)
}
