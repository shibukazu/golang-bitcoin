package privkey

import (
	"crypto/rand"
	"golang-bitcoin/pkg/secp256k1"
	"math/big"
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
