package helper

import (
	"encoding/hex"

	"github.com/decred/dcrwallet/walletseed"
)

func GenerateMnemonicSeed(seedSize uint) (string, string, error) {
	seed, err := walletseed.GenerateRandomSeed(seedSize)
	if err != nil {
		return "", "", err
	}

	return walletseed.EncodeMnemonic(seed), hex.EncodeToString(seed), nil
}
