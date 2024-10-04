package privkey

import (
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
