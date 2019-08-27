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
		w.Row(helper.ButtonHeight).Ratio(0.2, 0.2)
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
		chainParams,
		dcrec.STEcdsaSecp256k1)
	if err != nil {
		a.err = err
		return
	}

	privWif := dcrutil.NewWIF(priv, netPrivKeyID, dcrec.STEcdsaSecp256k1)
	a.privateKey = privWif.String()
	a.address = addr.Address()

	a.addressInput.Buffer = []rune(a.address)
	a.privateKeyInput.Buffer = []rune(a.privateKey)
}
