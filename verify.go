package main

import (
	"image/color"
	"strconv"

	"github.com/aarzilli/nucular"
)

func (handler *RenderHandler) renderVerify(window *nucular.Window) {
	drawHeader(window)

	window.Row(330).Dynamic(1)
	if newWindow := window.GroupBegin("Verify Content", 0); newWindow != nil {
		newWindow.Row(20).Dynamic(1)
		SetFont(newWindow, boldFont)
		newWindow.Label("Verify: ", "LC")

		SetFont(newWindow, normalFont)
		newWindow.Row(235).Dynamic(1)
		if group := newWindow.GroupBegin("", 0); group != nil {
			group.Row(220).Dynamic(5)
			currentItem := 0
			for index := range handler.columns {
				newInputColumn(group, handler.columns[index].inputs, &currentItem)
			}

			group.GroupEnd()
		}

		if handler.verifyMessage.message != "" {
			var color color.RGBA

			switch handler.verifyMessage.messageType {
			case "error":
				color = colorDanger
			case "success":
				color = colorSuccess
			}

			newWindow.Row(20).Dynamic(1)
			newWindow.LabelColored(handler.verifyMessage.message, "LC", color)
		}

		newWindow.Row(40).Ratio(0.5, 0.25, 0.25)
		newWindow.Label("", "LC")
		if newWindow.ButtonText("Verify") {
			msg := &verifyMessage{}
			if handler.doVerify(newWindow) {
				msg.message = "Verification successfull !!"
				msg.messageType = "success"
			} else {
				msg.message = "Invalid mnemonic"
				msg.messageType = "error"
			}
			handler.verifyMessage = msg
		}

		if newWindow.ButtonText("Back") {
			*handler.currentPage = "home"
			newWindow.Master().Changed()
		}
		newWindow.GroupEnd()
	}
}

func newInputColumn(window *nucular.Window, inputs []nucular.TextEditor, currentItem *int) {
	if group := window.GroupBegin(strconv.Itoa(*currentItem), 0); group != nil {
		for index := range inputs {
			group.Row(25).Ratio(0.25, 0.75)
			group.Label(strconv.Itoa(*currentItem+1)+". ", "LC")
			inputs[index].Edit(group)

			*currentItem++
		}
		group.GroupEnd()
	}
}

func (handler *RenderHandler) doVerify(window *nucular.Window) bool {
	for _ = range handler.columns {
		for columnIndex := range handler.columns {
			for itemIndex := range handler.columns[columnIndex].words {
				if handler.columns[columnIndex].words[itemIndex] != string(handler.columns[columnIndex].inputs[itemIndex].Buffer) {
					return false
				}
			}
		}
	}
	return true
}
