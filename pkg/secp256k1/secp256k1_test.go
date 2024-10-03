package secp256k1

import (
	"golang-bitcoin/pkg/curve"
	"math/big"
	"reflect"
	"testing"
)

func TestSecp256k1Point_Serialize(t *testing.T) {
	type fields struct {
		Point curve.Point
	}

	type args struct {
		compressed bool
	}
	tests := []struct {
		name            string
		fieldsGenerator func() fields
		args            args
		wantGenerator   func() []byte
	}{
		{
			name: "uncompressed case 1",
			fieldsGenerator: func() fields {
				e := big.NewInt(5000)
				P := NewSecp256k1G().Multiply(e)
				return fields{
					Point: P,
				}
			},
			args: args{
				compressed: false,
			},
			wantGenerator: func() []byte {
				wantHex := "04ffe558e388852f0120e46af2d1b370f85854a8eb0841811ece0e3e03d282d57c315dc72890a4f10a1481c031b03b351b0dc79901ca18a00cf009dbdb157a1d10"
				wantInt, _ := new(big.Int).SetString(wantHex, 16)

				return wantInt.Bytes()
			},
		},
		{
			name: "compressed case 1",
			fieldsGenerator: func() fields {
				e := big.NewInt(5001)
				P := NewSecp256k1G().Multiply(e)
				return fields{
					Point: P,
				}
			},
			args: args{
				compressed: true,
			},
			wantGenerator: func() []byte {
				wantHex := "0357a4f368868a8a6d572991e484e664810ff14c05c0fa023275251151fe0e53d1"
				wantInt, _ := new(big.Int).SetString(wantHex, 16)

				return wantInt.Bytes()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsGenerator()
			want := tt.wantGenerator()

			p := Secp256k1Point{
				Point: fields.Point,
			}
			if got := p.Serialize(tt.args.compressed); !reflect.DeepEqual(got, want) {
				t.Errorf("Secp256k1Point.SerializeUncompressedSEC() = %v, want %v", got, want)
			}
		})
	}
}

func TestDeserializeSecp256k1Point(t *testing.T) {
	type args struct {
		serialized []byte
	}
	tests := []struct {
		name          string
		argsGenerator func() args
		wantGenerator func() Secp256k1Point
	}{
		{
			name: "uncompressed case 1",
			argsGenerator: func() args {
				serializedHex := "04ffe558e388852f0120e46af2d1b370f85854a8eb0841811ece0e3e03d282d57c315dc72890a4f10a1481c031b03b351b0dc79901ca18a00cf009dbdb157a1d10"
				serializedInt, _ := new(big.Int).SetString(serializedHex, 16)

				return args{
					serialized: serializedInt.Bytes(),
				}
			},
			wantGenerator: func() Secp256k1Point {
				e := big.NewInt(5000)
				P := NewSecp256k1G().Multiply(e)
				return Secp256k1Point{
					Point: P,
				}
			},
		},
		{
			name: "compressed case 1",
			argsGenerator: func() args {
				serializedHex := "0357a4f368868a8a6d572991e484e664810ff14c05c0fa023275251151fe0e53d1"
				serializedInt, _ := new(big.Int).SetString(serializedHex, 16)

				return args{
					serialized: serializedInt.Bytes(),
				}
			},
			wantGenerator: func() Secp256k1Point {
				e := big.NewInt(5001)
				P := NewSecp256k1G().Multiply(e)
				return Secp256k1Point{
					Point: P,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.argsGenerator()
			want := tt.wantGenerator()
			if got := DeserializeSecp256k1Point(args.serialized); !got.Point.Equals(want.Point) {
				t.Errorf("DeserializeSecp256k1Point() = %v, want %v", got, want)
			}
		})
	}
}
