package script

import (
	"encoding/binary"
	"fmt"
	"golang-bitcoin/pkg/utils"
	"io"
)

const (
	OP_DUP       = 0x76
	OP_HASH160   = 0xa9
	OP_HASH256   = 0xaa
	OP_PUSHDATA1 = 0x4c
	OP_PUSHDATA2 = 0x4d
)

type Stack [][]byte

type Script struct {
	Instructions Stack
}

func NewScript() *Script {
	return &Script{
		Instructions: make(Stack, 0),
	}
}

func OpDup(stack Stack) Stack {
	if len(stack) < 1 {
		return stack
	}
	stack = append(stack, stack[len(stack)-1])
	return stack
}

func OpHash160(stack Stack) Stack {
	if len(stack) < 1 {
		return stack
	}
	element := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	hash := utils.Hash160(element)
	stack = append(stack, hash)
	return stack
}

func OpHash256(stack Stack) Stack {
	if len(stack) < 1 {
		return stack
	}
	element := stack[len(stack)-1]
	stack = stack[:len(stack)-1]

	hash := utils.Hash256(element)
	stack = append(stack, hash)
	return stack
}

func ParseScript(reader io.Reader) (*Script, error) {
	script := NewScript()
	len, err := utils.ParseVarInt(reader)
	if err != nil {
		return nil, err
	}
	var i uint64
	for i < len {
		buf := make([]byte, 1)
		if _, err := io.ReadFull(reader, buf); err != nil {
			return nil, err
		}
		i += 1

		if buf[0] >= 1 && buf[0] <= 75 {
			// NOTE: element
			elementLen := int(buf[0])
			buf = make([]byte, elementLen)
			if _, err := io.ReadFull(reader, buf); err != nil {
				return nil, err
			}
			i += uint64(elementLen)

			script.Instructions = append(script.Instructions, buf)
		} else if buf[0] == OP_PUSHDATA1 {
			buf = make([]byte, 1)
			if _, err := io.ReadFull(reader, buf); err != nil {
				return nil, err
			}
			i += 1

			elementLen := int(buf[0])
			buf = make([]byte, elementLen)
			if _, err := io.ReadFull(reader, buf); err != nil {
				return nil, err
			}
			i += uint64(elementLen)

			script.Instructions = append(script.Instructions, buf)
		} else if buf[0] == OP_PUSHDATA2 {
			buf = make([]byte, 2)
			if _, err := io.ReadFull(reader, buf); err != nil {
				return nil, err
			}
			i += 2

			elementLen := binary.LittleEndian.Uint16(buf)
			buf = make([]byte, elementLen)
			if _, err := io.ReadFull(reader, buf); err != nil {
				return nil, err
			}
			i += uint64(elementLen)

			script.Instructions = append(script.Instructions, buf)
		} else {
			// NOTE: opcode
			script.Instructions = append(script.Instructions, buf)
		}
	}
	return script, nil
}

func (s *Script) Serialize() ([]byte, error) {
	buf := make([]byte, 0)
	for _, inst := range s.Instructions {
		if utils.IsInteger(inst) {
			// NOTE: opcode
			buf = append(buf, inst...)
		} else {
			// NOTE: element
			length := len(inst)
			// NOTE: serialize length first
			if length <= 75 {
				buf = append(buf, byte(length))
			} else if length > 75 && length < 0x100 {
				buf = append(buf, OP_PUSHDATA1)
				buf = append(buf, byte(length))
			} else if length >= 0x100 && length < 520 {
				buf = append(buf, OP_PUSHDATA2)
				binary.LittleEndian.AppendUint16(buf, uint16(length))
			} else {
				return nil, fmt.Errorf("element is too long")
			}
			buf = append(buf, inst...)
		}
	}
	return buf, nil
}

func (s *Script) Add(other *Script) {
	s.Instructions = append(s.Instructions, other.Instructions...)
}
