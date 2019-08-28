package helper

import (
	"encoding/hex"

	"crypto/ecdsa"
	"crypto/rand"

	"github.com/decred/dcrwallet/walletseed"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrec"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"github.com/decred/dcrd/dcrutil"
)

var (
	curve = secp256k1.S256()
)

func GenerateMnemonicSeed(seedSize uint) (string, string, error) {
	seed, err := walletseed.GenerateRandomSeed(seedSize)
	if err != nil {
		return "", "", err
	}

	return walletseed.EncodeMnemonic(seed), hex.EncodeToString(seed), nil
}

func GenerateAddressAndPrivateKey(selectedNetwork string) (string, string, error) {
	var netPrivKeyID [2]byte
	var chainParams *chaincfg.Params
	switch selectedNetwork {
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
	return privWif.String(), addr.Address(), nil
}
