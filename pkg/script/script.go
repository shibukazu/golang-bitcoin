package script

import (
	"encoding/binary"
	"fmt"
	"golang-bitcoin/pkg/secp256k1"
	"golang-bitcoin/pkg/utils"
	"io"
	"math/big"
)

type Script struct {
	Instructions [][]byte
	Stack        [][]byte
	AltStack     [][]byte
}

func NewScript() *Script {
	return &Script{
		Instructions: make([][]byte, 0),
	}
}

func (s *Script) PopInstruction() ([]byte, error) {
	if len(s.Instructions) == 0 {
		return nil, fmt.Errorf("no instructions to pop")
	}
	element := s.Instructions[0]
	s.Instructions = s.Instructions[1:]
	return element, nil
}

func (s *Script) PopStack() ([]byte, error) {
	if len(s.Stack) == 0 {
		return nil, fmt.Errorf("no stack to pop")
	}
	element := s.Stack[len(s.Stack)-1]
	s.Stack = s.Stack[:len(s.Stack)-1]
	return element, nil
}

func (s *Script) PopAltStack() ([]byte, error) {
	if len(s.AltStack) == 0 {
		return nil, fmt.Errorf("no alt stack to pop")
	}
	element := s.AltStack[len(s.AltStack)-1]
	s.AltStack = s.AltStack[:len(s.AltStack)-1]
	return element, nil
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
		if IsOp(inst) {
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

func (s *Script) Evaluate(z *big.Int) error {
	for len(s.Instructions) > 0 {
		inst, err := s.PopInstruction()
		if err != nil {
			return err
		}
		if IsOp(inst) {
			/// NOTE: opcode
			intInst := decodeNum(inst)
			switch intInst {
			case OP_DUP:
				s.OpDup()
			case OP_HASH160:
				s.OpHash160()
			case OP_HASH256:
				s.OpHash256()
			case OP_EQUAL:
				s.OpEqual()
			case OP_EQUALVERIFY:
				s.OpEqualVerify()
			case OP_IF:
				s.OpIf()
			case OP_NOTIF:
				s.OpNotIf()
			case OP_TOALTSTACK:
				s.OpToAltStack()
			case OP_FROMALTSTACK:
				s.OpFromAltStack()
			case OP_CHECKSIG:
				s.OpCheckSig(z)
			default:
				return fmt.Errorf("unsupported opcode")
			}
		} else {
			// NOTE: element
			s.Stack = append(s.Stack, inst)
		}
	}
	if len(s.Stack) == 0 {
		return fmt.Errorf("stack is empty")
	}
	if len(s.Stack) > 1 {
		return fmt.Errorf("stack has multiple elements")
	}
	if len(s.AltStack) > 0 {
		return fmt.Errorf("alt stack is not empty")
	}
	element := s.Stack[0]
	if decodeNum(element) == 0 {
		return fmt.Errorf("stack top element is zero")
	}

	return nil
}

func NewP2PKHScriptPubkey(address string) (*Script, error) {
	hash160, err := secp256k1.ExtractHash160(address)
	if err != nil {
		return nil, err
	}
	script := NewScript()
	script.Instructions = append(script.Instructions, []byte{OP_DUP})
	script.Instructions = append(script.Instructions, []byte{OP_HASH160})
	script.Instructions = append(script.Instructions, hash160)
	script.Instructions = append(script.Instructions, []byte{OP_EQUALVERIFY})
	script.Instructions = append(script.Instructions, []byte{OP_CHECKSIG})
	return script, nil
}

func NewScriptSig(serializedSignature, serializedPubkey []byte) *Script {
	script := NewScript()
	script.Instructions = append(script.Instructions, serializedSignature)
	script.Instructions = append(script.Instructions, serializedPubkey)
	return script
}
