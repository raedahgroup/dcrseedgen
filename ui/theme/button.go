// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Button struct {
	btn material.ButtonStyle
}

type IconButton struct {
	btn        *widget.Clickable
	Background color.RGBA
	Color      color.RGBA
	Animated   bool
	Icon       *Icon
	Size       unit.Value
	Padding    unit.Value
	Inset      layout.Inset
}

func (t *Theme) Button(txt string, clickable *widget.Clickable) Button {
	btn := Button{
		btn: material.Button(t.Theme, clickable, txt),
	}
	btn.btn.Background = t.Color.Primary
	return btn
}

func (t *Theme) PrimaryButton(txt string, clickable *widget.Clickable) Button {
	btn := t.Button(txt, clickable)
	btn.btn.Background = t.Color.Primary
	return btn
}

func (t *Theme) SecondaryButton(txt string, clickable *widget.Clickable) Button {
	btn := t.Button(txt, clickable)
	btn.btn.Background = t.Color.Secondary
	return btn
}

func (t *Theme) DangerButton(txt string, clickable *widget.Clickable) Button {
	btn := t.Button(txt, clickable)
	btn.btn.Background = t.Color.Danger
	return btn
}

func (t *Theme) SuccessButton(txt string, clickable *widget.Clickable) Button {
	btn := t.Button(txt, clickable)
	btn.btn.Background = t.Color.Success
	return btn
}

func (btn *Button) Layout(gtx layout.Context) layout.Dimensions {
	return btn.btn.Layout(gtx)
}

func (t *Theme) IconButton(ic *Icon, btn *widget.Clickable) IconButton {

	return IconButton{
		btn:        btn,
		Background: t.Color.Primary,
		Color:      t.Color.Primary,
		Icon:       ic,
	}
}

func (btn *IconButton) Layout(gtx layout.Context) layout.Dimensions {
	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			size := gtx.Constraints.Min.X
			sizef := float32(size)
			rr := sizef * .5
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{X: sizef, Y: sizef}},
				NE:   rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			background := btn.Background
			if gtx.Queue == nil {
				background = mulAlpha(btn.Background, 150)
			}
			dims := fill(gtx, gtx.Constraints.Min.X, gtx.Constraints.Min.Y, background)
			for _, c := range btn.btn.History() {
				drawInk(gtx, c)
			}
			return dims
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return btn.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				size := gtx.Px(btn.Size) - 2*gtx.Px(btn.Padding)
				if btn.Icon != nil {
					btn.Icon.Color = btn.Color
					btn.Icon.Layout(gtx, unit.Px(float32(size)))
				}
				return layout.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
		}),
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			pointer.Ellipse(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
			return btn.btn.Layout(gtx)
		}),
	)
}

/**
type Button struct {
	Text string
	// Color is the text color.
	Color        color.RGBA
	Font         text.Font
	TextSize     unit.Value
	Background   color.RGBA
	CornerRadius unit.Value
	Inset        layout.Inset
	Radius       float32
	shaper       text.Shaper
}

type IconButton struct {
	Background color.RGBA
	Color      color.RGBA
	Animated   bool
	Icon       *Icon
	Size       unit.Value
	Padding    unit.Value
	Inset      layout.Inset
}

func (t *Theme) Button(txt string) Button {
	return Button{
		Text:       txt,
		Color:      rgb(0xffffff),
		Background: t.Color.Primary,
		TextSize:   t.TextSize.Scale(14.0 / 16.0),
		Inset: layout.Inset{
			Top: unit.Dp(10), Bottom: unit.Dp(10),
			Left: unit.Dp(12), Right: unit.Dp(12),
		},
		Radius: 4,
		shaper: t.Shaper,
	}
}

func (t *Theme) IconButton(icon *Icon) IconButton {
	return IconButton{
		Background: t.Color.Primary,
		Color:      t.Color.InvText,
		Icon:       icon,
		Size:       unit.Dp(56),
		Padding:    unit.Dp(16),
		Animated:   true,
	}
}

func (t *Theme) PlainIconButton(icon *Icon) IconButton {
	return IconButton{
		Background: color.RGBA{},
		Color:      t.Color.Primary,
		Icon:       icon,
		Size:       unit.Dp(56),
		Padding:    unit.Dp(0),
	}
}

func (b Button) Layout(gtx layout.Context, button *widget.Clickable) layout.Dimensions {
	col := b.Color
	bgcol := b.Background
	min := gtx.Constraints.Min

	return layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			rr := float32(gtx.Px(b.CornerRadius))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Min.X),
					Y: float32(gtx.Constraints.Min.Y),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			background := b.Background
			if gtx.Queue == nil {
				background = mulAlpha(b.Background, 150)
			}
			dims := FillMax(gtx, background)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
			return dims
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = min
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return b.Inset.Layout(gtx, w)
			})
		}),
		layout.Expanded(button.Layout),
	)
}

func (b IconButton) Layout(gtx layout.Context, button *widget.Clickable) layout.Dimensions {
	return layout.Dimensions{}

	/**layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			size := gtx.Constraints.Width.Min
			sizef := float32(size)
			rr := sizef * .5
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{X: sizef, Y: sizef}},
				NE:   rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			Fill(gtx, b.Background)
			if b.Animated {
				for _, c := range button.History() {
					drawInk(gtx, c)
				}
			}
		}),
		layout.Stacked(func() {
			layout.UniformInset(b.Padding).Layout(gtx, func() {
				size := gtx.Px(b.Size) - 2*gtx.Px(b.Padding)
				if b.Icon != nil {
					b.Icon.Color = b.Color
					b.Icon.Layout(gtx, unit.Px(float32(size)))
				}
				gtx.Dimensions = layout.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
			pointer.Ellipse(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)*
}**/

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func drawInk(gtx layout.Context, c widget.Press) {
	// duration is the number of seconds for the
	// completed animation: expand while fading in, then
	// out.
	const (
		expandDuration = float32(0.5)
		fadeDuration   = float32(0.9)
	)

	now := gtx.Now

	t := float32(now.Sub(c.Start).Seconds())

	end := c.End
	if end.IsZero() {
		// If the press hasn't ended, don't fade-out.
		end = now
	}

	endt := float32(end.Sub(c.Start).Seconds())

	// Compute the fade-in/out position in [0;1].
	var alphat float32
	{
		var haste float32
		if c.Cancelled {
			// If the press was cancelled before the inkwell
			// was fully faded in, fast forward the animation
			// to match the fade-out.
			if h := 0.5 - endt/fadeDuration; h > 0 {
				haste = h
			}
		}
		// Fade in.
		half1 := t/fadeDuration + haste
		if half1 > 0.5 {
			half1 = 0.5
		}

		// Fade out.
		half2 := float32(now.Sub(end).Seconds())
		half2 /= fadeDuration
		half2 += haste
		if half2 > 0.5 {
			// Too old.
			return
		}

		alphat = half1 + half2
	}

	// Compute the expand position in [0;1].
	sizet := t
	if c.Cancelled {
		// Freeze expansion of cancelled presses.
		sizet = endt
	}
	sizet /= expandDuration

	// Animate only ended presses, and presses that are fading in.
	if !c.End.IsZero() || sizet <= 1.0 {
		op.InvalidateOp{}.Add(gtx.Ops)
	}

	if sizet > 1.0 {
		sizet = 1.0
	}

	if alphat > .5 {
		// Start fadeout after half the animation.
		alphat = 1.0 - alphat
	}
	// Twice the speed to attain fully faded in at 0.5.
	t2 := alphat * 2
	// Beziér ease-in curve.
	alphaBezier := t2 * t2 * (3.0 - 2.0*t2)
	sizeBezier := sizet * sizet * (3.0 - 2.0*sizet)
	size := float32(gtx.Constraints.Min.X)
	if h := float32(gtx.Constraints.Min.Y); h > size {
		size = h
	}
	// Cover the entire constraints min rectangle.
	size *= 2 * float32(math.Sqrt(2))
	// Apply curve values to size and color.
	size *= sizeBezier
	alpha := 0.7 * alphaBezier
	const col = 0.8
	ba, bc := byte(alpha*0xff), byte(alpha*col*0xff)
	defer op.Push(gtx.Ops).Pop()
	ink := paint.ColorOp{Color: color.RGBA{A: ba, R: bc, G: bc, B: bc}}
	ink.Add(gtx.Ops)
	rr := size * .5
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(gtx.Ops)
	clip.Rect{
		Rect: f32.Rectangle{Max: f32.Point{
			X: float32(size),
			Y: float32(size),
		}},
		NE: rr, NW: rr, SE: rr, SW: rr,
	}.Op(gtx.Ops).Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
}

// mulAlpha scales all color components by alpha/255.
func mulAlpha(c color.RGBA, alpha uint8) color.RGBA {
	a := uint16(alpha)
	return color.RGBA{
		A: uint8(uint16(c.A) * a / 255),
		R: uint8(uint16(c.R) * a / 255),
		G: uint8(uint16(c.G) * a / 255),
		B: uint8(uint16(c.B) * a / 255),
	}
}
