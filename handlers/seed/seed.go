package seed

import (
	"fmt"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/widgets"
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
		grid                *fyne.Container
		isShowingVerifyPage bool
		isGenerating        bool
		container           *widget.Box
		masterWindow        fyne.Window
	}
)

const (
	seedSize          = 32 // affects number of words // noOfWords = (seedSize+1)
	noColumns         = 5
	noRows            = 7
	horizontalSpacing = 10
)

func (h *SeedHandler) BeforeRender(masterWindow fyne.Window) {
	h.masterWindow = masterWindow

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
	h.isGenerating = false
}

func (h *SeedHandler) buildColumns(done chan bool) {
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

	h.grid = fyne.NewContainerWithLayout(layout.NewGridLayout(noColumns))
	counter := 1
	for i := range h.columns {
		c := widget.NewVBox()
		for k := range h.columns[i].words {
			c.Append(widget.NewLabel(fmt.Sprintf("%d. %s", counter, h.columns[i].words[k])))
			counter++
		}
		h.grid.AddObject(c)
	}

	h.isGenerating = false
	done <- true
}

func (h *SeedHandler) Render() fyne.CanvasObject {
	if h.isShowingVerifyPage {
		h.renderVerifyPage()
	} else {
		h.renderSeedPage()
	}
	return h.container
}

func (h *SeedHandler) renderSeedPage() {
	var blocks *widget.Box
	defer func() {
		h.container.Children = []fyne.CanvasObject{
			blocks,
		}
	}()

	if len(h.columns) == 0 && !h.isGenerating {
		blocks = widget.NewHBox(widgets.LoadingText("Generating seed"))
		done := make(chan bool, 1)
		go h.buildColumns(done)
		<-done
		close(done)
		h.refreshPage()
		return
	}

	blocks = widget.NewHBox(
		widgets.NewHSpacer(horizontalSpacing),
		widget.NewVBox(
			widgets.BoldText("Mnemonic Words:"),
			h.grid,
			h.getSeedBlock(),
			h.buttonsBlock(),
		),
	)
}

func (h *SeedHandler) getSeedBlock() fyne.CanvasObject {
	return widget.NewVBox(
		widgets.NewVSpacer(10),
		widgets.BoldText("Hex Seed:"),
		widgets.Text(h.seed),
	)
}

func (h *SeedHandler) buttonsBlock() fyne.CanvasObject {
	return widget.NewHBox(
		widget.NewButton("Verify", h.goToVerify),
		widget.NewButton("Regenerate", h.regenerate),
	)
}

func (h *SeedHandler) goToVerify() {
	h.isShowingVerifyPage = true
	h.refreshPage()
}
func (h *SeedHandler) regenerate() {
	h.reset()
	done := make(chan bool, 1)
	go h.buildColumns(done)
	<-done
	h.refreshPage()
}
func (h *SeedHandler) refreshPage() {
	if h.container.Children != nil {
		h.container.Children = []fyne.CanvasObject{}
	}
	h.Render()
	widget.Refresh(h.container)
}
