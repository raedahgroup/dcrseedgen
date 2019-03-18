package handlers

import (
	"crypto/ecdsa"
	"crypto/rand"
	"strconv"

	"github.com/aarzilli/nucular"
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
}

var (
	curve = secp256k1.S256()
)

func (a *AddressGeneratorHandler) BeforeRender() {
	a.netOptions = []string{"Mainnet", "Testnet"}

	a.selectedNetIndex = 0

	a.quantityInput.Flags = nucular.EditClipboard
	a.quantityInput.Buffer = []rune("1")
}

func (a *AddressGeneratorHandler) Render(window *nucular.Window) {
	if a.err != nil {
		window.Row(20).Dynamic(1)
		window.Label(a.err.Error(), "LC")
	}

	window.Row(360).Dynamic(1)
	if w := helper.NewWindow("Address Page Content", window, 0); w != nil {
		w.Row(helper.ButtonHeight).Ratio(0.2, 0.2)
		w.ComboSimple(a.netOptions, a.selectedNetIndex, 30)

		if w.ButtonText("Generate Address") {
			a.generateAddressAndPrivateKey(w)
		}

		helper.UseFont(w, helper.FontBold)
		if len(a.inputPairs) > 0 {
			w.Row(10).Dynamic(1)
			w.Label("", "LC")

			w.Row(30).Ratio(0.4, 0.6)
			w.Label("Address:", "LC")
			w.Label("Private Key:", "LC")
		}

		helper.UseFont(w, helper.FontNormal)
		helper.StyleClipboardInput(w)

		w.Row(30).Ratio(0.4, 0.6)
		for i := range a.inputPairs {
			a.inputPairs[i].addressInput.Edit(w.Window)
			a.inputPairs[i].privateKeyInput.Edit(w.Window)
		}
		helper.ResetInputStyle(w)
		w.End()
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
	privWif, err := dcrutil.NewWIF(priv, &chainParam, dcrec.STEcdsaSecp256k1)
	if err != nil {
		return "", "", err
	}

	return addr.EncodeAddress(), privWif.String(), nil
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
