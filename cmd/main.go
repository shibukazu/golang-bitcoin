package main

import (
	"fmt"
	"golang-bitcoin/pkg/secp256k1"
	"math/big"
)

func main() {
	e := new(big.Int).SetBytes([]byte("kazukishibutanitestnetprivekey"))
	G := secp256k1.NewSecp256k1G()
	P := G.Multiply(e)

	pubkey := secp256k1.NewSecp256k1Point(P.X(), P.Y())
	serialized := pubkey.Address(true, true)
	fmt.Println(serialized)
}
