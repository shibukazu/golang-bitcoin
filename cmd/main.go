package main

import (
	"encoding/hex"
	"fmt"
	"golang-bitcoin/pkg/privkey"
	"golang-bitcoin/pkg/script"
	"golang-bitcoin/pkg/transaction"
	"math/big"
	"os"

	"github.com/joho/godotenv"
)

const (
	satoshiPerBitcoin = 100000000
	sendbackAddress   = "msijx6rX4HcwPrFQ5gPf8A9d9BkEKCZo5H"
)

// NOTE: PubKey Address: msijx6rX4HcwPrFQ5gPf8A9d9BkEKCZo5H

func main() {
	godotenv.Load()

	secretString := os.Getenv("SECRET_STRING")
	secret := new(big.Int).SetBytes([]byte(secretString))
	privKey := privkey.NewPrivKey(secret)
	pubKey := privKey.PubKey()
	fmt.Println("Pubkey Address:", pubKey.Address(true, true))

	// NOTE: 使いたいトランザクションのID
	prevOutputHash, _ := hex.DecodeString("75245e7b859c3cbd17bebea6bf691ca4bb646a8d1a951b8da060185cec7b5c6d")
	prevOutputIndex := uint32(1)
	txIn := transaction.NewInput(prevOutputHash, prevOutputIndex, nil, 0xffffffff)

	scriptPubKey, err := script.NewP2PKHScriptPubkey(sendbackAddress)
	if err != nil {
		panic(err)
	}
	txOut := transaction.NewOutput(0.00000905*satoshiPerBitcoin, scriptPubKey)

	lockTime := uint32(0)

	tx := transaction.NewTransaction(1, []*transaction.Input{txIn}, []*transaction.Output{txOut}, lockTime, false)

	sigHash, err := tx.SigHash(0, true)
	if err != nil {
		panic(err)
	}
	z := new(big.Int).SetBytes(sigHash)
	sig := privKey.Sign(z)
	serializedSig := sig.Serialize()
	serializedPubKey := pubKey.Serialize(true)
	scriptSig := script.NewScriptSig(serializedSig, serializedPubKey)
	tx.Inputs[0].ScriptSig = scriptSig

	serialized, err := tx.Serialize()
	if err != nil {
		panic(err)
	}

	txID, err := tx.ID()
	if err != nil {
		panic(err)
	}
	fmt.Println("TransactionID: ", txID)
	fmt.Printf("Transaction:\n%s\n", hex.EncodeToString(serialized))
}
