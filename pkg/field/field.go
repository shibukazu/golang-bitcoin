package field

type FieldElement struct {
	Num   int
	Prime int
}

func NewFieldElement(num, prime int) FieldElement {
	if num >= prime || num < 0 {
		panic("Num is not in field range")
	}
	return FieldElement{num, prime}
}

func (e FieldElement) Equals(other FieldElement) bool {
	return e.Num == other.Num && e.Prime == other.Prime
}

func (e FieldElement) Add(other FieldElement) FieldElement {
	if e.Prime != other.Prime {
		panic("Cannot add two numbers in different Fields")
	}
	num := (e.Num + other.Num) % e.Prime
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Subtract(other FieldElement) FieldElement {
	if e.Prime != other.Prime {
		panic("Cannot subtract two numbers in different Fields")
	}
	num := (e.Num - other.Num) % e.Prime
	if num < 0 {
		num += e.Prime
	}
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Multiply(other FieldElement) FieldElement {
	if e.Prime != other.Prime {
		panic("Cannot multiply two numbers in different Fields")
	}
	num := (e.Num * other.Num) % e.Prime
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Pow(exponent int) FieldElement {
	// NOTE: 負の値への対応
	n := exponent
	for n < 0 {
		n += e.Prime - 1
	}
	num := 1
	for i := 0; i < n; i++ {
		num = (num * e.Num) % e.Prime
	}
	return FieldElement{num, e.Prime}
}

func (e FieldElement) Divide(other FieldElement) FieldElement {
	if e.Prime != other.Prime {
		panic("Cannot divide two numbers in different Fields")
	}
	num := (e.Num * other.Pow(e.Prime-2).Num) % e.Prime
	return FieldElement{num, e.Prime}
}
