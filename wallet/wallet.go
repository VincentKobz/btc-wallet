package wallet

import "github.com/VincentKobz/btc-wallet/src"

type Wallet struct {
	amount           int64
	publicAddressBtc string
	utxo             []src.TxOut
}
