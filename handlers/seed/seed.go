package seed

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/raedahgroup/dcrseedgen/widgets"

	"strconv"
)

func (h *SeedHandler) drawSeedSubPage(ctx *layout.Context, refreshWindowFunc func()) {
	stack := layout.Stack{}

	mnemonicHeadingSection := stack.Rigid(ctx, func() {
		h.theme.H5("Mnemonic Words:").Layout(ctx)
	})

	wordsSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(3),
		}

		inset.Layout(ctx, func() {
			h.renderColumns(ctx)
		})
	})

	hexHeadingSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(158),
		}

		inset.Layout(ctx, func() {
			h.theme.H5("Hex Seed:").Layout(ctx)
		})
	})

	hexSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(180),
		}

		inset.Layout(ctx, func() {
			h.theme.Body1(h.seed.seedStr).Layout(ctx)
		})
	})

	buttonsSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(205),
		}

		inset.Layout(ctx, func() {
			h.renderSeedPageButtons(ctx, refreshWindowFunc)
		})
	})

	stack.Layout(ctx, mnemonicHeadingSection, wordsSection, hexHeadingSection, hexSection, buttonsSection)
}

func (h *SeedHandler) renderColumns(ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	inset := layout.Inset{
		Right: unit.Dp(30),
	}
	children := make([]layout.FlexChild, numberOfColumns)

	currentItem := 1
	for i := range h.seed.columns {
		children[i] = flex.Rigid(ctx, func() {
			inset.Layout(ctx, func() {
				h.renderColumn(h.seed.columns[i].words, &currentItem, ctx)
			})
		})
	}
	flex.Layout(ctx, children...)
}

func (h *SeedHandler) renderColumn(words []string, currentItem *int, ctx *layout.Context) {
	stack := layout.Stack{}

	nextLabelTopPosition := float32(18)
	labelHeight := 18

	child := stack.Rigid(ctx, func() {
		for _, word := range words {
			inset := layout.Inset{
				Top: unit.Dp(nextLabelTopPosition),
			}

			inset.Layout(ctx, func() {
				h.theme.Body1(strconv.Itoa(*currentItem) + ". " + word).Layout(ctx)
			})

			nextLabelTopPosition += float32(labelHeight)
			*currentItem++
		}
	})
	stack.Layout(ctx, child)
}

func (h *SeedHandler) renderSeedPageButtons(ctx *layout.Context, refreshWindowFunc func()) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	regenerateButtonSection := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Left: unit.Dp(0),
		}

		inset.Layout(ctx, func() {
			for h.regenerateButton.Clicked(ctx) {
				h.reset()
				refreshWindowFunc()
			}
			widgets.LayoutButton(h.regenerateButton, "Regenerate", h.theme, ctx)
		})
	})

	verifyButtonSection := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Left: unit.Dp(300),
		}

		inset.Layout(ctx, func() {
			for h.verifyButton.Clicked(ctx) {
				h.currentSubPage = verifySubPageName
				refreshWindowFunc()
			}
			widgets.LayoutButton(h.verifyButton, "Verify Words", h.theme, ctx)
		})
	})

	flex.Layout(ctx, regenerateButtonSection, verifyButtonSection)
}
