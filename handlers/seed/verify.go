package seed

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrseedgen/helper"
)

func (h *SeedHandler) renderVerifyPage() {
	h.container.Append(helper.BoldText("Verify:"))
	h.container.Append(h.renderVerifyInputs())
}

func (h *SeedHandler) renderVerifyInputs() fyne.CanvasObject {
	grid := fyne.NewContainerWithLayout(
		layout.NewGridLayout(noColumns),
	)
	counter := 1
	for i := range h.columns {
		c := widget.NewVBox()
		for k := range h.columns[i].words {
			c.Append(h.columns[i].inputs[k])
			counter++
		}
		grid.AddObject(c)
	}
	return grid
}

func (h *SeedHandler) renderVerifyButtons() fyne.CanvasObject {
	return widget.NewHBox(
		helper.NewHSpacer(5),
		widget.NewButton("Verify", h.verifySeed),
		widget.NewButton("BackRegenerate", h.goToSeedPage),
	)
}

func (h *SeedHandler) verifySeed() {

}

func (h *SeedHandler) goToSeedPage() {

}
