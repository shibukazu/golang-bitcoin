package script

import (
	"fmt"
	"golang-bitcoin/pkg/secp256k1"
	"golang-bitcoin/pkg/signature"
	"golang-bitcoin/pkg/utils"
	"math/big"
)

const (
	OP_DUP                 = 0x76
	OP_HASH160             = 0xa9
	OP_HASH256             = 0xaa
	OP_PUSHDATA1           = 0x4c
	OP_PUSHDATA2           = 0x4d
	OP_IF                  = 0x63
	OP_NOTIF               = 0x64
	OP_ELSE                = 0x67
	OP_ENDIF               = 0x68
	OP_TOALTSTACK          = 0x6b
	OP_FROMALTSTACK        = 0x6c
	OP_CHECKSIG            = 0xac
	OP_CHECKSIGVERIFY      = 0xad
	OP_CHECKMULTISIG       = 0xae
	OP_CHECKMULTISIGVERIFY = 0xaf
)

// NOTE: Op呼び出時のInstructionsはOp自体を含まない

func (s *Script) OpNumber(num int64) error {
	s.Stack = append(s.Stack, encodeNum(num))
	return nil
}

func (s *Script) OpDup() error {
	if len(s.Stack) < 1 {
		return fmt.Errorf("stack is empty")
	}
	s.Stack = append(s.Stack, s.Stack[len(s.Stack)-1])
	return nil
}

func (s *Script) OpHash160() error {
	if len(s.Stack) < 1 {
		return fmt.Errorf("stack is empty")
	}
	element, err := s.PopStack()
	if err != nil {
		return err
	}

	hash := utils.Hash160(element)
	s.Stack = append(s.Stack, hash)
	return nil
}

func (s *Script) OpHash256() error {
	if len(s.Stack) < 1 {
		return fmt.Errorf("stack is empty")
	}
	element, err := s.PopStack()
	if err != nil {
		return err
	}

	hash := utils.Hash256(element)
	s.Stack = append(s.Stack, hash)
	return nil
}

func (s *Script) OpIf() error {
	if len(s.Stack) < 1 {
		return fmt.Errorf("stack is empty")
	}
	true_instructions := make([][]byte, 0)
	false_instructions := make([][]byte, 0)
	new_instructions := &true_instructions
	numEndIf := 1
	valid := false

	for len(s.Instructions) > 0 {
		instruction, err := s.PopInstruction()
		if err != nil {
			return err
		}
		switch instruction[0] {
		case OP_IF:
			numEndIf += 1
			*new_instructions = append(*new_instructions, instruction)
		case OP_NOTIF:
			numEndIf += 1
			*new_instructions = append(*new_instructions, instruction)
		case OP_ELSE:
			if numEndIf == 1 {
				// NOTE: このときのみ負の分岐に入る
				new_instructions = &false_instructions
			} else {
				*new_instructions = append(*new_instructions, instruction)
			}
		case OP_ENDIF:
			numEndIf -= 1
			if numEndIf == 0 {
				valid = true
				break
			} else {
				*new_instructions = append(*new_instructions, instruction)
			}
		default:
			*new_instructions = append(*new_instructions, instruction)
		}
	}

	if !valid {
		return fmt.Errorf("invalid if-else-endif block")
	}

	// NOTE: Stackの先頭要素に基づき実行する分岐を決定
	element, err := s.PopStack()
	if err != nil {
		return err
	}
	decodedElement := decodeNum(element)
	if decodedElement != 0 {
		s.Instructions = append(s.Instructions, true_instructions...)
	} else {
		s.Instructions = append(s.Instructions, false_instructions...)
	}

	return nil
}

func (s *Script) OpNotIf() error {
	if len(s.Stack) < 1 {
		return fmt.Errorf("stack is empty")
	}
	true_instructions := make([][]byte, 0)
	false_instructions := make([][]byte, 0)
	new_instructions := &true_instructions
	numEndIf := 1
	valid := false

	for len(s.Instructions) > 0 {
		instruction, err := s.PopInstruction()
		if err != nil {
			return err
		}
		switch instruction[0] {
		case OP_IF:
			numEndIf += 1
			*new_instructions = append(*new_instructions, instruction)
		case OP_NOTIF:
			numEndIf += 1
			*new_instructions = append(*new_instructions, instruction)
		case OP_ELSE:
			if numEndIf == 1 {
				// NOTE: このときのみ負の分岐に入る
				new_instructions = &false_instructions
			} else {
				*new_instructions = append(*new_instructions, instruction)
			}
		case OP_ENDIF:
			numEndIf -= 1
			if numEndIf == 0 {
				valid = true
				break
			} else {
				*new_instructions = append(*new_instructions, instruction)
			}
		default:
			*new_instructions = append(*new_instructions, instruction)
		}
	}

	if !valid {
		return fmt.Errorf("invalid if-else-endif block")
	}

	// NOTE: Stackの先頭要素に基づき実行する分岐を決定
	element, err := s.PopStack()
	if err != nil {
		return err
	}
	decodedElement := decodeNum(element)
	if decodedElement == 0 {
		s.Instructions = append(s.Instructions, true_instructions...)
	} else {
		s.Instructions = append(s.Instructions, false_instructions...)
	}

	return nil
}

func (s *Script) OpToAltStack() error {
	if len(s.Stack) < 1 {
		return fmt.Errorf("stack is empty")
	}
	element, err := s.PopStack()
	if err != nil {
		return err
	}
	s.AltStack = append(s.AltStack, element)
	return nil
}

func (s *Script) OpFromAltStack() error {
	if len(s.AltStack) < 1 {
		return fmt.Errorf("alt stack is empty")
	}
	element, err := s.PopAltStack()
	if err != nil {
		return err
	}
	s.Stack = append(s.Stack, element)
	return nil
}

func (s *Script) OpCheckSig(z *big.Int) error {
	if len(s.Stack) < 2 {
		return fmt.Errorf("stack is empty")
	}
	secPubkey, err := s.PopStack()
	if err != nil {
		return err
	}
	derSig, err := s.PopStack()
	if err != nil {
		return err
	}

	pubkey := secp256k1.ParseSecp256k1Point(secPubkey)
	sig, err := signature.ParseSignature(derSig)
	if err != nil {
		return err
	}
	valid := pubkey.Verify(z, *sig)
	if valid {
		s.Stack = append(s.Stack, encodeNum(1))
	} else {
		s.Stack = append(s.Stack, encodeNum(0))
	}

	return nil
}
