package helper

import (
	"time"

	"gioui.org/gesture"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
)

type (
	Clicker struct {
		click   gesture.Click
		clicks  int
		history []widget.Click
	}
)

func NewClicker() Clicker {
	return Clicker{
		history: []widget.Click{},
	}
}

func (c *Clicker) Clicked(ctx *layout.Context) bool {
	for _, e := range c.click.Events(ctx) {
		switch e.Type {
		case gesture.TypeClick:
			c.clicks++
		case gesture.TypePress:
			c.history = append(c.history, widget.Click{
				Position: e.Position,
				Time:     ctx.Now(),
			})
		}
	}
	if c.clicks > 0 {
		c.clicks--
		if c.clicks > 0 {
			// Ensure timely delivery of remaining clicks.
			op.InvalidateOp{}.Add(ctx.Ops)
		}
		return true
	}
	return false
}

func (c *Clicker) Active() bool {
	return c.click.Active()
}

func (c *Clicker) History() []widget.Click {
	return c.history
}

func (c *Clicker) Register(ctx *layout.Context) {
	c.click.Add(ctx.Ops)
	if !c.Active() {
		c.clicks = 0
	}

	for len(c.history) > 0 {
		h := c.history[0]
		if ctx.Now().Sub(h.Time) < 1*time.Second {
			break
		}
		copy(c.history, c.history[1:])
		c.history = c.history[:len(c.history)-1]
	}
}
