package handlers

import (
	"crypto/ecdsa"
	"crypto/rand"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type AddressGeneratorHandler struct {
	err              error
	netOptions       []string
	selectedNetIndex int

	address    string
	privateKey string

	addressInput    nucular.TextEditor
	privateKeyInput nucular.TextEditor
}

var (
	curve = secp256k1.S256()
)

func (a *AddressGeneratorHandler) BeforeRender() {
	a.netOptions = []string{"Mainnet", "Testnet"}
	a.selectedNetIndex = 0

	a.addressInput.Flags = nucular.EditClipboard | nucular.EditNoCursor | nucular.EditBox
	a.privateKeyInput.Flags = nucular.EditClipboard | nucular.EditNoCursor | nucular.EditBox
}

func (a *AddressGeneratorHandler) Render(window *nucular.Window) {
	if a.err != nil {
		window.Row(20).Dynamic(1)
		window.Label(a.err.Error(), "LC")
		return
	}

	window.Row(360).Dynamic(1)
	if w := helper.NewWindow("Address Page Content", window, 0); w != nil {
		w.Row(30).Dynamic(4)
		w.ComboSimple(a.netOptions, a.selectedNetIndex, 30)

		if w.ButtonText("Generate Address") {
			a.generateAddressAndPrivateKey(w)
		}

		if a.address != "" && a.privateKey != "" {
			w.Row(10).Dynamic(1)
			w.Label("", "LC")

			w.Row(25).Dynamic(1)
			helper.UseFont(w, helper.FontBold)
			w.Label("Address:", "LC")

			w.Row(35).Dynamic(1)
			helper.UseFont(w, helper.FontNormal)
			helper.StyleClipboardInput(w)
			a.addressInput.Edit(w.Window)
			helper.ResetInputStyle(w)

			w.Row(10).Dynamic(1)
			w.Label("", "LC")

			w.Row(25).Dynamic(1)
			helper.UseFont(w, helper.FontBold)
			w.Label("Private Key:", "LC")

			w.Row(35).Dynamic(1)
			helper.StyleClipboardInput(w)
			a.privateKeyInput.Edit(w.Window)
			helper.ResetInputStyle(w)
		}
		w.End()
	}
}

func (a *AddressGeneratorHandler) generateAddressAndPrivateKey(window *helper.Window) {
	defer window.Master().Changed()

	var chainParam chaincfg.Params
	if a.netOptions[a.selectedNetIndex] == "Mainnet" {
		chainParam = chaincfg.MainNetParams
	} else {
		chainParam = chaincfg.TestNet3Params
	}

	key, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		a.err = err
		return
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
		&chainParam,
		dcrec.STEcdsaSecp256k1)
	if err != nil {
		a.err = err
		return
	}

	privWif, err := dcrutil.NewWIF(priv, &chainParam, dcrec.STEcdsaSecp256k1)
	if err != nil {
		a.err = err
		return
	}

	a.address = addr.EncodeAddress()
	a.privateKey = privWif.String()

	a.addressInput.Buffer = []rune(a.address)
	a.privateKeyInput.Buffer = []rune(a.privateKey)
}
