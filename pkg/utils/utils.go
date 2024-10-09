package utils

import (
	"crypto/sha256"
	"encoding/binary"
	"io"

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

func CompareBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
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

func ParseVarInt(rader io.Reader) (uint64, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(rader, buf); err != nil {
		return 0, err
	}
	marker := buf[0]
	if marker == 0xfd {
		buf = make([]byte, 2)
		if _, err := io.ReadFull(rader, buf); err != nil {
			return 0, err
		}
		return uint64(binary.LittleEndian.Uint16(buf)), nil
	} else if marker == 0xfe {
		buf = make([]byte, 4)
		if _, err := io.ReadFull(rader, buf); err != nil {
			return 0, err
		}
		return uint64(binary.LittleEndian.Uint32(buf)), nil
	} else if marker == 0xff {
		buf = make([]byte, 8)
		if _, err := io.ReadFull(rader, buf); err != nil {
			return 0, err
		}
		return binary.LittleEndian.Uint64(buf), nil
	} else {
		return uint64(marker), nil
	}
}

func SerializeVarInt(n uint64) ([]byte, error) {
	if n < 0xfd {
		return []byte{byte(n)}, nil
	} else if n <= 0xffff {
		buf := make([]byte, 3)
		buf[0] = 0xfd
		binary.LittleEndian.PutUint16(buf[1:], uint16(n))
		return buf, nil
	} else if n <= 0xffffffff {
		buf := make([]byte, 5)
		buf[0] = 0xfe
		binary.LittleEndian.PutUint32(buf[1:], uint32(n))
		return buf, nil
	} else {
		buf := make([]byte, 9)
		buf[0] = 0xff
		binary.LittleEndian.PutUint64(buf[1:], n)
		return buf, nil
	}
}
