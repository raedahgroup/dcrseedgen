package main

import (
	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrseedgen/handlers"
)

type Handler interface {
	BeforeRender()
	Render(*nucular.Window)
}

type page struct {
	name    string
	label   string
	handler Handler
}

func getPages() []page {
	return []page{
		{
			name:    "seed",
			label:   "Generate Seed",
			handler: &handlers.SeedGeneratorHandler{},
		},
		{
			name:    "address",
			label:   "Get Address",
			handler: &handlers.AddressGeneratorHandler{},
		},
	}
}
