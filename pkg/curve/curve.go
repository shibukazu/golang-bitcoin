package curve

import (
	"fmt"
	"golang-bitcoin/pkg/field"
)

type Point struct {
	x, y, a, b field.FieldElement
	isInf      bool
}

func NewPoint(x, y, a, b field.FieldElement, isInf bool) Point {
	left := y.Pow(2)
	right := x.Pow(3).Add(a.Multiply(x)).Add(b)
	if !isInf && !left.Equals(right) {
		panic("Point not on curve")
	}
	return Point{x, y, a, b, isInf}
}

func (p1 Point) IsOnSameCurve(p2 Point) bool {
	return p1.a.Equals(p2.a) && p1.b.Equals(p2.b)
}

func (p1 Point) Equals(p2 Point) bool {
	if !p1.IsOnSameCurve(p2) {
		return false
	} else {
		if p1.isInf && p2.isInf {
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
		if p1.isInf {
			return p2
		}
		if p2.isInf {
			return p1
		}
		// NOTE: 加法逆元との加算
		if p1.x.Equals(p2.x) && !p1.y.Equals(p2.y) {
			return NewPoint(p1.x, p1.y, p1.a, p1.b, true)
		}
		// NOTE: 異なる点同士の加算
		if !p1.x.Equals(p2.x) {
			s := (p2.y.Subtract(p1.y)).Divide(p2.x.Subtract(p1.x))
			x := s.Pow(2).Subtract(p1.x).Subtract(p2.x)
			y := s.Multiply(p1.x.Subtract(x)).Subtract(p1.y)
			fmt.Println(s, x, y, p1.a, p1.b)
			return NewPoint(x, y, p1.a, p1.b, false)
		}
		zero := field.NewFieldElement(0, p1.a.Prime)
		// NOTE: 同じ点同士でyがゼロでない場合の加算
		if p1.y.Equals(p2.y) && !p1.y.Equals(zero) {
			three := field.NewFieldElement(3, p1.a.Prime)
			two := field.NewFieldElement(2, p1.a.Prime)
			s := p1.x.Pow(2).Multiply(three).Add(p1.a).Divide(two.Multiply(p1.y))
			x := s.Pow(2).Subtract(two.Multiply(p1.x))
			y := s.Multiply(p1.x.Subtract(x)).Subtract(p1.y)

			return NewPoint(x, y, p1.a, p1.b, false)
		}
		// NOTE: 同じ点同士でyがゼロの場合の加算
		return NewPoint(zero, zero, p1.a, p1.b, true)
	}

}
