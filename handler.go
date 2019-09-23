package main

import (
	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrseedgen/handlers/seed"
	"github.com/raedahgroup/dcrseedgen/handlers"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type handler interface {
	BeforeRender()
	Render(*layout.Context, func())
}

type page struct {
	name     string
	navLabel string
	button  *widget.Button
	handler  handler
}

func getPages(theme *helper.Theme) []page {
	return []page{
		{
			name:     "seedhandler",
			navLabel: "Generate Seed",
			button:   new(widget.Button),
			handler:  seed.NewSeedHandler(theme),
		},
		{
			name:     "wallethandler",
			navLabel: "Generate Address",
			button:   new(widget.Button),
			handler:  handlers.NewAddressHandler(theme),
		},
	}
}
