package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/aarzilli/nucular"
	nstyle "github.com/aarzilli/nucular/style"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var (
	boldFont   font.Face
	normalFont font.Face

	logo *image.RGBA

	whiteColor = color.RGBA{255, 255, 255, 255}

	colorDanger     = color.RGBA{215, 58, 73, 255}
	colorSuccess    = color.RGBA{227, 98, 9, 255}
	colorPrimary    = color.RGBA{9, 20, 64, 255}
	colorAccent     = color.RGBA{237, 109, 71, 255}
	buttonTextColor = colorAccent

	windowPadding      = image.Point{20, 0}
	groupWindowPadding = image.Point{10, 0}
	noPadding          = image.Point{0, 0}
)

var colorTable = nstyle.ColorTable{
	ColorText:                  whiteColor,
	ColorWindow:                colorPrimary,
	ColorHeader:                color.RGBA{175, 175, 175, 255},
	ColorBorder:                colorAccent,
	ColorButton:                buttonTextColor,
	ColorButtonHover:           whiteColor,
	ColorButtonActive:          color.RGBA{0, 153, 204, 255},
	ColorToggle:                color.RGBA{150, 150, 150, 255},
	ColorToggleHover:           color.RGBA{120, 120, 120, 255},
	ColorToggleCursor:          color.RGBA{175, 175, 175, 255},
	ColorSelect:                color.RGBA{175, 175, 175, 255},
	ColorSelectActive:          color.RGBA{190, 190, 190, 255},
	ColorSlider:                color.RGBA{190, 190, 190, 255},
	ColorSliderCursor:          color.RGBA{80, 80, 80, 255},
	ColorSliderCursorHover:     color.RGBA{70, 70, 70, 255},
	ColorSliderCursorActive:    color.RGBA{60, 60, 60, 255},
	ColorProperty:              color.RGBA{175, 175, 175, 255},
	ColorEdit:                  color.RGBA{150, 150, 150, 255},
	ColorEditCursor:            color.RGBA{0, 0, 0, 255},
	ColorCombo:                 color.RGBA{175, 175, 175, 255},
	ColorChart:                 color.RGBA{160, 160, 160, 255},
	ColorChartColor:            color.RGBA{45, 45, 45, 255},
	ColorChartColorHighlight:   color.RGBA{255, 0, 0, 255},
	ColorScrollbar:             color.RGBA{180, 180, 180, 255},
	ColorScrollbarCursor:       color.RGBA{140, 140, 140, 255},
	ColorScrollbarCursorHover:  color.RGBA{150, 150, 150, 255},
	ColorScrollbarCursorActive: color.RGBA{160, 160, 160, 255},
	ColorTabHeader:             color.RGBA{0x89, 0x89, 0x89, 0xff},
}

func loadFonts() error {
	fontData, err := ioutil.ReadFile("assets/font/SourceSansPro-Regular.ttf")
	if err != nil {
		return err
	}

	normalFont, err = getFont(13, 72, fontData)
	if err != nil {
		return err
	}

	boldFont, err = getFont(19, 105, fontData)

	return nil
}

func loadLogo() error {
	logoHandler, err := os.Open("assets/logo.png")
	if err != nil {
		return err
	}
	defer logoHandler.Close()

	img, _ := png.Decode(logoHandler)
	logo = image.NewRGBA(img.Bounds())
	draw.Draw(logo, img.Bounds(), img, image.ZP, draw.Src)
	return nil
}

func getFont(fontSize, DPI int, fontData []byte) (font.Face, error) {
	ttfont, err := freetype.ParseFont(fontData)
	if err != nil {
		return nil, err
	}

	size := int(float64(fontSize) * scaling)
	options := &truetype.Options{
		Size:    float64(size),
		Hinting: font.HintingFull,
		DPI:     float64(DPI),
	}

	return truetype.NewFace(ttfont, options), nil
}

func drawHeader(window *nucular.Window) {
	window.Row(70).Dynamic(1)

	style := window.Master().Style()
	style.GroupWindow.FixedBackground.Data.Color = whiteColor
	style.NormalWindow.Padding = noPadding
	style.GroupWindow.Padding = noPadding
	window.Master().SetStyle(style)

	if group := window.GroupBegin("header", 0); window != nil {
		// style window
		group.Row(60).Dynamic(1)
		group.Image(logo)
		group.GroupEnd()
	}
	window.Row(10).Dynamic(1)
	window.Label("", "LC")

	style = window.Master().Style()
	style.GroupWindow.FixedBackground.Data.Color = colorTable.ColorWindow
	style.GroupWindow.Padding = groupWindowPadding
	window.Master().SetStyle(style)
}

func SetFont(window *nucular.Window, font font.Face) {
	masterWindow := window.Master()
	style := masterWindow.Style()
	style.Font = font
	masterWindow.SetStyle(style)
}

func setStyle(window nucular.MasterWindow) error {
	err := loadFonts()
	if err != nil {
		return fmt.Errorf("error loading font: %s", err.Error())
	}

	style := nstyle.FromTable(colorTable, scaling)
	style.Font = normalFont

	// window
	style.NormalWindow.Padding = windowPadding

	// buttons
	style.Button.Rounding = 0
	style.Button.TextHover = colorAccent

	// text input
	style.Edit.Normal.Data.Color = whiteColor
	style.Edit.Hover.Data.Color = whiteColor
	style.Edit.Active.Data.Color = whiteColor
	style.Edit.TextActive = colorAccent
	style.Edit.TextNormal = colorAccent
	style.Edit.TextHover = colorAccent
	style.Edit.CursorHover = colorAccent

	window.SetStyle(style)

	return nil
}
