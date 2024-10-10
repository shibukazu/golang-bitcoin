package secp256k1

import (
	"fmt"
	"golang-bitcoin/pkg/curve"
	"golang-bitcoin/pkg/field"
	"golang-bitcoin/pkg/signature"
	"golang-bitcoin/pkg/utils"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
)

const (
	s256aInt  = 0
	s256bInt  = 7
	s256pHex  = "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"
	s256gxHex = "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
	s256gyHex = "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
	s256nHex  = "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
)

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
	s256n := NewSecp256k1n()
	invS := new(big.Int).ModInverse(sig.S(), s256n)
	u := new(big.Int).Mul(z, invS)
	u.Mod(u, s256n)
	v := new(big.Int).Mul(sig.R(), invS)
	v.Mod(v, s256n)

	s256G := NewSecp256k1G()
	total := s256G.Multiply(u).Add(p.Multiply(v))
	return total.X().Cmp(sig.R()) == 0
}

func (p Secp256k1Point) Serialize(compressed bool) []byte {
	if !compressed {
		marker := byte(4)
		x := utils.PadTo32Bytes(p.X().Bytes())
		y := utils.PadTo32Bytes(p.Y().Bytes())
		serialized := make([]byte, 0, len(x)+len(y)+2)
		serialized = append(serialized, marker)
		serialized = append(serialized, x...)
		serialized = append(serialized, y...)

		return serialized
	} else {
		isEven := p.Y().Bit(0) == 0

		var marker byte
		if isEven {
			marker = byte(2)
		} else {
			marker = byte(3)
		}

		x := utils.PadTo32Bytes(p.X().Bytes())
		serialized := make([]byte, 0, len(x)+1)
		serialized = append(serialized, marker)
		serialized = append(serialized, x...)

		return serialized
	}
}

func ParseSecp256k1Point(serialized []byte) Secp256k1Point {
	marker := serialized[0]
	x := new(big.Int).SetBytes(serialized[1:33])
	var y *big.Int
	if marker == 4 {
		y = new(big.Int).SetBytes(serialized[33:])
	} else {
		right := new(big.Int).Exp(x, big.NewInt(3), NewSecp256p())
		right = right.Add(right, big.NewInt(7))
		right = right.Mod(right, NewSecp256p())

		y = new(big.Int).ModSqrt(right, NewSecp256p())

		var even_y *big.Int
		var odd_y *big.Int
		s256p := NewSecp256p()
		if y.Bit(0) == 0 {
			even_y = y
			odd_y = new(big.Int).Sub(s256p, y)
		} else {
			even_y = new(big.Int).Sub(s256p, y)
			odd_y = y
		}

		isEven := marker == 2
		if isEven {
			y = even_y
		} else {
			y = odd_y
		}
	}

	return NewSecp256k1Point(x, y)
}

func (p Secp256k1Point) Address(compressed bool, testnet bool) string {
	serialized := p.Serialize(compressed)

	serialized160 := utils.Hash160(serialized)

	var prefix []byte
	if testnet {
		prefix = []byte{0x6f}
	} else {
		prefix = []byte{0x00}
	}

	joint := append(prefix, serialized160...)
	checksum := utils.Hash256(joint)[:4]

	return base58.Encode(append(joint, checksum...))
}

func ExtractHash160(address string) ([]byte, error) {
	decoded := base58.Decode(address)
	if len(decoded) != 25 {
		return nil, fmt.Errorf("invalid address length")
	}

	prefix := decoded[0]
	serialized160 := decoded[1:21]
	checksum := decoded[21:]

	joint := append([]byte{prefix}, serialized160...)
	rawChecksum := utils.Hash256(joint)[:4]
	if !utils.CompareBytes(checksum, rawChecksum) {
		return nil, fmt.Errorf("invalid checksum")
	}

	return serialized160, nil
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

func NewSecp256k1n() *big.Int {
	s256n, _ := new(big.Int).SetString(s256nHex, 16)

	return s256n
}

func NewSecp256k1nHalf() *big.Int {
	s256n := NewSecp256k1n()
	s256nHalf := new(big.Int).Rsh(s256n, 1)

	return s256nHalf
}
