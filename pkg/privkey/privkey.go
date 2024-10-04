package privkey

import (
	"crypto/rand"
	"golang-bitcoin/pkg/secp256k1"
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

func (p PrivKey) Sign(z *big.Int) (r, s *big.Int) {
	k, err := rand.Int(rand.Reader, nil)
	if err != nil {
		panic(err)
	}
	R := secp256k1.NewSecp256k1G().Multiply(k)
	r = R.X()
	invK := new(big.Int).ModInverse(k, secp256k1.NewSecp256p())
	re := new(big.Int).Mul(p.secret, r)
	re.Add(re, z)
	s = new(big.Int).Mul(re, invK)
	s.Mod(s, secp256k1.NewSecp256p())
	return r, s
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
