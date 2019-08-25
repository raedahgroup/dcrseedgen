package main

import (
	"log"

	"github.com/aarzilli/nucular"
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

	window := nucular.NewMasterWindow(nucular.WindowNoScrollbar|nucular.WindowContextualReplace, appName, app.render)
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
	currentPage := app.pages[app.currentPage]

	if app.pageChanged {
		currentPage.handler.BeforeRender()
		app.pageChanged = false
	}

	helper.DrawPageHeader(window)
	app.renderNavButtons(window)
	currentPage.handler.Render(window)
}

func (app *App) renderNavButtons(window *nucular.Window) {
	window.Row(helper.ButtonHeight + 10).Dynamic(1)
	if group := window.GroupBegin("nav-window", 0); group != nil {
		group.Row(helper.ButtonHeight).Ratio(0.18, 0.23)
		helper.StyleNavButton(window)
		if group.ButtonText("Generate Seed") && app.currentPage != "seed" {
			app.changePage("seed")
		}

		if group.ButtonText("Generate Address") && app.currentPage != "address" {
			app.changePage("address")
		}
		helper.ResetButtonStyle(window)
		group.GroupEnd()
	}
}
