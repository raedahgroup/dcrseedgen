// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (t *Theme) H1(txt string) material.LabelStyle {
	return material.H1(t.Theme, txt)
}

func (t *Theme) H2(txt string) material.LabelStyle {
	return material.H2(t.Theme, txt)
}

func (t *Theme) H3(txt string) material.LabelStyle {
	return material.H3(t.Theme, txt)
}

func (t *Theme) H4(txt string) material.LabelStyle {
	return material.H4(t.Theme, txt)
}

func (t *Theme) H5(txt string) material.LabelStyle {
	return material.H5(t.Theme, txt)
}

func (t *Theme) H6(txt string) material.LabelStyle {
	return material.H6(t.Theme, txt)
}

func (t *Theme) Body1(txt string) material.LabelStyle {
	return material.Body1(t.Theme, txt)
}

func (t *Theme) Body2(txt string) material.LabelStyle {
	return material.Body2(t.Theme, txt)
}

func (t *Theme) Caption(txt string) material.LabelStyle {
	return material.Caption(t.Theme, txt)
}

func (t *Theme) ErrorLabel(txt string) material.LabelStyle {
	label := t.Caption(txt)
	label.Color = t.Color.Danger

	return label
}

func (t *Theme) Label(size unit.Value, txt string) material.LabelStyle {
	return material.Label(t.Theme, size, txt)
}
