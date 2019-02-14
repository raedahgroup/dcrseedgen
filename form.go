package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/decred/dcrwallet/walletseed"

	"github.com/aarzilli/nucular"
)

type verifyMessage struct {
	message     string
	messageType string
}

type wordInputColumn struct {
	words  []string
	inputs []nucular.TextEditor
}

type RenderHandler struct {
	words   string
	columns []wordInputColumn
	err     error

	seed        string
	currentPage *string

	verifyMessage *verifyMessage
}

const (
	seedSize  = 32 // affects number of words // noOfWords = (seedSize+1)
	noColumns = 5
	noRows    = 7
)

func (h *RenderHandler) beforeRender(currentPage *string) {
	h.currentPage = currentPage
	h.generate()
}

func (h *RenderHandler) generate() {
	// generate mnemonic words
	words, seed, err := generateWords()
	if err != nil {
		h.err = err
		return
	}
	h.words = words
	h.seed = seed

	h.buildColumns()
}

func (h *RenderHandler) buildColumns() {
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

func generateWords() (string, string, error) {
	seed, err := walletseed.GenerateRandomSeed(seedSize)
	if err != nil {
		return "", "", err
	}

	return walletseed.EncodeMnemonic(seed), hex.EncodeToString(seed), nil
}

func (h *RenderHandler) renderHome(window *nucular.Window) {
	if h.err != nil {
		window.Row(20).Dynamic(1)
		window.Label(h.err.Error(), "LC")
	}

	drawHeader(window)

	window.Row(380).Dynamic(1)
	if newWindow := window.GroupBegin("Main Content", 0); newWindow != nil {
		newWindow.Row(20).Dynamic(1)
		SetFont(newWindow, boldFont)
		newWindow.Label("Mnemonic Words:", "LC")

		newWindow.Row(187).Dynamic(1)
		if group := newWindow.GroupBegin("", 0); group != nil {
			group.Row(166).Dynamic(noColumns)
			SetFont(group, boldFont)

			currentItem := 0
			for _, column := range h.columns {
				newWordColumn(group, column.words, &currentItem)
			}
			group.GroupEnd()
		}

		if h.err != nil {
			newWindow.Row(30).Dynamic(1)
			newWindow.Label(fmt.Sprintf("error generating seed: %s", h.err.Error()), "LC")
		} else {
			newWindow.Row(1).Dynamic(1)
			newWindow.Label("", "LC")

			newWindow.Row(20).Dynamic(1)
			SetFont(newWindow, boldFont)
			newWindow.Label("Hex Seed", "LC")
			newWindow.Row(60).Dynamic(1)
			SetFont(newWindow, normalFont)
			newWindow.LabelWrap(h.seed)
		}

		SetFont(newWindow, normalFont)
		newWindow.Row(40).Ratio(0.5, 0.25, 0.25)
		newWindow.Label("", "LC")

		if newWindow.ButtonText("Verify") {
			h.verifyMessage = &verifyMessage{}
			*h.currentPage = "verify"
			newWindow.Master().Changed()
		}

		if newWindow.ButtonText("Regenerate") {
			h.generate()
			newWindow.Master().Changed()
		}
		newWindow.GroupEnd()
	}
}

func newWordColumn(window *nucular.Window, words []string, currentItem *int) {
	if group := window.GroupBegin(words[0], 0); group != nil {
		SetFont(group, normalFont)
		for _, word := range words {
			group.Row(20).Dynamic(1)
			group.Label(strconv.Itoa(*currentItem+1)+". "+word, "LC")
			*currentItem++
		}
		group.GroupEnd()
	}
}
