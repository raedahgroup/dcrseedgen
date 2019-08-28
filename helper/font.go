package helper

import (
	"io/ioutil"

	"fyne.io/fyne"
)

var (
	navFont     fyne.Resource
	headerFont  fyne.Resource
	contentFont fyne.Resource
	boldFont    fyne.Resource
)

func InitFonts() error {
	boldItalicsFontBytes, err := ioutil.ReadFile("assets/font/SourceSansPro-SemiboldIt.ttf")
	if err != nil {
		return err
	}

	semiBoldFontBytes, err := ioutil.ReadFile("assets/font/SourceSansPro-Semibold.ttf")
	if err != nil {
		return err
	}

	regularFontBytes, err := ioutil.ReadFile("assets/font/SourceSansPro-Regular.ttf")
	if err != nil {
		return err
	}

	navFont = getFont("Nav Font", regularFontBytes)
	headerFont = getFont("Header Font", boldItalicsFontBytes)
	contentFont = getFont("Content Font", regularFontBytes)
	boldFont = getFont("Bold Font", semiBoldFontBytes)
	return nil
}

func getFont(fontName string, fontBytes []byte) fyne.Resource {
	return &fyne.StaticResource{
		StaticName:    fontName,
		StaticContent: fontBytes,
	}
}
