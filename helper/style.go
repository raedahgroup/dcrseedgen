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

	normalWindowPadding = image.Point{10, 0}
	groupWindowPadding  = image.Point{10, 0}
	noPadding           = image.Point{0, 0}

	logo *image.RGBA
)

const (
	scaling = 1.1
)

func LoadLogo() error {
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
	style.Button.TextHover = colorAccent

	// text inputs
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

func loadFonts() error {
	fontData, err := ioutil.ReadFile("assets/font/SourceSansPro-Regular.ttf")
	if err != nil {
		return err
	}

	FontNormal, err = getFont(13, 72, fontData)
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
