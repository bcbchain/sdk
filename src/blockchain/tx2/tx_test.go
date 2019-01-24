package tx2

import (
	"blockchain/smcsdk/sdk/rlp"
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"blockchain/algorithm"
	"blockchain/smcsdk/sdk/bn"
	"blockchain/types"

	"github.com/tendermint/go-crypto"
)

func TestTransaction_TxParse(t *testing.T) {
	Init("bcb")
	crypto.SetChainId("bcb")

	methodID1 := algorithm.BytesToUint32(algorithm.CalcMethodId("Transfer(types.Address,bn.Number)"))
	toContract1 := "bcbMWedWqzzW8jkt5tntTomQQEN7fSwWFhw6"

	toAccount := "bcbCpeczqoSoxLxx1x3UyuKsaS4J8yamzWzz"
	value := bn.N(1000000000)
	itemInfo1 := WrapInvokeParams(toAccount, value)
	message1 := types.Message{
		Contract: toContract1,
		MethodID: methodID1,
		Items:    itemInfo1,
	}
	nonce := uint64(1)
	gasLimit := int64(500)
	note := "Example for cascade invoke smart contract."
	txPayloadBytesRlp := WrapPayload(nonce, gasLimit, note, message1)
	privKeyStr := "0x4a2c14697282e658b3ed7dd5324de1a102d216d6fa50d5937ffe89f35cbc12aa68eb9a09813bdf7c0869bf34a244cc545711509fe70f978d121afd3a4ae610e6"
	finalTx := WrapTx(txPayloadBytesRlp, privKeyStr)

	privKeyBytes, _ := hex.DecodeString(privKeyStr[2:])
	privKey := crypto.PrivKeyEd25519FromBytes(privKeyBytes)
	address := privKey.PubKey().Address()

	transaction, txPubKey, err := TxParse(string(finalTx))
	if err != nil {
		panic("TxParse error:" + err.Error())
	}
	if !reflect.DeepEqual(txPubKey.Address(), address) {
		fmt.Println("Parsed pubkey: ", txPubKey)
		fmt.Println("Wrapper pubkey:", privKey.PubKey())
		fmt.Println("Parsed pubkey address: ", txPubKey.Address())
		fmt.Println("Wrapper pubkey address:", address)

		panic("Sender address is wrong")
	}
	if !reflect.DeepEqual(transaction.Nonce, nonce) {
		panic("nonce is wrong")
	}
	if !reflect.DeepEqual(transaction.GasLimit, gasLimit) {
		panic("gaslimit is wrong")
	}

	if !reflect.DeepEqual(len(transaction.Messages), int(1)) {
		panic("Message size mismatch")
	}
	message := transaction.Messages[0]
	if !reflect.DeepEqual(message.Contract, message1.Contract) {
		panic("Contract in message mismatch")
	}
	if !reflect.DeepEqual(message.MethodID, message1.MethodID) {
		panic("MethodID in message mismatch")
	}
	if !reflect.DeepEqual(len(message.Items), len(message1.Items)) {
		panic("Length of items in message mismatch")
	}
	if !reflect.DeepEqual(len(message.Items), int(2)) {
		panic("Length of items in message is wrong")
	}
	to := ""
	err = rlp.DecodeBytes(message.Items[0], &to)
	if err != nil {
		panic(err)
	}
	if to != toAccount {
		panic("item to is wrong")
	}
	var v bn.Number
	err = rlp.DecodeBytes(message.Items[1], &v)
	if err != nil {
		panic(err)
	}
	if v.String() != value.String() {
		panic("item value is wrong")
	}
}
