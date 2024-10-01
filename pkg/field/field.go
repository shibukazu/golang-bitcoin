package field

import (
	"math/big"
)

type FieldElement struct {
	Num   *big.Int
	Prime *big.Int
}

func NewFieldElement(num, prime *big.Int) FieldElement {
	if num.Cmp(big.NewInt(0)) < 0 || num.Cmp(prime) >= 0 {
		panic("Num is not in field range")
	}
	// コピーを作成してフィールド要素を初期化
	return FieldElement{new(big.Int).Set(num), new(big.Int).Set(prime)}
}

func (e FieldElement) Equals(other FieldElement) bool {
	return e.Num.Cmp(other.Num) == 0 && e.Prime.Cmp(other.Prime) == 0
}

func (e FieldElement) Add(other FieldElement) FieldElement {
	if e.Prime.Cmp(other.Prime) != 0 {
		panic("Cannot add two numbers in different Fields")
	}
	num := new(big.Int).Add(e.Num, other.Num)
	num.Mod(num, e.Prime)
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Subtract(other FieldElement) FieldElement {
	if e.Prime.Cmp(other.Prime) != 0 {
		panic("Cannot subtract two numbers in different Fields")
	}
	num := new(big.Int).Sub(e.Num, other.Num)
	num.Mod(num, e.Prime)
	if num.Cmp(big.NewInt(0)) < 0 {
		num.Add(num, e.Prime)
	}
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Multiply(other FieldElement) FieldElement {
	if e.Prime.Cmp(other.Prime) != 0 {
		panic("Cannot multiply two numbers in different Fields")
	}
	num := new(big.Int).Mul(e.Num, other.Num)
	num.Mod(num, e.Prime)
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Pow(exponent *big.Int) FieldElement {
	// Prime-1 を加えることで負の指数に対応
	n := new(big.Int).Mod(exponent, new(big.Int).Sub(e.Prime, big.NewInt(1)))
	num := new(big.Int).Exp(e.Num, n, e.Prime)
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Divide(other FieldElement) FieldElement {
	if e.Prime.Cmp(other.Prime) != 0 {
		panic("Cannot divide two numbers in different Fields")
	}
	// Fermatの小定理を使って逆数を計算し、それを乗算する
	inverse := new(big.Int).Exp(other.Num, new(big.Int).Sub(e.Prime, big.NewInt(2)), e.Prime)
	num := new(big.Int).Mul(e.Num, inverse)
	num.Mod(num, e.Prime)
	return FieldElement{num, e.Prime}
}
