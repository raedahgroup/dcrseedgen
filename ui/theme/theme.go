package theme

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	// decred primary colors

	keyblue = rgb(0x2970ff)
	//turquiose = rgb(0x2ed6a1)
	darkblue = rgb(0x091440)

	// decred complemetary colors

	//lightblue = rgb(0x70cbff)
	//orange    = rgb(0xed6d47)
	green = rgb(0x41bf53)
)

type Theme struct {
	*material.Theme
	Shaper text.Shaper
	Color  struct {
		Primary    color.RGBA
		Secondary  color.RGBA
		Text       color.RGBA
		Hint       color.RGBA
		Overlay    color.RGBA
		InvText    color.RGBA
		Success    color.RGBA
		Danger     color.RGBA
		Background color.RGBA
		Gray       color.RGBA
		Black      color.RGBA
		Surface    color.RGBA
	}
	TextSize           unit.Value
	radioCheckedIcon   *Icon
	radioUncheckedIcon *Icon
}

func New(col *text.Collection) *Theme {
	t := &Theme{
		Theme: material.NewTheme(col),
	}
	t.Color.Primary = keyblue
	t.Color.Text = darkblue
	t.Color.Hint = rgb(0xbbbbbb)
	t.Color.InvText = rgb(0xffffff)
	t.Color.Overlay = rgb(0x000000)
	t.Color.Background = argb(0x33444444)
	t.Color.Success = green
	t.Color.Danger = rgb(0xff0000)
	t.Color.Gray = rgb(0x808080)
	t.Color.Black = rgb(0x000000)
	t.Color.Surface = rgb(0xffffff)
	t.TextSize = unit.Sp(16)

	t.radioCheckedIcon = MustIcon(NewIcon(icons.ToggleRadioButtonChecked))
	t.radioUncheckedIcon = MustIcon(NewIcon(icons.ToggleRadioButtonUnchecked))

	return t
}

func (t *Theme) alert(gtx layout.Context, txt string, bgColor color.RGBA) layout.Dimensions {
	bgColor.A = 200

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rr := float32(gtx.Px(unit.Dp(2)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			return Fill(gtx, bgColor)
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min.X = gtx.Constraints.Max.Y
			return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				lbl := t.Body2(txt)
				lbl.Color = t.Color.Surface
				return lbl.Layout(gtx)
			})
		}),
	)
}

func (t *Theme) ErrorAlert(gtx layout.Context, txt string) layout.Dimensions {
	return t.alert(gtx, txt, t.Color.Danger)
}

func (t *Theme) SuccessAlert(gtx layout.Context, txt string) layout.Dimensions {
	return t.alert(gtx, txt, t.Color.Success)
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func ToMax(gtx layout.Context) {
	gtx.Constraints.Min = gtx.Constraints.Max
}

func Fill(gtx layout.Context, col color.RGBA) layout.Dimensions {
	cs := gtx.Constraints.Min
	return fill(gtx, cs.X, cs.Y, col)
}

func FillMax(gtx layout.Context, col color.RGBA) layout.Dimensions {
	cs := gtx.Constraints.Max
	return fill(gtx, cs.X, cs.Y, col)
}

func fill(gtx layout.Context, width, height int, col color.RGBA) layout.Dimensions {
	d := image.Point{X: width, Y: height}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	return layout.Dimensions{Size: d}
}

func Bounds(gtx layout.Context) f32.Rectangle {
	cs := gtx.Constraints
	d := cs.Min
	return f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
}

func dynamicColor(i int) color.RGBA {
	sn, cs := math.Sincos(float64(i) * math.Phi)
	return color.RGBA{
		R: 0xA0 + byte(0x30*sn),
		G: 0xA0 + byte(0x30*cs),
		B: 0xD0,
		A: 0xFF,
	}
}

func MustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
