package ui

import (
	"image"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/raedahgroup/dcrseedgen/ui/pages"
	"github.com/raedahgroup/dcrseedgen/ui/theme"
)

const (
	appName      = "Dcrseedgen"
	windowHeight = 500
	windowWidth  = 800
)

type Page interface {
	BeforeRender()
	Render(*layout.Context)
}

type Window struct {
	window          *app.Window
	theme           *theme.Theme
	pages           map[string]Page
	currentPage     string
	isRenderingPage bool
	navTabs         *theme.Tabs
}

func NewWindow(decredIcons map[string]image.Image) *Window {
	win := new(Window)
	win.currentPage = pages.SeedPageID
	win.isRenderingPage = false
	win.window = app.NewWindow(
		app.Title(appName),
	)
	win.theme = theme.New()
	win.registerPages(decredIcons)

	return win
}

func (win *Window) registerPages(decredIcons map[string]image.Image) {
	navTabs := theme.NewTabs()
	navTabs.SetTabs([]theme.TabItem{
		{
			Label: win.theme.Body1("Generate Seed"),
			Icon:  decredIcons["receive"],
		},
		{
			Label: win.theme.Body1("Generate Address"),
			Icon:  decredIcons["receive"],
		},
	})
	win.navTabs = navTabs

	win.pages = map[string]Page{
		pages.SeedPageID:    pages.NewSeedPage(win.theme),
		pages.AddressPageID: pages.NewAddressPage(win.theme),
	}
}

func (win *Window) Loop() {
	gtx := layout.NewContext(win.window.Queue())

	for {
		e := <-win.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			gtx.Reset(e.Config, e.Size)
			win.drawWindow(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func (win *Window) drawWindow(gtx *layout.Context) {
	theme.ToMax(gtx)
	theme.Fill(gtx, win.theme.Color.Background)

	navs := []string{pages.SeedPageID, pages.AddressPageID}

	win.navTabs.Separator = true
	if win.navTabs.Changed() {
		win.currentPage = navs[win.navTabs.Selected]
		win.isRenderingPage = false
	}

	if page, ok := win.pages[win.currentPage]; ok {
		if !win.isRenderingPage {
			page.BeforeRender()
		}

		win.navTabs.Layout(gtx, func() {
			layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
				page.Render(gtx)
				win.isRenderingPage = true
			})
		})
	}
}
