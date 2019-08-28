package seed

import (
	"fmt"
	"strings"

	"fyne.io/fyne/layout"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type (
	column struct {
		words  []string
		inputs []*widget.Entry
	}
	SeedHandler struct {
		err                 error
		seed                string
		columns             []column
		isShowingVerifyPage bool
		container           *widget.Box
	}
)

const (
	seedSize  = 32 // affects number of words // noOfWords = (seedSize+1)
	noColumns = 5
	noRows    = 7
)

func (h *SeedHandler) BeforeRender() {
	if h.isShowingVerifyPage {
		return
	}

	h.container = widget.NewVBox()
	h.reset()
}

func (h *SeedHandler) reset() {
	h.err = nil
	h.seed = ""
	h.columns = nil
	h.buildColumns()
}

func (h *SeedHandler) buildColumns() {
	words, seed, err := helper.GenerateMnemonicSeed(seedSize)
	if err != nil {
		h.err = err
		return
	}

	h.columns = make([]column, noColumns)
	h.seed = seed

	wordSlice := strings.Split(words, " ")
	currentColumn := 0
	for index, word := range wordSlice {
		h.columns[currentColumn].words = append(h.columns[currentColumn].words, word)
		h.columns[currentColumn].inputs = append(h.columns[currentColumn].inputs, widget.NewEntry())

		if index > 0 && (index+1)%noRows == 0 {
			currentColumn++
		}
	}
}

func (h *SeedHandler) Render() fyne.CanvasObject {
	if h.isShowingVerifyPage {
		h.renderVerifyPage()
		return h.container
	}

	h.container.Append(helper.BoldText("Mnemonic Words:"))
	h.container.Append(h.renderWords())
	h.container.Append(helper.NewVSpacer(40))
	h.container.Append(h.renderSeed())
	h.container.Append(helper.NewVSpacer(10))
	h.container.Append(h.renderButtons())
	return h.container
}

func (h *SeedHandler) renderWords() fyne.CanvasObject {
	grid := fyne.NewContainerWithLayout(
		layout.NewGridLayout(noColumns),
	)
	counter := 1
	for i := range h.columns {
		c := widget.NewVBox()
		for k := range h.columns[i].words {
			c.Append(
				widget.NewLabel(fmt.Sprintf("%d. %s", counter, h.columns[i].words[k])),
			)
			counter++
		}
		grid.AddObject(c)
	}
	return grid
}

func (h *SeedHandler) renderSeed() fyne.CanvasObject {
	return widget.NewVBox(
		helper.NewVSpacer(1),
		helper.BoldText("Hex Seed:"),
		helper.Text(h.seed),
	)
}

func (h *SeedHandler) renderButtons() fyne.CanvasObject {
	return widget.NewHBox(
		helper.NewHSpacer(5),
		widget.NewButton("Verify", h.goToVerify),
		widget.NewButton("Regenerate", h.regenerate),
	)
}

func (h *SeedHandler) goToVerify() {
	h.isShowingVerifyPage = true
	h.refresh(false)
}

func (h *SeedHandler) regenerate() {
	h.refresh(true)
}

func (h *SeedHandler) refresh(reset bool) {
	if reset {
		h.reset()
	}

	if h.container.Children != nil {
		h.container.Children = []fyne.CanvasObject{}
	}
	h.Render()
	widget.Refresh(h.container)
}
