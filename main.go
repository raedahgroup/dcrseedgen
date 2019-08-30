package main

import (
	"fmt"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type App struct {
}

const (
	appName = "DCR Seed Generator"
)

func main() {
	err := helper.InitFonts()
	if err != nil {
		log.Fatal(fmt.Errorf("error loading fonts: %s", err.Error()))
	}

	// create export folder in a separate goroutine
	go func() {
		err := helper.CreateExportFolder()
		if err != nil {
			log.Fatalf("error creating export folder: %s", err.Error())
		}
	}()

	app := app.New()
	app.Settings().SetTheme(helper.DefaultTheme)
	masterWindow := app.NewWindow(appName)

	navTabs := widget.NewTabContainer(getPages(masterWindow)...)
	navTabs.SetTabLocation(widget.TabLocationTop)

	masterWindow.Resize(fyne.NewSize(1070, 750))
	masterWindow.SetFixedSize(true)
	masterWindow.CenterOnScreen()
	masterWindow.SetContent(navTabs)
	masterWindow.ShowAndRun()

}

func getPages(masterWindow fyne.Window) []*widget.TabItem {
	handlers := getHandlers()
	pages := make([]*widget.TabItem, len(handlers))

	for i, v := range handlers {
		pages[i] = widget.NewTabItemWithIcon(v.label, v.icon, render(v.handler, masterWindow))
	}
	return pages
}

func render(h handler, masterWindow fyne.Window) fyne.CanvasObject {
	// call before render method to load required data and setup variables
	h.BeforeRender(masterWindow)
	container := widget.NewScrollContainer(
		h.Render(),
	)
	//container.Resize(fyne.NewSize(970, 750))
	return container
}
