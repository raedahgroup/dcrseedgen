package handlers

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/widgets"
)

type (
	AddressHandler struct {
		networks  []string
		selected  string
		addresses map[string]string
		container *widget.Box
	}
)

func (h *AddressHandler) BeforeRender() {
	h.container = widget.NewVBox()
	h.networks = []string{
		"Mainnet",
		"Testnet3",
		"Simnet",
		"Regnet",
	}
	h.addresses = make(map[string]string)
}

func (h *AddressHandler) Render() fyne.CanvasObject {
	h.container.Children = []fyne.CanvasObject{
		widgets.NewVSpacer(20),
		widget.NewHBox(
			widgets.NewHSpacer(15),
			h.renderForm(),
		),
		widget.NewHBox(
			h.renderTable(),
		),
	}
	widget.Refresh(h.container)

	return h.container
}

func (h *AddressHandler) renderForm() fyne.CanvasObject {
	numberToGenerateInput := widget.NewEntry()
	numberToGenerateInput.SetPlaceHolder("Number of adresses")

	return widget.NewHBox(
		widgets.Text("Network:"),
		widget.NewSelect(h.networks, h.setSelected),
		widgets.NewHSpacer(10),
		widgets.Text("Qty:"),
		numberToGenerateInput,
		widget.NewButton("Generate", h.generateAddressesAndPrivateKeys),
	)
}

func (h *AddressHandler) setSelected(selected string) {
	h.selected = selected
}

func (h *AddressHandler) renderTable() fyne.CanvasObject {
	table := widgets.NewTable()

	if len(h.addresses) > 0 {
		table.AddRowWithTextCells(
			widgets.BoldText("Address"),
			widgets.BoldText("Private Key"),
		)
	}

	for i, v := range h.addresses {
		table.AddRowWithTextCells(
			widgets.Text(i),
			widgets.Text(v),
		)
	}

	return widget.NewHBox(
		table.CondensedTable(),
	)
}

func (h *AddressHandler) generateAddressesAndPrivateKeys() {
	address, key, _ := helper.GenerateAddressAndPrivateKey(h.selected)
	h.addresses[address] = key

	if h.container.Children != nil {
		h.container.Children = []fyne.CanvasObject{}
	}

	h.Render()
}
