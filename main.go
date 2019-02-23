package main

import (
	"log"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type App struct {
	currentPage  string
	pageChanged  bool
	masterWindow nucular.MasterWindow
	pages        map[string]page
}

const (
	appName  = "DCR Seed Generator"
	homePage = "seed"

	navPaneWidth            = 220
	contentPaneXOffset      = 25
	contentPaneWidthPadding = 55
)

func main() {
	app := &App{
		pageChanged: true,
		currentPage: homePage,
	}

	// register pages
	pages := getPages()
	app.pages = make(map[string]page, len(pages))
	for _, page := range pages {
		app.pages[page.name] = page
	}

	// load logo once
	err := helper.LoadLogo()
	if err != nil {
		log.Fatal(err)
	}

	window := nucular.NewMasterWindow(nucular.WindowNoScrollbar, appName, app.render)
	if err := helper.InitStyle(window); err != nil {
		log.Fatal(err)
	}

	app.masterWindow = window
	window.Main()
}

func (app *App) changePage(page string) {
	app.currentPage = page
	app.pageChanged = true
	app.masterWindow.Changed()
}

func (app *App) render(window *nucular.Window) {
	area := window.Row(0).SpaceBegin(2)

	// render nav
	app.renderNavPane(area.H, window)

	// render content
	app.renderContentPane(area, window)

}

func (app *App) renderNavPane(height int, window *nucular.Window) {
	// create navigation pane
	navPane := rect.Rect{
		X: 0,
		Y: 0,
		W: navPaneWidth,
		H: height,
	}
	window.LayoutSpacePushScaled(navPane)

	helper.StyleNav(app.masterWindow)
	if navWindow := helper.NewWindow("Navigation Window", window, 0); navWindow != nil {
		helper.DrawPageHeader(navWindow.Window)

		navWindow.Row(40).Dynamic(1)
		for _, page := range getPages() {
			if navWindow.Button(label.TA(page.label, "LC"), false) {
				app.changePage(page.name)
			}
		}
		navWindow.End()
	}
}

func (app *App) renderContentPane(area rect.Rect, window *nucular.Window) {
	// create content pane
	contentPane := rect.Rect{
		X: navPaneWidth - contentPaneXOffset,
		Y: 0,
		W: area.W - navPaneWidth,
		H: area.H,
	}

	helper.StylePage(app.masterWindow)
	window.LayoutSpacePushScaled(contentPane)
	if app.currentPage == "" {
		app.changePage(homePage)
		return
	}

	currentPage := app.pages[app.currentPage]

	if app.pageChanged {
		currentPage.handler.BeforeRender()
		app.pageChanged = false
	}

	currentPage.handler.Render(window)
}
