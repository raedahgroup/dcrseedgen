package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/raedahgroup/dcrseedgen/handlers"
	"github.com/raedahgroup/dcrseedgen/handlers/seed"
)

type handler interface {
	BeforeRender(fyne.Window)
	Render() fyne.CanvasObject
}

type page struct {
	label   string
	icon    fyne.Resource
	handler handler
}

func getHandlers() []page {
	return []page{
		{
			label:   "Generate Seed",
			icon:    theme.ConfirmIcon(),
			handler: &seed.SeedHandler{},
		},
		{
			label:   "Generate Addresses",
			icon:    theme.HomeIcon(),
			handler: &handlers.AddressHandler{},
		},
	}
}
