package signature

import "math/big"

type Signature struct {
	r, s *big.Int
}

func NewSignature(r, s *big.Int) Signature {
	return Signature{r, s}
}

func (s Signature) R() *big.Int {
	return s.r
}

func (s Signature) S() *big.Int {
	return s.s
}

func (s Signature) Equals(other Signature) bool {
	return s.r.Cmp(other.r) == 0 && s.s.Cmp(other.s) == 0
}
