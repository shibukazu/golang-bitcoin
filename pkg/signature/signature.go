package signature

import (
	"fmt"
	"math/big"
)

type Signature struct {
	r, s *big.Int
}

func NewSignature(r, s *big.Int) *Signature {
	return &Signature{r, s}
}

func (s *Signature) R() *big.Int {
	return s.r
}

func (s *Signature) S() *big.Int {
	return s.s
}

func (s *Signature) Equals(other *Signature) bool {
	return s.r.Cmp(other.r) == 0 && s.s.Cmp(other.s) == 0
}

func (s *Signature) Serialize() []byte {
	marker := []byte{0x30}
	rbin := s.r.Bytes()
	for len(rbin) > 0 && rbin[0] == 0 {
		rbin = rbin[1:]
	}
	if rbin[0]&0x80 != 0 {
		rbin = append([]byte{0x00}, rbin...)
	}
	rbin = append([]byte{0x02, byte(len(rbin))}, rbin...)

	sbin := s.s.Bytes()
	for len(sbin) > 0 && sbin[0] == 0 {
		sbin = sbin[1:]
	}
	if sbin[0]&0x80 != 0 {
		sbin = append([]byte{0x00}, sbin...)
	}
	sbin = append([]byte{0x02, byte(len(sbin))}, sbin...)

	remainLen := len(rbin) + len(sbin)

	return append(append(append(marker, byte(remainLen)), rbin...), sbin...)
}

func ParseSignature(signature []byte) (*Signature, error) {
	if len(signature) < 6 {
		return nil, fmt.Errorf("signature is too short")
	}

	if signature[0] != 0x30 {
		return nil, fmt.Errorf("invalid der marker")
	}

	if int(signature[1]) != len(signature)-2 {
		return nil, fmt.Errorf("invalid der length: %d != %d", signature[1], len(signature)-2)
	}

	if signature[2] != 0x02 {
		return nil, fmt.Errorf("invalud r marker")
	}

	rLen := int(signature[3])
	if rLen == 0 || rLen > 33 {
		return nil, fmt.Errorf("invalid r length")
	}

	if int(signature[4+rLen]) != 0x02 {
		return nil, fmt.Errorf("invalid s marker")
	}

	sLen := int(signature[3+rLen+3])
	if sLen == 0 || sLen > 33 {
		return nil, fmt.Errorf("invalid s lemgth")
	}

	r := new(big.Int).SetBytes(signature[4 : 4+rLen])
	s := new(big.Int).SetBytes(signature[4+rLen+2:])

	return NewSignature(r, s), nil
}
