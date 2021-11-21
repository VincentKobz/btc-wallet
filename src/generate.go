package src

import (
	"crypto/sha256"

	"github.com/btcsuite/btcutil/base58"
	"github.com/decred/dcrd/dcrec/secp256k1"
	"golang.org/x/crypto/ripemd160"
)

// GenerateBtcAddress: Generate new BTC address
func GenerateBtcAddress() (string, *secp256k1.PrivateKey, error) {
	privateKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return "", nil, err
	}

	publicKey := privateKey.PubKey()

	publicHashSha256 := sha256.Sum256(publicKey.SerializeCompressed())
	publicHashRipemd := ripemd160.New()
	publicHashRipemd.Write(publicHashSha256[:])

	temp := publicHashRipemd.Sum(nil)
	temp = append([]byte{0}, temp...)

	publicHash := sha256.Sum256(temp)
	publicHash = sha256.Sum256(publicHash[:])

	for i := 0; i < 4; i++ {
		temp = append(temp, publicHash[i])
	}

	btcAddress := base58.Encode(temp)

	return btcAddress, privateKey, nil
}
