package curve

import (
	"golang-bitcoin/pkg/field"
	"math/big"
)

type Point struct {
	x, y, a, b field.FieldElement
}

func NewPoint(x, y, a, b field.FieldElement) Point {
	isInf := x.Num.Cmp(big.NewInt(0)) == 0 && y.Num.Cmp(big.NewInt(0)) == 0
	if !isInf {
		left := y.Pow(big.NewInt(2))
		right := x.Pow(big.NewInt(3)).Add(a.Multiply(x)).Add(b)
		if  !left.Equals(right) {
			panic("Point not on curve")
		}
	}
		
	return Point{x, y, a, b}
}

func (p1 Point) IsInf() bool {
	return p1.x.Num.Cmp(big.NewInt(0)) == 0 && p1.y.Num.Cmp(big.NewInt(0)) == 0
}

func (p1 Point) IsOnSameCurve(p2 Point) bool {
	return p1.a.Equals(p2.a) && p1.b.Equals(p2.b)
}

func (p1 Point) Equals(p2 Point) bool {
	if !p1.IsOnSameCurve(p2) {
		return false
	} else {
		if p1.IsInf() && p2.IsInf() {
			return true
		} else {
			return p1.x.Equals(p2.x) && p1.y.Equals(p2.y)
		}
	}
}

func (p1 Point) Add(p2 Point) Point {
	if !p1.IsOnSameCurve(p2) {
		panic("Points are not on the same curve")
	} else {
		// NOTE: 加法単位元との加算
		if p1.IsInf() {
			return p2
		}
		if p2.IsInf() {
			return p1
		}
		// NOTE: 加法逆元との加算
		if p1.x.Equals(p2.x) && !p1.y.Equals(p2.y) {
			return NewPoint(field.NewFieldElement(big.NewInt(0), p1.a.Prime), field.NewFieldElement(big.NewInt(0), p1.a.Prime), p1.a, p1.b)
		}
		// NOTE: 異なる点同士の加算
		if !p1.x.Equals(p2.x) {
			s := p2.y.Subtract(p1.y).Divide(p2.x.Subtract(p1.x))
			x := s.Pow(big.NewInt(2)).Subtract(p1.x).Subtract(p2.x)
			y := s.Multiply(p1.x.Subtract(x)).Subtract(p1.y)
			return NewPoint(x, y, p1.a, p1.b)
		}
		zero := field.NewFieldElement(big.NewInt(0), p1.a.Prime)
		// NOTE: 同じ点同士でyがゼロでない場合の加算
		if p1.y.Equals(p2.y) && !p1.y.Equals(zero) {
			three := field.NewFieldElement(big.NewInt(3), p1.a.Prime)
			two := field.NewFieldElement(big.NewInt(2), p1.a.Prime)
			s := p1.x.Pow(big.NewInt(2)).Multiply(three).Add(p1.a).Divide(two.Multiply(p1.y))
			x := s.Pow(big.NewInt(2)).Subtract(two.Multiply(p1.x))
			y := s.Multiply(p1.x.Subtract(x)).Subtract(p1.y)

			return NewPoint(x, y, p1.a, p1.b)
		}
		// NOTE: 同じ点同士でyがゼロの場合の加算
		return NewPoint(zero, zero, p1.a, p1.b)
	}
}

func (p1 Point) Multiply(scalar *big.Int) Point {
	zero := big.NewInt(0)
	if scalar.Cmp(zero) == 0 {
		return NewPoint(field.NewFieldElement(zero, p1.a.Prime), field.NewFieldElement(zero, p1.a.Prime), p1.a, p1.b)
	}
	if scalar.Cmp(zero) < 0 {
		panic("Scalar must be a positive integer")
	}
	result := NewPoint(field.NewFieldElement(zero, p1.a.Prime), field.NewFieldElement(zero, p1.a.Prime), p1.a, p1.b)
	// TODO: 二進展開を理解する
	addend := p1
	for scalar.Cmp(zero) > 0 {
		if new(big.Int).And(scalar, big.NewInt(1)).Cmp(big.NewInt(1)) == 0 {
			result = result.Add(addend)
		}
		addend = addend.Add(addend)
		scalar.Rsh(scalar, 1)
	}
	return result
}
