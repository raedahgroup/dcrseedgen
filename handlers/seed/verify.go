package seed

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrseedgen/widgets"
)

func (h *SeedHandler) renderVerifyPage() {
	verifyLabel := widgets.BoldText("Verify:")
	inputs := h.renderVerifyInputs()

	h.container.Children = []fyne.CanvasObject{
		widget.NewVBox(
			verifyLabel,
			inputs,
			h.renderVerifyButtons(),
		),
	}
}

func (h *SeedHandler) renderVerifyInputs() fyne.CanvasObject {
	grid := fyne.NewContainerWithLayout(layout.NewGridLayout(noColumns))
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
		widget.NewButton("Verify", h.verifySeed),
		widget.NewButton("Back", h.goToSeedPage),
	)
}

func (h *SeedHandler) verifySeed() {

}

func (h *SeedHandler) goToSeedPage() {
	h.isShowingVerifyPage = false
	h.refreshPage()
}
