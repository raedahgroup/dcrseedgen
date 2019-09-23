package main

import (
	"context"
	"image"
	"log"

	"gioui.org/layout"
	"github.com/raedahgroup/dcrseedgen/helper"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
)

type App struct {
	window *app.Window
	ctx    context.Context

	pageChanged bool
	currentPage string
	pages       []page

	theme *helper.Theme
}

const (
	appName      = "DCR Seed Generator"
	windowWidth  = 500
	windowHeight = 350
)

func main() {
	a := &App{
		theme: helper.NewTheme(),
	}
	a.setHandlers()

	go func() {
		a.window = app.NewWindow(
			app.Size(
				unit.Dp(windowWidth),
				unit.Dp(windowHeight),
			),
			app.Title(appName),
		)

		if err := a.startRenderLoop(); err != nil {
			log.Fatal(err)
		}
	}()

	app.Main()
}

func (a *App) setHandlers() {
	pages := getPages(a.theme)
	a.pages = make([]page, len(pages))

	for index, page := range pages {
		a.pages[index] = page
	}

	if len(a.pages) > 0 {
		a.changePage(a.pages[0].name)
	}
}

func (a *App) changePage(pageName string) {
	if a.currentPage == pageName {
		return
	}

	a.pageChanged = true
	a.currentPage = pageName

	if a.window != nil {
		a.refreshWindow()
	}
}

func (a *App) startRenderLoop() error {
	ctx := &layout.Context{
		Queue: a.window.Queue(),
	}

	for {
		e := <-a.window.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ctx.Reset(e.Config, e.Size)
			a.drawWindowContents(ctx)
			e.Frame(ctx.Ops)

		}
	}
}

func (a *App) drawWindowContents(ctx *layout.Context) {
	stack := layout.Stack{}

	navSection := stack.Rigid(ctx, func() {
		a.drawNavSection(ctx)
	})

	contentSection := stack.Rigid(ctx, func() {
		a.drawContentSection(ctx)
	})

	stack.Layout(ctx, navSection, contentSection)
}

func (a *App) drawNavSection(ctx *layout.Context) {
	navAreaBounds := image.Point{
		X: windowWidth * 2,
		Y: 53,
	}
	helper.PaintArea(ctx, helper.GrayColor, navAreaBounds)
	inset := layout.Inset{
		Top:  unit.Sp(0),
		Left: unit.Sp(0),
	}

	inset.Layout(ctx, func() {
		flex := layout.Flex{
			Axis: layout.Horizontal,
		}
		children := make([]layout.FlexChild, len(a.pages))
		inset := layout.UniformInset(unit.Dp(0))
		for index, page := range a.pages {
			children[index] = flex.Rigid(ctx, func() {
				inset.Layout(ctx, func() {
					for page.button.Clicked(ctx) {
						a.changePage(page.name)
					}
					btn := a.theme.Button(page.navLabel)
					btn.Color = helper.DecredDarkBlueColor
					if a.currentPage == page.name {
						btn.Background = helper.WhiteColor
					} else {
						btn.Background = helper.GrayColor
					}
					btn.Layout(ctx, page.button)
				})
			})
		}
		flex.Layout(ctx, children...)
	})
}

func (a *App) drawContentSection(ctx *layout.Context) {
	var page page
	for i := range a.pages {
		if a.pages[i].name == a.currentPage {
			page = a.pages[i]
			break
		}
	}

	if a.pageChanged {
		page.handler.BeforeRender()
		a.pageChanged = false
	}
	stack := layout.Stack{}
	inset := layout.Inset{
		Top:   unit.Dp(48),
		Left:  unit.Dp(15),
		Right: unit.Dp(15),
	}

	inset.Layout(ctx, func() {
		page.handler.Render(ctx, a.refreshWindow)
	})
	stack.Layout(ctx)
}

func (a *App) refreshWindow() {
	a.window.Invalidate()
}
