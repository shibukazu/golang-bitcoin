package signature

import (
	"encoding/hex"
	"math/big"
	"reflect"
	"testing"
)

func TestSignature_Serialize(t *testing.T) {
	type fields struct {
		r *big.Int
		s *big.Int
	}
	tests := []struct {
		name            string
		fieldsGenerator func() fields
		wantGenerator   func() []byte
	}{
		{
			name: "valid 1",
			fieldsGenerator: func() fields {
				r, _ := new(big.Int).SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
				s, _ := new(big.Int).SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
				return fields{
					r: r,
					s: s,
				}
			},
			wantGenerator: func() []byte {
				want, _ := hex.DecodeString("3045022037206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c60221008ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec")
				return want
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsGenerator()
			want := tt.wantGenerator()
			s := &Signature{
				r: fields.r,
				s: fields.s,
			}
			if got := s.Serialize(); !reflect.DeepEqual(got, want) {
				t.Errorf("Signature.Serialize() = %x, want %x", got, want)
			}
		})
	}
}
