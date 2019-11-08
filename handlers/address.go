package handlers

import (
	"errors"
	"strconv"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/widgets"
)

type (
	AddressHandler struct {
		theme               *helper.Theme
		networkCombo        *widgets.Combo
		numberOfItemsEditor *widget.Editor
		generateButton      *widget.Button
		addresses           []string
		privateKeys         []string

		err error
	}
)

func NewAddressHandler(theme *helper.Theme) *AddressHandler {
	networks := []string{
		"Testnet3",
		"Mainnet",
		"Regnet",
	}

	numberOfItemsEditor := &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	numberOfItemsEditor.SetText("1")

	return &AddressHandler{
		theme:               theme,
		networkCombo:        widgets.NewCombo(networks),
		numberOfItemsEditor: numberOfItemsEditor,
		generateButton:      new(widget.Button),
	}
}

func (h *AddressHandler) BeforeRender() {

}

func (h *AddressHandler) Render(ctx *layout.Context, refreshWindowFunc func()) {
	stack := layout.Stack{}

	formSection := stack.Rigid(ctx, func() {
		h.renderFormSection(ctx)
	})

	errSection := stack.Rigid(ctx, func() {
		if h.err != nil {
			inset := layout.Inset{
				Top: unit.Dp(35),
			}
			inset.Layout(ctx, func() {
				label := h.theme.H5(h.err.Error())
				label.Color = helper.DangerColor
				label.Layout(ctx)
			})
		}
	})

	contentSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(50),
		}

		inset.Layout(ctx, func() {
			if len(h.addresses) > 0 && h.err == nil {
				h.renderTableHeader(ctx)

				insetTop := float32(20)
				for index, address := range h.addresses {
					h.renderPairRow(address, h.privateKeys[index], insetTop, ctx)

					insetTop += float32(20)
				}
			}
		})
	})

	stack.Layout(ctx, formSection, errSection, contentSection)
}

func (h *AddressHandler) renderTableHeader(ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	addressColumn := flex.Rigid(ctx, func() {
		h.theme.H5("Address").Layout(ctx)
	})

	keyColumn := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Left: unit.Dp(125),
		}

		inset.Layout(ctx, func() {
			h.theme.H5("Private Key").Layout(ctx)
		})
	})

	flex.Layout(ctx, addressColumn, keyColumn)
}

func (h *AddressHandler) renderPairRow(key, val string, insetTop float32, ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	addressColumn := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(insetTop),
		}
		inset.Layout(ctx, func() {
			h.theme.Body2(key).Layout(ctx)
		})
	})

	spacer := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(insetTop),
		}

		inset.Layout(ctx, func() {
			h.theme.Body2("  ").Layout(ctx)
		})
	})

	keyColumn := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(insetTop),
		}

		inset.Layout(ctx, func() {
			h.theme.Body2(val).Layout(ctx)
		})
	})

	flex.Layout(ctx, addressColumn, spacer, keyColumn)
}

func (h *AddressHandler) renderFormSection(ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	networkLabelSection := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(8),
		}
		inset.Layout(ctx, func() {
			h.theme.Body1("Network: ").Layout(ctx)
		})
	})

	networkComboSection := flex.Rigid(ctx, func() {
		h.networkCombo.Layout(ctx, h.theme)
	})

	numberOfItemsLabel := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(8),
		}
		inset.Layout(ctx, func() {
			h.theme.Body1("     Number Of Addresses : ").Layout(ctx)
		})
	})

	numberOfItemsSection := flex.Rigid(ctx, func() {
		e := h.theme.Editor("")
		e.Font = text.Font{
			Size: unit.Sp(10),
		}

		inset := layout.Inset{
			Top: unit.Dp(8),
		}
		inset.Layout(ctx, func() {
			e.Layout(ctx, h.numberOfItemsEditor)
		})
	})

	spacer := flex.Rigid(ctx, func() {
		h.theme.Body1("   ").Layout(ctx)
	})

	generateButtonSection := flex.Rigid(ctx, func() {
		for h.generateButton.Clicked(ctx) {
			h.doGenerate(ctx)
		}
		widgets.LayoutButton(h.generateButton, "Generate", h.theme, ctx)
	})

	flex.Layout(ctx, networkLabelSection, networkComboSection, numberOfItemsLabel, numberOfItemsSection, spacer, generateButtonSection)
}

func (h *AddressHandler) doGenerate(ctx *layout.Context) {
	numberOfItemsToGenerateStr := h.numberOfItemsEditor.Text()
	if numberOfItemsToGenerateStr == "" {
		h.err = errors.New("Please type in the required number of pairs")
		return
	}

	numberOfItemsToGenerate, err := strconv.Atoi(numberOfItemsToGenerateStr)
	if err != nil {
		h.err = errors.New("Invalid number")
		return
	}

	network := h.networkCombo.GetSelected()

	h.addresses = make([]string, numberOfItemsToGenerate)
	h.privateKeys = make([]string, numberOfItemsToGenerate)

	for i := 0; i < numberOfItemsToGenerate; i++ {
		privateKey, address, err := helper.GenerateAddressAndPrivateKey(network)
		if err != nil {
			h.err = err
			return
		}

		h.addresses[i] = address
		h.privateKeys[i] = privateKey
	}
}
