package ui

import (
	"image"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrseedgen/ui/pages"
	"github.com/raedahgroup/dcrseedgen/ui/theme"
)

const (
	appName      = "Dcrseedgen"
	windowHeight = 600
	windowWidth  = 850
)

type Page interface {
	BeforeRender()
	Render(layout.Context) layout.Dimensions
}

type Window struct {
	window          *app.Window
	theme           *theme.Theme
	pages           map[string]Page
	currentPage     string
	isRenderingPage bool
	navTabs         *theme.Tabs
}

func NewWindow(decredIcons map[string]image.Image, col *text.Collection) *Window {
	win := new(Window)
	win.currentPage = pages.SeedPageID
	win.isRenderingPage = false
	win.window = app.NewWindow(
		app.Size(unit.Dp(windowWidth), unit.Dp(windowHeight)),
		app.Title(appName),
	)
	win.theme = theme.New(col)
	win.registerPages(decredIcons)

	return win
}

func (win *Window) registerPages(decredIcons map[string]image.Image) {
	win.pages = map[string]Page{
		pages.SeedPageID:    pages.NewSeedPage(win.theme),
		pages.AddressPageID: pages.NewAddressPage(win.theme),
	}

	win.navTabs = win.theme.NewTabs()
	win.navTabs.AddItems([]theme.Tab{
		{
			ID:      pages.SeedPageID,
			Title:   "Generate Seed",
			Content: win.pages[pages.SeedPageID].Render,
		},
		{
			ID:      pages.AddressPageID,
			Title:   "Generate Address",
			Content: win.pages[pages.AddressPageID].Render,
		},
	})
}

func (win *Window) Loop() {
	var ops op.Ops
	for {
		select {
		case e := <-win.window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				win.drawWindow(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func (win *Window) drawWindow(gtx layout.Context) {
	theme.ToMax(gtx)
	theme.Fill(gtx, win.theme.Color.Background)

	if win.navTabs.Changed() {
		win.currentPage = win.navTabs.SelectedID()
		win.isRenderingPage = false
	}

	if page, ok := win.pages[win.currentPage]; ok {
		if !win.isRenderingPage {
			page.BeforeRender()
		}

		win.isRenderingPage = true
		win.navTabs.Layout(gtx)
	}
}
