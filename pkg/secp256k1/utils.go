package secp256k1

func padTo32Bytes(input []byte) []byte {
	if len(input) == 32 {
		return input
	}

	padded := make([]byte, 32)
	copy(padded[32-len(input):], input)
	return padded
}
