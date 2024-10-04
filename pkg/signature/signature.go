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

func (s Signature) Seriarize() []byte {
	marker := []byte{0x30}
	rbin := s.r.Bytes()
	for len(rbin) > 0 && rbin[0] == 0 {
		rbin = rbin[1:]
	}
	if rbin[0]&0x80 != 0 {
		rbin = append([]byte{0x00}, rbin...)
	}
	rbin = append([]byte{0x02}, rbin...)

	sbin := s.s.Bytes()
	for len(sbin) > 0 && sbin[0] == 0 {
		sbin = sbin[1:]
	}
	if sbin[0]&0x80 != 0 {
		sbin = append([]byte{0x00}, sbin...)
	}
	sbin = append([]byte{0x02}, sbin...)

	return append(marker, append(rbin, sbin...)...)
}
