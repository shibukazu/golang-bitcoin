package script

import "math/big"

// NOTE: 数値の符号を考慮し、little-endianでエンコード
func encodeNum(num int64) []byte {
	if num == 0 {
		return []byte{}
	}

	absNum := big.NewInt(num)
	negative := num < 0
	result := []byte{}

	// NOTE: convert to little-endian
	for absNum.Cmp(big.NewInt(0)) != 0 {
		result = append(result, byte(absNum.Int64()&0xff))
		absNum.Rsh(absNum, 8)
	}

	// NOTE: 正値と負値の判定のために明示的な符号ビットを追加
	if result[len(result)-1]&0x80 != 0 {
		if negative {
			result = append(result, 0x80)
		} else {
			result = append(result, 0)
		}
	} else if negative {
		result[len(result)-1] |= 0x80
	}

	return result
}

func decodeNum(element []byte) int64 {
	if len(element) == 0 {
		return 0
	}

	bigEndian := make([]byte, len(element))
	for i := 0; i < len(element); i++ {
		bigEndian[i] = element[len(element)-1-i]
	}
	negative := bigEndian[0]&0x80 != 0
	// NOTE: 上記より先頭ビットは符号の意味しかないため、無視できる
	//       その場合、以降のビットが絶対値を表す
	result := big.NewInt(int64(bigEndian[0] & 0x7f))

	for _, c := range bigEndian[1:] {
		result.Lsh(result, 8)
		result.Add(result, big.NewInt(int64(c)))
	}

	if negative {
		result.Neg(result)
	}
	return result.Int64()
}
