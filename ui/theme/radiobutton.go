// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"gioui.org/layout"
	//"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type RadioButton struct {
	material.RadioButtonStyle
}

// RadioButton returns a RadioButton with a label. The key specifies
// the value for the Enum.
func (t *Theme) RadioButton(key, label string, enum *widget.Enum) RadioButton {
	return RadioButton{
		material.RadioButton(t.Theme, enum, key, label),
	}
}

func (r RadioButton) Layout(gtx layout.Context) layout.Dimensions {
	return r.RadioButtonStyle.Layout(gtx)
}
