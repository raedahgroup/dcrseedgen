package seed

import (
	"fmt"

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
	wrong := false

	for _ = range h.columns {
		for columnIndex := range h.columns {
			for itemIndex := range h.columns[columnIndex].words {
				if h.columns[columnIndex].words[itemIndex] != h.columns[columnIndex].inputs[itemIndex].Text {
					wrong = true
				}
			}
		}
	}

	if !wrong {
		widgets.InfoDialog("Successfully verified words", h.masterWindow)
		return
	}

	widgets.ErrorDialog(fmt.Errorf("Incorrect verification words"), h.masterWindow)
}

func (h *SeedHandler) goToSeedPage() {
	h.isShowingVerifyPage = false
	h.refreshPage()
}
