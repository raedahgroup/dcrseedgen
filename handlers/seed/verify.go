package seed

import (
	"gioui.org/layout"
	"gioui.org/unit"

	"gioui.org/text"
	"gioui.org/widget"

	"strconv"

	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/widgets"
)

func (h *SeedHandler) drawVerifySubPage(ctx *layout.Context, refreshWindowFunc func()) {
	stack := layout.Stack{}

	verifyHeadingSection := stack.Rigid(ctx, func() {
		h.theme.H5("Verify Mnemonic Words:").Layout(ctx)
	})

	inputsSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(10),
		}

		inset.Layout(ctx, func() {
			h.renderInputColumns(ctx)
		})
	})

	messageSection := stack.Rigid(ctx, func() {
		if h.verifyMessage != nil {
			inset := layout.Inset{
				Top: unit.Dp(165),
			}

			inset.Layout(ctx, func() {
				label := h.theme.H5(h.verifyMessage.message)
				if h.verifyMessage.messageType == "success" {
					label.Color = helper.SuccessColor
				} else {
					label.Color = helper.DangerColor
				}

				label.Layout(ctx)
			})
		}
	})

	buttonsSection := stack.Rigid(ctx, func() {
		inset := layout.Inset{
			Top: unit.Dp(190),
		}

		inset.Layout(ctx, func() {
			h.renderVerifyPageButtons(ctx, refreshWindowFunc)
		})
	})

	stack.Layout(ctx, verifyHeadingSection, inputsSection, messageSection, buttonsSection)
}

func (h *SeedHandler) renderInputColumns(ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	inset := layout.Inset{
		Right: unit.Dp(75),
	}
	children := make([]layout.FlexChild, numberOfColumns)

	currentItem := 1
	for i := range h.seed.columns {
		children[i] = flex.Rigid(ctx, func() {
			inset.Layout(ctx, func() {
				h.renderInputColumn(h.seed.columns[i].editors, &currentItem, ctx)
			})
		})
	}
	flex.Layout(ctx, children...)
}

func (h *SeedHandler) renderInputColumn(columnEditors []*widget.Editor, currentItem *int, ctx *layout.Context) {
	stack := layout.Stack{}

	nextInputTopPosition := float32(18)
	editorHeight := 18

	child := stack.Rigid(ctx, func() {
		for _, editor := range columnEditors {
			inset := layout.Inset{
				Top: unit.Dp(nextInputTopPosition),
			}

			inset.Layout(ctx, func() {
				h.renderInputCell(editor, *currentItem, ctx)
			})
			nextInputTopPosition += float32(editorHeight)
			*currentItem++
		}
	})
	stack.Layout(ctx, child)
}

func (h *SeedHandler) renderInputCell(editor *widget.Editor, currentItem int, ctx *layout.Context) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	inset := layout.Inset{}
	numberColumn := flex.Rigid(ctx, func() {
		inset.Layout(ctx, func() {
			h.theme.Body1(strconv.Itoa(currentItem) + ". ").Layout(ctx)
		})
	})

	editorColumn := flex.Rigid(ctx, func() {
		inset.Layout(ctx, func() {
			e := h.theme.Editor("")
			e.Font = text.Font{
				Size: unit.Sp(10),
			}
			e.Layout(ctx, editor)
		})
	})

	flex.Layout(ctx, numberColumn, editorColumn)
}

func (h *SeedHandler) renderVerifyPageButtons(ctx *layout.Context, refreshWindowFunc func()) {
	flex := layout.Flex{
		Axis: layout.Horizontal,
	}

	backButtonSection := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Left: unit.Dp(0),
		}

		inset.Layout(ctx, func() {
			for h.backButton.Clicked(ctx) {
				h.currentSubPage = seedSubPageName
				h.verifyMessage = nil
				refreshWindowFunc()
			}
			widgets.LayoutButton(h.backButton, "Go Back", h.theme, ctx)
		})
	})

	doVerifyButtonSection := flex.Rigid(ctx, func() {
		inset := layout.Inset{
			Left: unit.Dp(300),
		}

		inset.Layout(ctx, func() {
			for h.doVerifyButton.Clicked(ctx) {
				msg := &verifyMessage{}
				if h.isVerificationCorrect() {
					msg.message = "Verification successfull!!"
					msg.messageType = "success"
				} else {
					msg.message = "Invalid verification words"
					msg.messageType = "error"
				}
				h.verifyMessage = msg
				refreshWindowFunc()
			}
			widgets.LayoutButton(h.doVerifyButton, "Verify Words", h.theme, ctx)
		})
	})

	flex.Layout(ctx, backButtonSection, doVerifyButtonSection)
}

func (h *SeedHandler) isVerificationCorrect() bool {
	for _ = range h.seed.columns {
		for columnIndex := range h.seed.columns {
			for itemIndex := range h.seed.columns[columnIndex].words {
				if h.seed.columns[columnIndex].words[itemIndex] != h.seed.columns[columnIndex].editors[itemIndex].Text() {
					return false
				}
			}
		}
	}
	return true
}
