package privkey

import (
	"crypto/rand"
	"golang-bitcoin/pkg/secp256k1"
	"golang-bitcoin/pkg/signature"
	"golang-bitcoin/pkg/utils"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
)

type PrivKey struct {
	secret *big.Int
}

func NewPrivKey(secret *big.Int) PrivKey {
	return PrivKey{secret}
}

func (p PrivKey) Secret() *big.Int {
	return p.secret
}

func (p PrivKey) Equals(other PrivKey) bool {
	return p.secret.Cmp(other.secret) == 0
}

func (p PrivKey) Sign(z *big.Int) *signature.Signature {
	var r *big.Int
	var k *big.Int
	var err error
	for r == nil || r.Sign() == 0 {
		k, err = rand.Int(rand.Reader, secp256k1.NewSecp256k1n())
		if err != nil {
			panic(err)
		}
		R := secp256k1.NewSecp256k1G().Multiply(k)
		r = R.X()
		r.Mod(r, secp256k1.NewSecp256k1n())
	}

	invK := new(big.Int).ModInverse(k, secp256k1.NewSecp256k1n())

	re := new(big.Int).Mul(p.secret, r)
	re.Mod(re, secp256k1.NewSecp256k1n())

	rez := new(big.Int).Add(re, z)
	rez.Mod(rez, secp256k1.NewSecp256k1n())

	s := new(big.Int).Mul(rez, invK)
	s.Mod(s, secp256k1.NewSecp256k1n())

	if s.Cmp(secp256k1.NewSecp256k1nHalf()) == 1 {
		s = new(big.Int).Sub(secp256k1.NewSecp256k1n(), s)
	}

	return signature.NewSignature(r, s)
}

func (p PrivKey) SignWithK(z *big.Int, k *big.Int) *signature.Signature {
	R := secp256k1.NewSecp256k1G().Multiply(k)
	r := R.X()

	// rが0の場合は例外処理が必要
	if r.Sign() == 0 {
		panic("r is 0, invalid signature")
	}

	// kの逆数を計算
	invK := new(big.Int).ModInverse(k, secp256k1.NewSecp256k1n())

	// rez = (p.secret * r + z) mod p
	rez := new(big.Int).Add(new(big.Int).Mul(p.secret, r), z)
	rez.Mod(rez, secp256k1.NewSecp256k1n())

	// s = rez * invK mod p
	s := new(big.Int).Mul(rez, invK)
	s.Mod(s, secp256k1.NewSecp256k1n())

	if s.Cmp(secp256k1.NewSecp256k1nHalf()) == 1 {
		s = new(big.Int).Sub(secp256k1.NewSecp256k1n(), s)
	}

	// 署名を生成
	return signature.NewSignature(r, s)
}

func (p PrivKey) WIF(compressed bool, testnet bool) string {
	secretBytes := utils.PadTo32Bytes(p.secret.Bytes())
	var prefix []byte
	if testnet {
		prefix = []byte{0xef}
	} else {
		prefix = []byte{0x80}
	}
	secretBytes = append(prefix, secretBytes...)
	if compressed {
		secretBytes = append(secretBytes, 0x01)
	}
	checksum := utils.Hash256(secretBytes)
	checksum = checksum[:4]
	secretBytes = append(secretBytes, checksum...)
	return base58.Encode(secretBytes)
}

func (p *PrivKey) PubKey() secp256k1.Secp256k1Point {
	P := secp256k1.NewSecp256k1G().Multiply(p.secret)
	return secp256k1.NewSecp256k1Point(P.X(), P.Y())
}
