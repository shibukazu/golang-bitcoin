package main

import (
	"encoding/hex"
	"golang-bitcoin/pkg/privkey"
	"golang-bitcoin/pkg/script"
	"golang-bitcoin/pkg/transaction"
	"math/big"
)

const (
	satoshiPerBitcoin = 100000000
	secretString      = "kazukishibutanitestnetprivekey"
)

func main() {
	// NOTE: 使いたいトランザクションのID
	prevOutputHash, _ := hex.DecodeString("")
	prevOutputIndex := uint32(0)
	txIn := transaction.NewInput(prevOutputHash, prevOutputIndex, nil, 0xffffffff)

	address := "mzBc4XEFS4g3v7m3UuZs4zr1vZ3f6z1j6C"
	scriptPubKey, err := script.NewP2PKHScriptPubkey(address)
	if err != nil {
		panic(err)
	}
	txOut := transaction.NewOutput(0.01*satoshiPerBitcoin, scriptPubKey)

	lockTime := uint32(0)

	tx := transaction.NewTransaction(1, []*transaction.Input{txIn}, []*transaction.Output{txOut}, lockTime)

	secret := new(big.Int).SetBytes([]byte(secretString))
	privKey := privkey.NewPrivKey(secret)
	pubKey := privKey.PubKey()
	serializedPubKey := pubKey.Serialize(true)
	sigHash, err := tx.SigHash(0, true)
	if err != nil {
		panic(err)
	}
	z := new(big.Int).SetBytes(sigHash)
	sig := privKey.Sign(z)
	serializedSig := sig.Serialize()
	scriptSig := script.NewScriptSig(serializedSig, serializedPubKey)
	tx.Inputs[0].ScriptSig = scriptSig

}
