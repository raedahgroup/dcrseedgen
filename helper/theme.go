package helper

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

const (
	padding            = 7
	fontSize           = 19
	scrollbarSize      = 5
	smallScrollbarSize = 3
)

var (
	primaryColor    = color.RGBA{9, 20, 64, 255} // color.RGBA{112, 203, 255, 255}
	secondaryColor  = color.RGBA{237, 109, 71, 255}
	backgroundColor = primaryColor
	whiteColor      = color.RGBA{255, 255, 255, 255}
	greyColor       = color.RGBA{200, 200, 200, 255}

	DefaultTheme = &Theme{}
)

type Theme struct {
}

func (Theme) BackgroundColor() color.Color {
	return backgroundColor
}

func (Theme) ButtonColor() color.Color {
	return secondaryColor
}

func (Theme) DisabledButtonColor() color.Color {
	return secondaryColor
}

func (Theme) HyperlinkColor() color.Color {
	return secondaryColor
}

func (Theme) TextColor() color.Color {
	return whiteColor
}

func (Theme) DisabledTextColor() color.Color {
	return greyColor
}

func (Theme) IconColor() color.Color {
	return whiteColor
}

func (Theme) DisabledIconColor() color.Color {
	return greyColor
}

func (Theme) PlaceHolderColor() color.Color {
	return whiteColor
}

func (Theme) PrimaryColor() color.Color {
	return primaryColor
}

func (Theme) HoverColor() color.Color {
	return secondaryColor
}

func (Theme) FocusColor() color.Color {
	return secondaryColor
}

func (Theme) ScrollBarColor() color.Color {
	return secondaryColor
}

func (Theme) ShadowColor() color.Color {
	return &color.RGBA{0xcc, 0xcc, 0xcc, 0xcc}
}

func (Theme) TextSize() int {
	return fontSize
}

func (Theme) TextFont() fyne.Resource {
	return contentFont
}

func (Theme) TextBoldFont() fyne.Resource {
	return boldFont
}

func (Theme) TextItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

func (Theme) TextBoldItalicFont() fyne.Resource {
	return theme.DefaultTextBoldItalicFont()
}

func (Theme) TextMonospaceFont() fyne.Resource {
	return theme.DefaultTextMonospaceFont()
}

func (Theme) Padding() int {
	return padding
}

func (Theme) IconInlineSize() int {
	return 50
}

func (Theme) ScrollBarSize() int {
	return scrollbarSize
}

func (Theme) ScrollBarSmallSize() int {
	return smallScrollbarSize
}
