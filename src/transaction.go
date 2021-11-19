package src

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// transation struct
type Tx struct {
	version     []byte
	flag        []byte
	nbIn, nbOut []byte
	inputs      []TxIn
	outputs     []TxOut
	lock_time   time.Time
}

// txIn struct
type TxIn struct {
	previousTx    []byte
	previousTxOut []byte
	txInScriptLen []byte
	txInScript    []byte
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
	binary.LittleEndian.PutUint64(in, uint64(i))

	out := make([]byte, 8)
	i = int64(0)
	binary.LittleEndian.PutUint64(out, uint64(i))

	new_tx := Tx{
		version:   []byte{0, 0, 0, 1},
		flag:      []byte{0, 1},
		lock_time: now,
		nbIn:      in,
		nbOut:     out,
	}

	return new_tx
}

func AddOutput(tx *Tx, out *TxOut) {
	tx.outputs = append(tx.outputs, *out)

	temp := int64(binary.LittleEndian.Uint64(tx.nbOut))
	temp += 1
	binary.LittleEndian.PutUint64(tx.nbOut, uint64(temp))
}

func AddInput(tx *Tx, input *TxIn) {
	tx.inputs = append(tx.inputs, *input)

	temp := int64(binary.LittleEndian.Uint64(tx.nbIn))
	temp += 1
	binary.LittleEndian.PutUint64(tx.nbIn, uint64(temp))
}

func Transaction(destAddress, txHash string, amount int64, out *TxOut) (Tx, error) {
	tx := InitializeTransaction()

	value := int64(binary.LittleEndian.Uint64(out.Value))

	if int64(value) < amount {
		return tx, fmt.Errorf("not enough btc in balance")
	}

	temp, _ := hex.DecodeString(txHash)
	fmt.Println(len(temp))

	txInput := TxIn{
		txInScript:    []byte(destAddress),
		txInScriptLen: []byte(strconv.Itoa(len(destAddress))),
		previousTx:    temp,
		previousTxOut: []byte{0, 0, 0, 0},
		sequence:      []byte{255, 255, 255, 255},
	}

	txOutput := TxOut{
		Value:          out.Value,
		TxOutScriptLen: []byte(strconv.Itoa(len(destAddress))),
		ScriptPubKey:   []byte(destAddress),
	}

	AddInput(&tx, &txInput)
	AddOutput(&tx, &txOutput)
	return tx, nil
}

func SerializeTransaction(tx *Tx) ([]byte, error) {
	res := []byte{}

	res = append(res, tx.version...)
	res = append(res, tx.flag...)

	if len(tx.nbIn) > 9 {
		return nil, fmt.Errorf("error in nb_count")
	}
	res = append(res, tx.nbIn...)

	for _, elt := range tx.inputs {

		if len(elt.previousTx) != 32 {
			return nil, fmt.Errorf("error in previous_tx")
		}

		res = append(res, elt.previousTx...)

		fmt.Println(len(elt.previousTxOut))
		if len(elt.previousTxOut) != 4 {
			return nil, fmt.Errorf("error in previous_index")
		}

		res = append(res, elt.previousTxOut...)

		if len(elt.txInScriptLen) > 9 {
			return nil, fmt.Errorf("error in scriptlen")
		}

		res = append(res, elt.txInScriptLen...)

		res = append(res, elt.txInScript...)
		res = append(res, elt.sequence...)
	}

	if len(tx.nbOut) > 9 {
		return nil, fmt.Errorf("error in nb_count")
	}
	res = append(res, tx.nbOut...)

	for _, elt := range tx.outputs {
		if len(elt.Value) != 8 {
			return nil, fmt.Errorf("error in value")
		}

		res = append(res, elt.Value...)

		if len(elt.TxOutScriptLen) > 9 {
			return nil, fmt.Errorf("error in scriptlen")
		}

		res = append(res, elt.TxOutScriptLen...)
		res = append(res, elt.ScriptPubKey...)
	}

	return res, nil
}

/*func SignTransaction(privateKey string, sourcePkScript []byte, tx *wire.MsgTx) (string, error) {
	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "8", err
	}

	signatureScript, err := txscript.SignatureScript(tx, 0, sourcePkScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return "9", err
	}
	tx.TxIn[0].SignatureScript = signatureScript
	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)

	return hex.EncodeToString(signedTx.Bytes()), nil
}

func CreateTransaction(privateKey, destination string, amount int64, hash string) (string, error) {
	destinationAddress, err := btcutil.DecodeAddress(destination, &chaincfg.TestNet3Params)
	if err != nil {
		return "1", err
	}

	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return "2", err
	}

	publicAddr, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeCompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return "3", err
	}

	sourceAddress, err := btcutil.DecodeAddress(publicAddr.EncodeAddress(), &chaincfg.TestNet3Params)
	if err != nil {
		return "4", err
	}

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		return "5", err
	}

	destinationPkScript, err := txscript.PayToAddrScript(destinationAddress)
	if err != nil {
		return "6", err
	}

	// Create transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	utxoHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return "7", err
	}

	utxo := wire.NewOutPoint(utxoHash, 0)
	txIn := wire.NewTxIn(utxo, nil, nil)
	txOut := wire.NewTxOut(amount, destinationPkScript)

	tx.AddTxIn(txIn)
	tx.AddTxOut(txOut)

	finalTx, err := SignTransaction(privateKey, sourcePkScript, tx)
	return finalTx, err
}*/
