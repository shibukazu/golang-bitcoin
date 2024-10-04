package utils

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

func PadTo32Bytes(input []byte) []byte {
	if len(input) == 32 {
		return input
	}

	padded := make([]byte, 32)
	copy(padded[32-len(input):], input)
	return padded
}

func Hash160(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)

	h2 := ripemd160.New()
	h2.Write(hashed)
	return h2.Sum(nil)
}

func Hash256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)

	h.Reset()
	h = sha256.New()
	h.Write(hashed)
	return h.Sum(nil)
}
