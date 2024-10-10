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
	sendbackAddress   = "mrs6r8TKaYZkXxrCw9kDg1C4XatTsss5Dm"
)

// NOTE: PubKey Address: mrs6r8TKaYZkXxrCw9kDg1C4XatTsss5Dm

func main() {
	godotenv.Load()

	secretString := os.Getenv("SECRET_STRING")
	secret := new(big.Int).SetBytes([]byte(secretString))
	fmt.Println("Secret:", secret.Text(16))
	privKey := privkey.NewPrivKey(secret)
	fmt.Println("WIF:", privKey.WIF(true, true))
	pubKey := privKey.PubKey()
	fmt.Printf("Pubkey:\n%s\n", hex.EncodeToString(pubKey.Serialize(true)))
	fmt.Println("Pubkey Address:", pubKey.Address(true, true))

	// NOTE: 使いたいトランザクションのID
	prevOutputHash, _ := hex.DecodeString("ec1728d31875b50e0f17f2e475eb43819d54b696ab8b114dbda029ed52a03941")
	prevOutputIndex := uint32(1)
	txIn := transaction.NewInput(prevOutputHash, prevOutputIndex, nil, 0xffffffff)

	scriptPubKey, err := script.NewP2PKHScriptPubkey(sendbackAddress)
	if err != nil {
		panic(err)
	}
	txOut := transaction.NewOutput(0.00015627*satoshiPerBitcoin, scriptPubKey)

	lockTime := uint32(0)

	tx := transaction.NewTransaction(1, []*transaction.Input{txIn}, []*transaction.Output{txOut}, lockTime, false)

	sigHash, err := tx.SigHash(0, true)
	if err != nil {
		panic(err)
	}
	z := new(big.Int).SetBytes(sigHash)
	sig := privKey.SignWithK(z, big.NewInt(1))
	serializedSig := sig.Serialize()
	serializedPubKey := pubKey.Serialize(true)
	scriptSig := script.NewScriptSig(serializedSig, serializedPubKey)
	tx.Inputs[0].ScriptSig = scriptSig

	serializedScriptSig, _ := scriptSig.Serialize()
	fmt.Printf("ScriptSig:\n%s\n", hex.EncodeToString(serializedScriptSig))
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
