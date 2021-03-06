package src

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/decred/dcrd/dcrec/secp256k1/v2"
)

// transation struct
type Tx struct {
	version     []byte
	nbIn, nbOut int
	inputs      []TxIn
	outputs     []TxOut
	lock_time   time.Time
}

// txIn struct
type TxIn struct {
	previousTx    []byte
	previousTxOut []byte
	sig           []byte
	pubKey        []byte
	sequence      []byte
}

// txOut struct
type TxOut struct {
	Value          []byte
	TxOutScriptLen []byte
	ScriptPubKey   []byte
}

func InitializeTransaction() (tx Tx) {
	now := time.Now()

	in := make([]byte, 8)
	i := int64(0)
	binary.BigEndian.PutUint64(in, uint64(i))

	out := make([]byte, 8)
	i = int64(0)
	binary.BigEndian.PutUint64(out, uint64(i))

	new_tx := Tx{
		version:   []byte{1, 0, 0, 0},
		lock_time: now,
		nbIn:      0,
		nbOut:     0,
	}
	return new_tx
}

func (tx *Tx) AddOutput(out *TxOut) {
	tx.outputs = append(tx.outputs, *out)
	tx.nbOut += 1
}

func (tx *Tx) AddInput(input *TxIn) {
	tx.inputs = append(tx.inputs, *input)

	tx.nbIn += 1
}

func Transaction(pubKey, destAddress, txHash string, amount int64, out *TxOut) (Tx, error) {
	tx := InitializeTransaction()

	value := int64(binary.LittleEndian.Uint64(out.Value))

	if int64(value) < amount {
		return tx, fmt.Errorf("not enough btc in balance")
	}

	temp, _ := hex.DecodeString(txHash)
	fmt.Println(len(temp))

	txInput := TxIn{
		pubKey:        []byte(pubKey),
		previousTx:    temp,
		previousTxOut: []byte{1, 0, 0, 0},
		sequence:      []byte{255, 255, 255, 255},
	}

	txOutput := TxOut{
		Value:          out.Value,
		TxOutScriptLen: VarIntConverter(len(destAddress)),
		ScriptPubKey:   []byte(destAddress),
	}

	tx.AddInput(&txInput)
	tx.AddOutput(&txOutput)
	return tx, nil
}

func (tx *Tx) Serialize() ([]byte, error) {
	res := []byte{}

	// Append version number
	res = append(res, tx.version...)

	// Append input counter
	in := VarIntConverter(tx.nbIn)
	res = append(res, in...)

	// Append input
	for _, elt := range tx.inputs {
		if len(elt.previousTx) != 32 {
			return nil, fmt.Errorf("error in previous_tx")
		}

		// Append previous transaction hash
		res = append(res, elt.previousTx...)

		if len(elt.previousTxOut) != 4 {
			return nil, fmt.Errorf("error in previous_index")
		}

		// Append previous TxOut-index
		res = append(res, elt.previousTxOut...)

		// Append previous sequence
		res = append(res, elt.sequence...)
	}
	// Append ouput counter
	out := VarIntConverter(tx.nbOut)
	res = append(res, out...)

	// Append output
	for _, elt := range tx.outputs {
		if len(elt.Value) != 8 {
			return nil, fmt.Errorf("error in value")
		}

		// Append value
		res = append(res, elt.Value...)
	}

	return res, nil
}

func (tx *Tx) SerializeSignature(privateKey *secp256k1.PrivateKey) ([]byte, error) {
	txTemp, err := tx.Serialize()
	if err != nil {
		return nil, err
	}

	sig, err := SigTransaction(txTemp, privateKey)
	if err != nil {
		return nil, err
	}

	res := []byte{}

	res = append(res, tx.version...)

	// Append input counter
	in := VarIntConverter(tx.nbIn)
	res = append(res, in...)

	for _, elt := range tx.inputs {
		if len(elt.previousTx) != 32 {
			return nil, fmt.Errorf("error in previous_tx")
		}

		// Append previous transaction hash
		res = append(res, elt.previousTx...)

		if len(elt.previousTxOut) != 4 {
			return nil, fmt.Errorf("error in previous_index")
		}

		// Append previous TxOut-index
		res = append(res, elt.previousTxOut...)

		scriptSig := txscript.NewScriptBuilder()
		scriptSig.AddData(sig)
		scriptSig.AddData(elt.pubKey)

		txInScript, err := scriptSig.Script()
		if err != nil {
			return nil, err
		}
		res = append(res, VarIntConverter(len(txInScript))...)
		fmt.Print("scriptSig: ")
		fmt.Println(txInScript)
		res = append(res, txInScript...)
		res = append(res, elt.sequence...)
	}

	// Append ouput counter
	out := VarIntConverter(tx.nbOut)
	res = append(res, out...)

	for _, elt := range tx.outputs {
		if len(elt.Value) != 8 {
			return nil, fmt.Errorf("error in value")
		}

		res = append(res, elt.Value...)

		if len(elt.TxOutScriptLen) > 9 {
			return nil, fmt.Errorf("error in scriptlen")
		}
		address, err := btcutil.DecodeAddress(string(elt.ScriptPubKey), &chaincfg.MainNetParams)

		if err != nil {
			return nil, err
		}

		script, err := txscript.PayToAddrScript(address)

		if err != nil {
			return nil, err
		}

		res = append(res, VarIntConverter(len(script))...)
		fmt.Print("test")
		fmt.Println(script)
		res = append(res, script...)
	}
	return res, nil
}

func SigTransaction(txSerialized []byte, privateKey *secp256k1.PrivateKey) ([]byte, error) {
	txHash := sha256.Sum256(txSerialized)

	signature, err := privateKey.Sign(txHash[:])
	if err != nil {
		return nil, err
	}
	return signature.Serialize(), nil
}
