package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aarzilli/nucular"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type verifyMessage struct {
	message     string
	messageType string
}

type wordInputColumn struct {
	words  []string
	inputs []nucular.TextEditor
}

type SeedGeneratorHandler struct {
	seedErr error
	words   string
	columns []wordInputColumn

	seed                string
	isShowingVerifyPage bool

	verifyMessage *verifyMessage
}

const (
	seedSize  = 32 // affects number of words // noOfWords = (seedSize+1)
	noColumns = 5
	noRows    = 7
)

func (h *SeedGeneratorHandler) BeforeRender() {
	if h.isShowingVerifyPage {
		return
	}
	h.generateSeed()
}

func (h *SeedGeneratorHandler) generateSeed() {
	h.words, h.seed, h.seedErr = helper.GenerateMnemonicSeed(seedSize)
	if h.seedErr != nil {
		return
	}

	h.buildColumns()
}

func (h *SeedGeneratorHandler) buildColumns() {
	wordSlice := strings.Split(h.words, " ")
	h.columns = make([]wordInputColumn, noColumns)

	currentColumn := 0
	for index, word := range wordSlice {
		h.columns[currentColumn].words = append(h.columns[currentColumn].words, word)
		h.columns[currentColumn].inputs = append(h.columns[currentColumn].inputs, nucular.TextEditor{})

		if index > 0 && (index+1)%noRows == 0 {
			currentColumn++
		}
	}
}

func (h *SeedGeneratorHandler) Render(window *nucular.Window) {
	if !h.isShowingVerifyPage {
		h.renderSeedPage(window)
		return
	}

	h.renderVerifyPage(window)
}

func (h *SeedGeneratorHandler) renderSeedPage(window *nucular.Window) {
	if h.seedErr != nil {
		window.Row(20).Dynamic(1)
		window.Label(h.seedErr.Error(), "LC")
		return
	}

	// draw page header
	helper.DrawPageHeader(window)

	window.Row(380).Dynamic(1)
	if w := helper.NewWindow("Seed Page Content", window, 0); w != nil {
		w.Row(20).Dynamic(1)

		// set font
		helper.UseFont(w, helper.FontBold)
		w.Label("Mnemonic Words:", "LC")

		w.Row(187).Dynamic(1)
		if colWindow := w.NewWindow("Word Columns", 0); colWindow != nil {
			colWindow.Row(166).Dynamic(noColumns)
			helper.UseFont(colWindow, helper.FontBold)

			currentItem := 0
			for _, column := range h.columns {
				newWordColumn(colWindow, column.words, &currentItem)
			}
			colWindow.End()
		}

		if h.seedErr != nil {
			w.Row(30).Dynamic(1)
			w.Label(fmt.Sprintf("error generating seed: %s", h.seedErr.Error()), "LC")
		} else {
			w.Row(1).Dynamic(1)
			w.Label("", "LC")

			w.Row(20).Dynamic(1)
			helper.UseFont(w, helper.FontBold)

			w.Label("Hex Seed", "LC")
			w.Row(60).Dynamic(1)
			helper.UseFont(w, helper.FontNormal)
			w.LabelWrap(h.seed)
		}

		helper.UseFont(w, helper.FontNormal)
		w.Row(40).Ratio(0.5, 0.25, 0.25)
		w.Label("", "LC")

		if w.ButtonText("Verify") {
			h.isShowingVerifyPage = true
			w.Master().Changed()
		}

		if w.ButtonText("Regenerate") {
			h.generateSeed()
			w.Master().Changed()
		}
		w.End()
	}
}

func newWordColumn(window *helper.Window, words []string, currentItem *int) {
	if w := window.NewWindow(words[0], 0); w != nil {
		helper.UseFont(w, helper.FontNormal)
		for _, word := range words {
			w.Row(20).Dynamic(1)
			w.Label(strconv.Itoa(*currentItem+1)+". "+word, "LC")
			*currentItem++
		}
		w.End()
	}
}

func (h *SeedGeneratorHandler) renderVerifyPage(window *nucular.Window) {

}
