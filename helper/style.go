package helper

import (
	"fmt"
	"image"
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
	FontBold   font.Face
	FontNormal font.Face

	normalWindowPadding = image.Point{8, 0}
	groupWindowPadding  = image.Point{8, 0}
	noPadding           = image.Point{0, 0}

	logo *image.RGBA

	ExportIcon *image.RGBA
)

const (
	scaling      = 1.2
	ButtonHeight = 40
)

func LoadLogo() error {
	l, err := loadImage("assets/logo.png")
	if err != nil {
		return err
	}
	logo = l
	return nil
}

func LoadIcons() error {
	l, err := loadImage("assets/export_icon.png")
	if err != nil {
		return err
	}
	ExportIcon = l
	return nil
}

func loadImage(path string) (*image.RGBA, error) {
	handler, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer handler.Close()

	img, _ := png.Decode(handler)
	dest := image.NewRGBA(img.Bounds())
	draw.Draw(dest, img.Bounds(), img, image.ZP, draw.Src)
	return dest, nil
}

func InitStyle(window nucular.MasterWindow) error {
	if err := loadFonts(); err != nil {
		return fmt.Errorf("error loading fonts: %s", err.Error())
	}

	style := nstyle.FromTable(colorTable, scaling)
	style.Font = FontNormal

	// style normal window
	style.NormalWindow.Padding = noPadding

	// style buttons
	style.Button.Rounding = 0
	style.Button.TextHover = colorSecondary

	// text inputs
	style.Edit.Normal.Data.Color = whiteColor
	style.Edit.Hover.Data.Color = whiteColor
	style.Edit.Active.Data.Color = whiteColor
	style.Edit.TextActive = colorAccent
	style.Edit.TextNormal = colorAccent
	style.Edit.TextHover = colorAccent
	style.Edit.CursorHover = colorAccent

	// combo
	style.Combo.LabelActive = colorSecondary
	style.Combo.LabelHover = colorSecondary
	style.Combo.LabelNormal = colorSecondary
	style.Combo.Active.Data.Color = whiteColor
	style.Combo.Hover.Data.Color = whiteColor
	style.Combo.Normal.Data.Color = whiteColor
	style.Combo.Button.Normal.Data.Color = colorSecondary
	style.Combo.Button.Hover.Data.Color = colorSecondary
	style.Combo.Button.Active.Data.Color = colorSecondary
	style.Combo.Button.Border = 1
	style.Combo.Button.BorderColor = colorSecondary

	// combo window
	style.ComboWindow.Spacing = noPadding
	style.ComboWindow.Padding = noPadding
	style.ComboWindow.ScrollbarSize = noPadding

	style.ComboWindow.Background = colorAccent
	style.ComboWindow.Scaler.Data.Color = colorAccent
	style.ComboWindow.FixedBackground.Data.Color = colorAccent

	window.SetStyle(style)

	return nil
}

func loadFonts() error {
	fontData, err := ioutil.ReadFile("assets/font/SourceSansPro-Regular.ttf")
	if err != nil {
		return err
	}

	FontNormal, err = getFont(14, 72, fontData)
	if err != nil {
		return err
	}

	FontBold, err = getFont(19, 105, fontData)

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

func UseFont(window *Window, font font.Face) {
	style := window.Master().Style()
	style.Font = font
	window.Master().SetStyle(style)
}

func StyleClipboardInput(window *Window) {
	style := window.Master().Style()
	style.Font = FontNormal
	style.Edit.Border = 0
	style.Edit.Normal.Data.Color = colorPrimary
	style.Edit.Hover.Data.Color = colorPrimary
	style.Edit.Active.Data.Color = colorPrimary
	style.Edit.TextActive = whiteColor
	style.Edit.TextNormal = whiteColor
	style.Edit.TextHover = whiteColor
	style.Edit.CursorHover = whiteColor
	window.Master().SetStyle(style)
}

func ResetInputStyle(window *Window) {
	style := window.Master().Style()
	style.Font = FontNormal
	style.Edit.Border = 1
	style.Edit.Normal.Data.Color = whiteColor
	style.Edit.Hover.Data.Color = whiteColor
	style.Edit.Active.Data.Color = whiteColor
	style.Edit.TextActive = colorSecondary
	style.Edit.TextNormal = colorSecondary
	style.Edit.TextHover = colorSecondary
	window.Master().SetStyle(style)
}

func StyleNavButton(window *nucular.Window) {
	style := window.Master().Style()
	style.Button.Border = 0
	style.Button.Normal.Data.Color = colorAccent
	style.Button.Hover.Data.Color = whiteColor
	style.Button.TextHover = colorAccent
	window.Master().SetStyle(style)
}

func ResetButtonStyle(window *nucular.Window) {
	style := window.Master().Style()
	style.Button.Border = 1
	style.Button.Normal.Data.Color = colorSecondary
	style.Button.Hover.Data.Color = whiteColor
	window.Master().SetStyle(style)
}
