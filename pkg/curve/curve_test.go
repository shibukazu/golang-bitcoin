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
