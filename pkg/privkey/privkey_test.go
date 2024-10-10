package privkey

import (
	"fmt"
	"golang-bitcoin/pkg/secp256k1"
	"golang-bitcoin/pkg/signature"
	"golang-bitcoin/pkg/utils"
	"math/big"
	"testing"
)

func TestPrivKey_WIF(t *testing.T) {
	type fields struct {
		secret *big.Int
	}
	type args struct {
		compressed bool
		testnet    bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "testnet compressed",
			fields: fields{
				secret: big.NewInt(5003),
			},
			args: args{
				compressed: true,
				testnet:    true,
			},
			want: "cMahea7zqjxrtgAbB7LSGbcQUr1uX1ojuat9jZodMN8rFTv2sfUK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PrivKey{
				secret: tt.fields.secret,
			}
			if got := p.WIF(tt.args.compressed, tt.args.testnet); got != tt.want {
				t.Errorf("PrivKey.WIF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivKey_SignWithK(t *testing.T) {
	type fields struct {
		secret *big.Int
	}
	type args struct {
		z *big.Int
		k *big.Int
	}
	tests := []struct {
		name            string
		fieldsGenerator func() fields
		argsGenerator   func() args
		wantGenerator   func() *signature.Signature
	}{
		{
			name: "valid signature 1",
			fieldsGenerator: func() fields {
				secretString := "my secret"
				secretHash := utils.Hash256([]byte(secretString))
				secret := new(big.Int).SetBytes(secretHash)

				fmt.Printf("Public key: %s\n", secp256k1.NewSecp256k1G().Multiply(secret).X().Text(16))

				return fields{
					secret: secret,
				}
			},
			argsGenerator: func() args {
				zString := "my message"
				zHash := utils.Hash256([]byte(zString))
				z := new(big.Int).SetBytes(zHash)
				k := big.NewInt(1234567890)
				return args{
					z: z,
					k: k,
				}
			},
			wantGenerator: func() *signature.Signature {
				rHex := "2b698a0f0a4041b77e63488ad48c23e8e8838dd1fb7520408b121697b782ef22"
				r, _ := new(big.Int).SetString(rHex, 16)
				sHex := "bb14e602ef9e3f872e25fad328466b34e6734b7a0fcd58b1eb635447ffae8cb9"
				s, _ := new(big.Int).SetString(sHex, 16)
				return signature.NewSignature(r, s)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fieldsGenerator()
			args := tt.argsGenerator()
			want := tt.wantGenerator()
			p := PrivKey{
				secret: fields.secret,
			}
			if got := p.SignWithK(args.z, args.k); !got.Equals(want) {
				fmt.Printf("Want r: %s\n", want.R().Text(16))
				fmt.Printf("Got  r: %s\n", got.R().Text(16))
				fmt.Printf("Want s: %s\n", want.S().Text(16))
				fmt.Printf("Got  s: %s\n", got.S().Text(16))
				t.Errorf("PrivKey.SignWithK() = %v, want %v", got, want)
			}
		})
	}
}
