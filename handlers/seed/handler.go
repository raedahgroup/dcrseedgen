package seed

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget"

	"github.com/raedahgroup/dcrseedgen/helper"
)

type (
	column struct {
		words   []string
		editors []*widget.Editor
	}

	seed struct {
		seedStr string
		columns []column
	}

	verifyMessage struct {
		message     string
		messageType string
	}

	SeedHandler struct {
		currentSubPage string
		err            error
		seed           *seed

		theme         *helper.Theme
		verifyMessage *verifyMessage

		regenerateButton *widget.Button
		verifyButton     *widget.Button
		backButton       *widget.Button
		doVerifyButton   *widget.Button
	}
)

const (
	seedSubPageName   = "seed"
	verifySubPageName = "verify"

	seedSize          = 32 // affects number of words // noOfWords = (seedSize+1)
	numberOfColumns   = 5
	numberOfRows      = 7
	horizontalSpacing = 10
)

func NewSeedHandler(theme *helper.Theme) *SeedHandler {
	return &SeedHandler{
		theme:            theme,
		regenerateButton: new(widget.Button),
		verifyButton:     new(widget.Button),
		backButton:       new(widget.Button),
		doVerifyButton:   new(widget.Button),
	}
}

func (h *SeedHandler) BeforeRender() {
	h.verifyMessage = nil

	if h.currentSubPage == verifySubPageName {
		return
	}
	h.reset()
}

func (h *SeedHandler) reset() {
	h.seed = nil
	h.generateMnemonic()
}

func (h *SeedHandler) generateMnemonic() {
	words, seedStr, err := helper.GenerateMnemonicSeed(seedSize)
	if err != nil {
		h.err = err
		return
	}

	h.seed = &seed{
		seedStr: seedStr,
		columns: make([]column, numberOfColumns),
	}

	wordSlice := strings.Split(words, " ")
	currentColumn := 0

	for index, word := range wordSlice {
		h.seed.columns[currentColumn].words = append(h.seed.columns[currentColumn].words, word)
		editor := &widget.Editor{
			SingleLine: true,
			Submit:     true,
		}
		h.seed.columns[currentColumn].editors = append(h.seed.columns[currentColumn].editors, editor)

		if index > 0 && (index+1)%numberOfRows == 0 {
			currentColumn++
		}
	}
}

func (h *SeedHandler) Render(ctx *layout.Context, refreshWindowFunc func()) {
	if h.currentSubPage == verifySubPageName {
		h.drawVerifySubPage(ctx, refreshWindowFunc)
		return
	}

	h.currentSubPage = seedSubPageName
	h.drawSeedSubPage(ctx, refreshWindowFunc)
}
