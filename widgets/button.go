package widgets

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/raedahgroup/dcrseedgen/helper"
)

func NewButton() *widget.Button {
	return new(widget.Button)
}

func LayoutButton(button *widget.Button, txt string, theme *helper.Theme, gtx *layout.Context) {
	col := helper.WhiteColor
	bgcol := helper.DecredDarkBlueColor
	if !button.Active() {
		col = color.RGBA{255, 255, 255, 255}
		bgcol = helper.DecredDarkBlueColor
	}
	st := layout.Stack{Alignment: layout.Center}
	hmin := gtx.Constraints.Width.Min

	font := text.Font{
		Size: unit.Dp(theme.TextSize.V),
	}

	lbl := st.Rigid(gtx, func() {
		gtx.Constraints.Width.Min = hmin
		gtx.Constraints.Height.Min = 30
		layout.Align(layout.Center).Layout(gtx, func() {
			layout.UniformInset(unit.Dp(9)).Layout(gtx, func() {
				paint.ColorOp{Color: col}.Add(gtx.Ops)
				widget.Label{Alignment: text.Middle}.Layout(gtx, theme.Shaper, font, txt)
			})
		})
		pointer.RectAreaOp{Rect: image.Rectangle{Max: gtx.Dimensions.Size}}.Add(gtx.Ops)
		button.Layout(gtx)
	})
	bg := st.Expand(gtx, func() {
		rr := float32(gtx.Px(unit.Dp(0)))
		rrect(gtx.Ops,
			float32(gtx.Constraints.Width.Min),
			float32(gtx.Constraints.Height.Min),
			rr, rr, rr, rr,
		)
		fill(gtx, bgcol)
		for _, c := range button.History() {
			drawInk(gtx, c)
		}
	})
	st.Layout(gtx, bg, lbl)
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: f32.Point{X: float32(r.Min.X), Y: float32(r.Min.Y)},
		Max: f32.Point{X: float32(r.Max.X), Y: float32(r.Max.Y)},
	}
}
