package curve

import (
	"golang-bitcoin/pkg/field"
	"math/big"
	"reflect"
	"testing"
)

func TestNewPoint(t *testing.T) {
	type args struct {
		x field.FieldElement
		y field.FieldElement
	}
	prime := big.NewInt(223)
	a := field.NewFieldElement(big.NewInt(0), prime)
	b := field.NewFieldElement(big.NewInt(7), prime)
	tests := []struct {
		name      string
		args      args
		wantPanic bool
	}{
		{
			name: "valid point",
			args: args{
				x: field.NewFieldElement(big.NewInt(192), prime),
				y: field.NewFieldElement(big.NewInt(105), prime),
			},
			wantPanic: false,
		},
		{
			name: "invalid point",
			args: args{
				x: field.NewFieldElement(big.NewInt(200), prime),
				y: field.NewFieldElement(big.NewInt(119), prime),
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewPoint() did not panic")
					}
				}()
			}
			NewPoint(tt.args.x, tt.args.y, a, b)
		})
	}
}

func TestPoint_Add(t *testing.T) {
	type fields struct {
		x field.FieldElement
		y field.FieldElement
	}
	type args struct {
		x field.FieldElement
		y field.FieldElement
	}
	prime := big.NewInt(223)
	a := field.NewFieldElement(big.NewInt(0), prime)
	b := field.NewFieldElement(big.NewInt(7), prime)
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Point
	}{
		{
			name: "Addition of two points on the same curve 1",
			fields: fields{
				x: field.NewFieldElement(big.NewInt(170), prime),
				y: field.NewFieldElement(big.NewInt(142), prime),
			},
			args: args{
				x: field.NewFieldElement(big.NewInt(60), prime),
				y: field.NewFieldElement(big.NewInt(139), prime),
			},
			want: Point{
				x: field.NewFieldElement(big.NewInt(220), prime),
				y: field.NewFieldElement(big.NewInt(181), prime),
				a: a,
				b: b,
			},
		},
		{
			name: "Addition of two points on the same curve 2",
			fields: fields{
				x: field.NewFieldElement(big.NewInt(47), prime),
				y: field.NewFieldElement(big.NewInt(71), prime),
			},
			args: args{
				x: field.NewFieldElement(big.NewInt(17), prime),
				y: field.NewFieldElement(big.NewInt(56), prime),
			},
			want: Point{
				x: field.NewFieldElement(big.NewInt(215), prime),
				y: field.NewFieldElement(big.NewInt(68), prime),
				a: a,
				b: b,
			},
		},
		{
			name: "Addition with the zero point",
			fields: fields{
				x: field.NewFieldElement(big.NewInt(47), prime),
				y: field.NewFieldElement(big.NewInt(71), prime),
			},
			args: args{
				x: field.NewFieldElement(big.NewInt(0), prime),
				y: field.NewFieldElement(big.NewInt(0), prime),
			},
			want: Point{
				x: field.NewFieldElement(big.NewInt(47), prime),
				y: field.NewFieldElement(big.NewInt(71), prime),
				a: a,
				b: b,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p1 := Point{
				x: tt.fields.x,
				y: tt.fields.y,
				a: a,
				b: b,
			}
			p2 := Point{
				x: tt.args.x,
				y: tt.args.y,
				a: a,
				b: b,
			}
			if got := p1.Add(p2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Point.Add() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPoint_Multiply(t *testing.T) {
	type fields struct {
		x field.FieldElement
		y field.FieldElement
		a field.FieldElement
		b field.FieldElement
	}
	type args struct {
		scalar *big.Int
	}
	type want struct {
		x field.FieldElement
		y field.FieldElement
	}
	tests := []struct {
		name          string
		fieldsFactory func() fields
		argsFactory   func() args
		wantFactory   func() want
	}{
		{
			name: "small 1",
			fieldsFactory: func() fields {
				prime := big.NewInt(223)
				a := big.NewInt(0)
				b := big.NewInt(7)
				x := big.NewInt(47)
				y := big.NewInt(71)
				return fields{
					x: field.NewFieldElement(x, prime),
					y: field.NewFieldElement(y, prime),
					a: field.NewFieldElement(a, prime),
					b: field.NewFieldElement(b, prime),
				}
			},
			argsFactory: func() args {
				scalar := big.NewInt(21)
				return args{
					scalar: scalar,
				}
			},
			wantFactory: func() want {
				prime := big.NewInt(223)
				x := big.NewInt(0)
				y := big.NewInt(0)
				return want{
					x: field.NewFieldElement(x, prime),
					y: field.NewFieldElement(y, prime),
				}
			},
		},
		{
			name: "small 2",
			fieldsFactory: func() fields {
				prime := big.NewInt(223)
				a := big.NewInt(0)
				b := big.NewInt(7)
				x := big.NewInt(192)
				y := big.NewInt(105)
				return fields{
					x: field.NewFieldElement(x, prime),
					y: field.NewFieldElement(y, prime),
					a: field.NewFieldElement(a, prime),
					b: field.NewFieldElement(b, prime),
				}
			},
			argsFactory: func() args {
				scalar := big.NewInt(2)
				return args{
					scalar: scalar,
				}
			},
			wantFactory: func() want {
				prime := big.NewInt(223)
				x := big.NewInt(49)
				y := big.NewInt(71)
				return want{
					x: field.NewFieldElement(x, prime),
					y: field.NewFieldElement(y, prime),
				}
			},
		},
		{
			name: "secp256k1",
			fieldsFactory: func() fields {
				primeHex := "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"
				prime := new(big.Int)
				prime.SetString(primeHex, 16)
				a := big.NewInt(0)
				b := big.NewInt(7)
				xHex := "79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"
				x := new(big.Int)
				x.SetString(xHex, 16)
				yHex := "483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8"
				y := new(big.Int)
				y.SetString(yHex, 16)
				return fields{
					x: field.NewFieldElement(x, prime),
					y: field.NewFieldElement(y, prime),
					a: field.NewFieldElement(a, prime),
					b: field.NewFieldElement(b, prime),
				}
			},
			argsFactory: func() args {
				scalarHex := "fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
				scalar := new(big.Int)
				scalar.SetString(scalarHex, 16)
				return args{
					scalar: scalar,
				}
			},
			wantFactory: func() want {
				primeHex := "fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"
				prime := new(big.Int)
				prime.SetString(primeHex, 16)
				return want{
					x: field.NewFieldElement(big.NewInt(0), prime),
					y: field.NewFieldElement(big.NewInt(0), prime),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsFactory()
			p1 := Point{
				x: fields.x,
				y: fields.y,
				a: fields.a,
				b: fields.b,
			}
			args := tt.argsFactory()
			want := tt.wantFactory()
			wantPoint := Point{
				x: want.x,
				y: want.y,
				a: fields.a,
				b: fields.b,
			}
			if got := p1.Multiply(args.scalar); !got.Equals(wantPoint) {
				t.Errorf("Point.Multiply() = %v, want %v", got, want)
			}
		})
	}
}
