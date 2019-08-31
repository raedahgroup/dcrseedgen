package handlers

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/csv"
	"image"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aarzilli/nucular"
	"github.com/aarzilli/nucular/label"
	"github.com/aarzilli/nucular/rect"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type inputPair struct {
	addressInput    nucular.TextEditor
	privateKeyInput nucular.TextEditor
}

type AddressGeneratorHandler struct {
	err              error
	netOptions       []string
	selectedNetIndex int

	address    string
	privateKey string

	quantityInput nucular.TextEditor

	inputPairs []inputPair

	popupBounds     rect.Rect
	csvFolderInput  nucular.TextEditor
	csvInputError   string
	isGeneratingCSV bool
}

var (
	curve = secp256k1.S256()
)

func (a *AddressGeneratorHandler) BeforeRender() {
	a.netOptions = []string{"Mainnet", "Testnet"}

	a.selectedNetIndex = 0

	a.quantityInput.Flags = nucular.EditClipboard
	a.quantityInput.Buffer = []rune("1")

	a.csvFolderInput.Flags = nucular.EditClipboard | nucular.EditSimple

	a.popupBounds = rect.Rect{
		X: 130,
		Y: 60,
		W: 390,
		H: 150,
	}
}

func (a *AddressGeneratorHandler) Render(window *nucular.Window) {
	if a.err != nil {
		window.Row(20).Dynamic(1)
		window.Label(a.err.Error(), "LC")
	}

	window.Row(300).Dynamic(1)
	if w := helper.NewWindow("Address Page Content", window, 0); w != nil {
		w.Row(10).Dynamic(3)
		w.Label("Net Type:", "LC")
		w.Label("Number to generate:", "LC")
		w.Label("", "LC")

		w.Row(35).Static(200, 200, 150)
		a.selectedNetIndex = w.ComboSimple(a.netOptions, a.selectedNetIndex, 30)
		a.quantityInput.Edit(w.Window)

		if w.ButtonText("Generate") {
			a.generateAddressAndPrivateKey(w)
		}

		helper.UseFont(w, helper.FontBold)
		if len(a.inputPairs) > 0 {
			w.Row(10).Dynamic(1)
			w.Label("", "LC")

			w.Row(26).Ratio(0.4, 0.6)
			w.Label("Address:", "LC")
			w.Label("Private Key:", "LC")
		}

		helper.UseFont(w, helper.FontNormal)
		helper.StyleClipboardInput(w)

		w.Row(27).Ratio(0.39, 0.61)
		for i := range a.inputPairs {
			a.inputPairs[i].addressInput.Edit(w.Window)
			a.inputPairs[i].privateKeyInput.Edit(w.Window)
		}
		helper.ResetInputStyle(w)
		w.End()
	}

	if len(a.inputPairs) > 0 {
		window.Row(50).Dynamic(1)
		if w := helper.NewWindow("Create csv button", window, 0); w != nil {
			w.Row(25).Dynamic(10)
			if w.ButtonText("") {
				w.Master().PopupOpen("Export as csv", nucular.WindowTitle|nucular.WindowDynamic|nucular.WindowNoScrollbar, a.popupBounds, true, a.renderCSVPopup)
			}

			if w.Input().Mouse.HoveringRect(w.LastWidgetBounds) {
				w.Tooltip("Export as csv")
			}

			w.End()
		}
	}
}

func (a *AddressGeneratorHandler) doGenerate(window *helper.Window) (string, string, error) {
	defer window.Master().Changed()

	var netPrivKeyID [2]byte
	var chainParams *chaincfg.Params
	switch a.netOptions[a.selectedNetIndex] {
	case "Testnet3":
		netPrivKeyID = [2]byte{0x23, 0x0e} // starts with Pt
		chainParams = chaincfg.TestNet3Params()
	case "Regnet":
		netPrivKeyID = [2]byte{0x22, 0xfe} // starts with Pr
		chainParams = chaincfg.RegNetParams()
	case "Simnet":
		netPrivKeyID = [2]byte{0x23, 0x07} // starts with Ps
		chainParams = chaincfg.SimNetParams()
	default:
		netPrivKeyID = [2]byte{0x22, 0xde} // starts with Pm
		chainParams = chaincfg.MainNetParams()

	}

	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return "", "", err
	}
	pub := secp256k1.PublicKey{
		Curve: curve,
		X:     key.PublicKey.X,
		Y:     key.PublicKey.Y,
	}
	priv := secp256k1.PrivateKey{
		PublicKey: key.PublicKey,
		D:         key.D,
	}

	addr, err := dcrutil.NewAddressPubKeyHash(
		dcrutil.Hash160(pub.SerializeCompressed()),
		chainParams,
		dcrec.STEcdsaSecp256k1)
	if err != nil {
		return "", "", err
	}

	privWif := dcrutil.NewWIF(priv, netPrivKeyID, dcrec.STEcdsaSecp256k1)

	a.privateKey = privWif.String()
	a.address = addr.Address()
	if err != nil {
		return "", "", err
	}

	return addr.String(), privWif.String(), nil
}

func (a *AddressGeneratorHandler) generateAddressAndPrivateKey(window *helper.Window) {
	numberToGenerate, err := strconv.Atoi(string(a.quantityInput.Buffer))
	if err != nil {
		a.err = err
		return
	}

	a.inputPairs = make([]inputPair, numberToGenerate)
	for i := 0; i < numberToGenerate; i++ {
		address, privateKey, err := a.doGenerate(window)
		if err != nil {
			a.err = err
			return
		}

		a.inputPairs[i].addressInput.Flags = nucular.EditClipboard | nucular.EditNoCursor | nucular.EditBox
		a.inputPairs[i].privateKeyInput.Flags = nucular.EditClipboard | nucular.EditNoCursor | nucular.EditBox

		a.inputPairs[i].addressInput.Buffer = []rune(address)
		a.inputPairs[i].privateKeyInput.Buffer = []rune(privateKey)
	}
}

func (a *AddressGeneratorHandler) renderCSVPopup(window *nucular.Window) {
	masterWindow := window.Master()

	// set popup style
	style := window.Master().Style()
	style.NormalWindow.Padding = image.Point{20, 50}
	//style.NormalWindow.Background = color.RGBA{0xff, 0xff, 0xff, 0xff}
	masterWindow.SetStyle(style)

	defer func() {
		// reset page style
		style.NormalWindow.Padding = image.Point{0, 0}
		masterWindow.SetStyle(style)
	}()

	// render popup
	window.Row(20).Dynamic(1)
	window.Label("Target folder:", "LC")

	window.Row(25).Dynamic(1)
	a.csvFolderInput.Edit(window)

	if a.csvInputError != "" {
		window.Row(10).Dynamic(1)
		window.LabelColored(a.csvInputError, "LC", color.RGBA{205, 32, 32, 255})
	}

	window.Row(25).Static(65, 65)
	if window.Button(label.T("Close"), false) {
		window.Close()
	}

	buttonText := "Submit"
	if a.isGeneratingCSV {
		buttonText = "Submitting..."
	}

	if window.Button(label.T(buttonText), false) {
		if a.validateCSVForm() {
			if a.generateCSV() {
				window.Close()
			}
			return
		}
		masterWindow.Changed()
	}
}

func (a *AddressGeneratorHandler) validateCSVForm() bool {
	targetFolder := string(a.csvFolderInput.Buffer)

	if targetFolder == "" {
		a.csvInputError = "Target folder is required"
		return false
	}

	a.csvInputError = ""
	return true
}

func (a *AddressGeneratorHandler) generateCSV() bool {
	filename := "dcrseedgen_address_" + strconv.Itoa(int(time.Now().Unix())) + ".csv"
	fp := filepath.Join(string(a.csvFolderInput.Buffer), filename)

	file, err := os.Create(fp)
	if err != nil {
		a.csvInputError = err.Error()
		return false
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := [][]string{{"Address", "Private Key"}}

	for _, inputPair := range a.inputPairs {
		item := []string{string(inputPair.addressInput.Buffer), string(inputPair.privateKeyInput.Buffer)}
		data = append(data, item)
	}

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			a.csvInputError = err.Error()
			return false
		}
	}

	return true
}
