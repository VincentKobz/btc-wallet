package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/VincentKobz/btc-wallet/src"
)

func main() {
	btcAddress, privateKey := src.GenerateBtcAddress()
	fmt.Println(btcAddress)
	fmt.Println(privateKey)

	var test src.TxOut
	temp := make([]byte, 8)
	i := int64(1)
	binary.LittleEndian.PutUint64(temp, uint64(i))

	test.Value = temp

	transaction, err := src.Transaction("13iFeyezF7SHcHhpYsosTQpqWcxbL5aSCL", "4e8378675bcf6a389c8cfe246094aafa44249e48ab88a40e6fda3bf0f44f916a", 1, &test)
	if err != nil {
		fmt.Println("Error")
	}
	res, _ := transaction.Serialize()

	final := hex.EncodeToString(res)

	fmt.Println(final)
}
