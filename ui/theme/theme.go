package theme

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"

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

func New() *Theme {
	t := &Theme{
		Shaper: font.Default(),
	}
	t.Color.Primary = keyblue
	t.Color.Text = darkblue
	t.Color.Hint = rgb(0xbbbbbb)
	t.Color.InvText = rgb(0xffffff)
	t.Color.Overlay = rgb(0x000000)
	t.Color.Background = argb(0x22444444)
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

func (t *Theme) alert(gtx *layout.Context, txt string, bgColor color.RGBA) {
	bgColor.A = 200

	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(2)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			Fill(gtx, bgColor)
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
			layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
				lbl := t.Body2(txt)
				lbl.Color = t.Color.Surface
				lbl.Layout(gtx)
			})
		}),
	)
}

func (t *Theme) ErrorAlert(gtx *layout.Context, txt string) {
	t.alert(gtx, txt, t.Color.Danger)
}

func (t *Theme) SuccessAlert(gtx *layout.Context, txt string) {
	t.alert(gtx, txt, t.Color.Success)
}

func rgb(c uint32) color.RGBA {
	return argb(0xff000000 | c)
}

func argb(c uint32) color.RGBA {
	return color.RGBA{A: uint8(c >> 24), R: uint8(c >> 16), G: uint8(c >> 8), B: uint8(c)}
}

func ToMax(gtx *layout.Context) {
	gtx.Constraints.Width.Min = gtx.Constraints.Width.Max
	gtx.Constraints.Height.Min = gtx.Constraints.Height.Max
}

func Fill(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	fill(gtx, cs.Width.Min, cs.Height.Min, col)
}

func FillMax(gtx *layout.Context, col color.RGBA) {
	cs := gtx.Constraints
	fill(gtx, cs.Width.Max, cs.Height.Max, col)
}

func fill(gtx *layout.Context, width, height int, col color.RGBA) {
	d := image.Point{X: width, Y: height}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: col}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}

func MustIcon(ic *Icon, err error) *Icon {
	if err != nil {
		panic(err)
	}
	return ic
}
