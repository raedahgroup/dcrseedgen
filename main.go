package main

import (
	"image"
	"log"
	"os"
	"strings"

	app "gioui.org/app"
	"gioui.org/font/opentype"
	"gioui.org/text"

	"github.com/markbates/pkger"
	"github.com/raedahgroup/dcrseedgen/helper"
	"github.com/raedahgroup/dcrseedgen/ui"
)

func main() {
	// make data directory if not exists
	err := helper.CreateDataDirectory()
	if err != nil {
		log.Fatalf("error creating data directory: %s", err.Error())
	}

	// load font
	col, err := loadFont()
	if err != nil {
		log.Fatalf("error loading font: %s", err.Error())
	}

	// load decred icons
	decredIcons, err := loadDecredIcons()
	if err != nil {
		log.Fatalf("error loading decred icons: %s", err.Error())
	}

	win := ui.NewWindow(decredIcons, col)
	go win.Loop()

	app.Main()
}

func loadFont() (*text.Collection, error) {
	sans, err := pkger.Open("/assets/fonts/source_sans_pro_regular.otf")
	if err != nil {
		return nil, err
	}

	stat, err := sans.Stat()
	if err != nil {
		return nil, err
	}
	bytes := make([]byte, stat.Size())
	sans.Read(bytes)
	fnt, err := opentype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	col := new(text.Collection)
	col.Register(text.Font{}, fnt)

	return col, nil
}

func loadDecredIcons() (map[string]image.Image, error) {
	decredIcons := make(map[string]image.Image)
	err := pkger.Walk("/assets/decredicons", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() || !strings.HasSuffix(path, ".png") {
			return nil
		}

		f, _ := pkger.Open(path)
		img, _, err := image.Decode(f)
		if err != nil {
			return err
		}
		split := strings.Split(info.Name(), ".")
		decredIcons[split[0]] = img
		return nil
	})

	return decredIcons, err
}
