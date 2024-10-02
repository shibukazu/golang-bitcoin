package secp256k1

import (
	"golang-bitcoin/pkg/curve"
	"golang-bitcoin/pkg/field"
	"golang-bitcoin/pkg/signature"
	"math/big"
)

const (
	s256aInt  = 0
	s256bInt  = 7
	s256pHex  = "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"
	s256gxHex = "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	s256gyHex = "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
)

type Secp256k1FieldElement struct {
	field.FieldElement
}

type Secp256k1Point struct {
	curve.Point
}

func NewSecp256k1Point(x, y *big.Int) Secp256k1Point {
	s256a := big.NewInt(s256aInt)
	s256b := big.NewInt(s256bInt)
	s256p, _ := new(big.Int).SetString(s256pHex, 16)

	point := curve.NewPoint(field.NewFieldElement(x, s256p), field.NewFieldElement(y, s256p), field.NewFieldElement(s256a, s256p), field.NewFieldElement(s256b, s256p))
	secP256k1Point := Secp256k1Point{point}

	return secP256k1Point
}

func (p Secp256k1Point) Verify(z *big.Int, sig signature.Signature) bool {
	s256p, _ := new(big.Int).SetString(s256pHex, 16)
	invS := new(big.Int).ModInverse(sig.S(), s256p)
	u := new(big.Int).Mul(z, invS)
	u.Mod(u, s256p)
	v := new(big.Int).Mul(sig.R(), invS)
	v.Mod(v, s256p)

	s256G := NewSecp256k1G()
	total := s256G.Multiply(u).Add(p.Multiply(v))
	return total.X().Cmp(sig.R()) == 0
}

func (p Secp256k1Point) Serialize(compressed bool) []byte {
	if !compressed {
		marker := byte(4)
		x := padTo32Bytes(p.X().Bytes())
		y := padTo32Bytes(p.Y().Bytes())
		serialized := make([]byte, 0, len(x)+len(y)+2)
		serialized = append(serialized, marker)
		serialized = append(serialized, x...)
		serialized = append(serialized, y...)

		return serialized
	} else {
		return nil
	}
}

func NewSecp256k1FieldElement(num *big.Int) Secp256k1FieldElement {
	s256p, _ := new(big.Int).SetString(s256pHex, 16)
	fieldElement := field.NewFieldElement(num, s256p)
	secp256k1FieldElement := Secp256k1FieldElement{fieldElement}

	return secp256k1FieldElement
}

func NewSecp256p() *big.Int {
	s256p, _ := new(big.Int).SetString(s256pHex, 16)

	return s256p
}

func NewSecp256k1G() Secp256k1Point {
	gx, _ := new(big.Int).SetString(s256gxHex, 16)
	gy, _ := new(big.Int).SetString(s256gyHex, 16)

	return NewSecp256k1Point(gx, gy)
}
