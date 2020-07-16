package theme

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type Tab struct {
	btn     widget.Clickable
	ID      string
	Title   string
	Content layout.Widget
}

type Tabs struct {
	theme    *Theme
	list     layout.List
	tabs     []Tab
	selected int
	slider   Slider
	changed  bool
}

func (t *Theme) NewTabs() *Tabs {
	return &Tabs{
		theme:   t,
		list:    layout.List{Axis: layout.Vertical},
		changed: false,
	}
}

func (t *Tabs) AddItems(items []Tab) {
	for i := range items {
		items[i].btn = widget.Clickable{}
	}
	t.tabs = items
}

func (t *Tabs) SelectedID() string {
	return t.tabs[t.selected].ID
}

func (t *Tabs) Changed() bool {
	return t.changed
}

func (t *Tabs) Layout(gtx layout.Context) layout.Dimensions {
	t.changed = false

	return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return t.list.Layout(gtx, len(t.tabs), func(gtx layout.Context, tabID int) layout.Dimensions {
				current := &t.tabs[tabID]
				if current.btn.Clicked() {
					if tabID != t.selected {
						t.changed = true
					}
					t.selected = tabID
				}

				gtx.Constraints.Max.X = 150
				return material.Clickable(gtx, &current.btn, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Constraints.Max.X
					return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							if tabID == t.selected {
								return line(gtx, 3, 45, t.theme.Color.Primary)(gtx)
							}
							return layout.Dimensions{}
						}),
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
								return material.Body1(t.theme.Theme, current.Title).Layout(gtx)
							})
						}),
					)
				})
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return line(gtx, 2, gtx.Constraints.Max.Y, rgb(0xcccccc))(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			gtx.Constraints.Min = gtx.Constraints.Max
			return t.tabs[t.selected].Content(gtx)
		}),
	)
}

// line returns a rectangle using a defined width, height and color.
func line(gtx layout.Context, width, height int, col color.RGBA) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		paint.ColorOp{Color: col}.Add(gtx.Ops)
		paint.PaintOp{Rect: f32.Rectangle{
			Max: f32.Point{
				X: float32(width),
				Y: float32(height),
			},
		}}.Add(gtx.Ops)
		return layout.Dimensions{
			Size: image.Point{X: width, Y: height},
		}
	}
}
