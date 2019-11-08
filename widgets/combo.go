package widgets

import (
	"image"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrseedgen/helper"
)

type (
	comboItem struct {
		text      string
		isTrigger bool
		clicker   helper.Clicker
	}

	Combo struct {
		items []comboItem
		theme *helper.Theme

		selectedIndex int

		isOpen bool
	}
)

func NewCombo(items []string) *Combo {
	c := &Combo{
		isOpen: false,
		items:  make([]comboItem, len(items)+1),
	}

	if len(items) > 0 {
		// item at zeroeth index is the trigger
		c.items[0] = comboItem{
			text:      items[0],
			isTrigger: true,
			clicker:   helper.NewClicker(),
		}

		for i := range items {
			c.items[i+1] = comboItem{
				text:      items[i],
				isTrigger: true,
				clicker:   helper.NewClicker(),
			}
		}
	}

	return c
}

func (c *Combo) Layout(ctx *layout.Context, theme *helper.Theme) {
	c.theme = theme
	stack := layout.Stack{}
	children := make([]layout.StackChild, len(c.items))

	currentPositionTop := float32(0)
	for i := range c.items {
		if !c.isOpen && i != 0 {
			break
		}

		for c.items[i].clicker.Clicked(ctx) {
			if i != 0 {
				c.setSelected(i)
			}
			c.isOpen = !c.isOpen
		}

		insetTop := currentPositionTop
		if i == 0 {
			insetTop = 0
		}

		children[i] = stack.Rigid(ctx, func() {
			inset := layout.Inset{
				Top: unit.Dp(insetTop),
			}

			inset.Layout(ctx, func() {
				c.drawItem(ctx, &c.items[i])
			})
		})
		currentPositionTop += float32(25)
	}

	stack.Layout(ctx, children...)

}

func (c *Combo) setSelected(itemIndex int) {
	c.selectedIndex = itemIndex
	c.items[0].text = c.items[itemIndex].text
}

func (c *Combo) GetSelected() string {
	return c.items[c.selectedIndex].text
}

func (c *Combo) drawItem(ctx *layout.Context, item *comboItem) {
	col := helper.DecredDarkBlueColor
	bgcol := helper.GrayColor

	st := layout.Stack{Alignment: layout.Center}

	font := text.Font{
		Size: unit.Dp(c.theme.TextSize.V),
	}

	//hmin := ctx.Constraints.Width.Max
	vmin := ctx.Constraints.Height.Min
	lbl := st.Rigid(ctx, func() {
		ctx.Constraints.Width.Min = 120
		ctx.Constraints.Height.Min = vmin
		layout.Align(layout.Start).Layout(ctx, func() {
			layout.UniformInset(unit.Dp(8)).Layout(ctx, func() {
				paint.ColorOp{Color: col}.Add(ctx.Ops)
				widget.Label{}.Layout(ctx, c.theme.Shaper, font, item.text)
			})
		})
		pointer.RectAreaOp{Rect: image.Rectangle{Max: ctx.Dimensions.Size}}.Add(ctx.Ops)
		item.clicker.Register(ctx)
	})
	bg := st.Expand(ctx, func() {
		rr := float32(ctx.Px(unit.Dp(0)))
		rrect(ctx.Ops,
			float32(ctx.Constraints.Width.Min),
			float32(ctx.Constraints.Height.Min),
			rr, rr, rr, rr,
		)
		fill(ctx, bgcol)
		for _, c := range item.clicker.History() {
			drawInk(ctx, c)
		}
	})

	st.Layout(ctx, bg, lbl)
}
