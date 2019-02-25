package handlers

import (
	"path/filepath"

	"github.com/aarzilli/nucular"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrseedgen/helper"
)

type AddressGeneratorHandler struct {
	serverAddressInput    nucular.TextEditor
	rpccertInput          nucular.TextEditor
	tlsOptions            []string
	walletPassphraseInput nucular.TextEditor

	err error
}

const (
	walletRPCServerAddress = "localhost:19111"
	useWalletRPCTLS        = true
)

var (
	defaultDcrwalletAppDataDir = dcrutil.AppDataDir("dcrwallet", false)
	defaultRPCCertFile         = filepath.Join(defaultDcrwalletAppDataDir, "rpc.cert")
)

func (a *AddressGeneratorHandler) BeforeRender() {
	a.tlsOptions = []string{"Yes", "No"}
	a.walletPassphraseInput.PasswordChar = '*'

	a.rpccertInput.Buffer = []rune(defaultRPCCertFile)
	a.serverAddressInput.Buffer = []rune(walletRPCServerAddress)
}

func (a *AddressGeneratorHandler) Render(window *nucular.Window) {
	if a.err != nil {
		window.Row(20).Dynamic(1)
		window.Label(a.err.Error(), "LC")
		return
	}

	window.Row(360).Dynamic(1)
	if w := helper.NewWindow("Address Page Content", window, 0); w != nil {
		w.Row(20).Dynamic(2)
		w.Label("Wallet RPC Server Address:", "LC")
		w.Label("RPC Cert:", "LC")

		w.Row(30).Dynamic(2)
		a.serverAddressInput.Edit(w.Window)
		a.rpccertInput.Edit(w.Window)

		w.Row(20).Dynamic(2)
		w.Label("RPC Cert File:", "LC")
		w.Label("Wallet Passphrase:", "LC")

		w.Row(30).Dynamic(2)
		w.ComboSimple(a.tlsOptions, 0, 30)
		a.walletPassphraseInput.Edit(w.Window)

		w.Row(10).Dynamic(1)
		w.Label("", "LC")

		w.Row(40).Dynamic(3)
		if w.ButtonText("Generate Address") {

		}

		w.End()
	}
}
