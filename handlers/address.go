package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/widgets"
)

type (
	keyPair struct {
		address    string
		privateKey string
	}

	components struct {
		selectComponent            *widget.Select
		numberOfAddressesComponent *widget.Entry
	}

	AddressHandler struct {
		networks     []string
		selected     string
		addresses    []keyPair
		container    *widget.Box
		masterWindow fyne.Window
		components
	}
)

func (h *AddressHandler) BeforeRender(masterWindow fyne.Window) {
	h.masterWindow = masterWindow
	h.container = widget.NewVBox()
	h.networks = []string{
		"Mainnet",
		"Testnet3",
		"Simnet",
		"Regnet",
	}

	h.addresses = []keyPair{}
	if h.selectComponent == nil {
		h.selectComponent = widget.NewSelect(h.networks, h.setSelected)
		h.selectComponent.Selected = h.networks[0]
		widget.Refresh(h.selectComponent)
	}

	if h.numberOfAddressesComponent == nil {
		h.numberOfAddressesComponent = widget.NewEntry()
		h.numberOfAddressesComponent.SetText("1")
	}
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

	if len(h.addresses) > 0 {
		h.container.Append(
			widget.NewHBox(
				widget.NewButton("Export as csv", h.exportAsCSV),
			),
		)
	}

	widget.Refresh(h.container)

	return h.container
}

func (h *AddressHandler) renderForm() fyne.CanvasObject {
	return widget.NewHBox(
		widgets.Text("Network:"),
		h.selectComponent,
		widgets.NewHSpacer(10),
		widgets.Text("How many addresses?:"),
		h.numberOfAddressesComponent,
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

	for _, v := range h.addresses {
		table.AddRowWithTextCells(
			widgets.Text(v.address),
			widgets.Text(v.privateKey),
		)
	}

	return widget.NewHBox(
		table.CondensedTable(),
	)
}

func (h *AddressHandler) exportAsCSV() {
	data := make([][]string, len(h.addresses))
	for i, v := range h.addresses {
		data[i] = []string{v.address, v.privateKey}
	}

	filename, err := helper.GenerateCSV(data)
	if err != nil {
		widgets.ErrorDialog(fmt.Errorf("error exporting csv: %s", err.Error()), h.masterWindow)
		return
	}

	widgets.InfoDialog(fmt.Sprintf("Succesfully exported file to %s", filename), h.masterWindow)
}

func (h *AddressHandler) generateAddressesAndPrivateKeys() {
	// first clear addresses
	h.addresses = []keyPair{}
	numberOfAddresses, err := h.getNumberOfAddresses()
	if err != nil {
		widgets.ErrorDialog(err, h.masterWindow)
		return
	}

	selectedNetwork, err := h.getSelectedNetwork()
	if err != nil {
		widgets.ErrorDialog(err, h.masterWindow)
		return
	}

	h.addresses = make([]keyPair, numberOfAddresses)

	for i := 0; i < numberOfAddresses; i++ {
		k := keyPair{}
		k.address, k.privateKey, _ = helper.GenerateAddressAndPrivateKey(selectedNetwork)
		k.address += "       "
		h.addresses[i] = k
	}

	if h.container.Children != nil {
		h.container.Children = []fyne.CanvasObject{}
	}
	h.Render()
}

func (h *AddressHandler) getSelectedNetwork() (string, error) {
	selectedNetwork := h.components.selectComponent.Selected
	if selectedNetwork == "" {
		return "", errors.New("Please select a network")
	}

	return selectedNetwork, nil
}

func (h *AddressHandler) getNumberOfAddresses() (int, error) {
	numStr := h.numberOfAddressesComponent.Text
	if numStr == "" {
		return 0, errors.New("Please specify the number of addresses to generate")
	}

	numberOfAddresses, err := strconv.Atoi(h.numberOfAddressesComponent.Text)
	if err != nil {
		return 0, errors.New("Number of addresses to generate must be a valid number")
	}

	if numberOfAddresses == 0 {
		return 0, errors.New("Number of addresses must be greater than 0")
	}

	return numberOfAddresses, nil
}
