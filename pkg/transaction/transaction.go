package transaction

import (
	"encoding/binary"
	"io"
)

type Transaction struct {
	Version  uint32
	Inputs   []*Input
	Outputs  []*Output
	Locktime uint32
}

type Input struct {
	PreviousOutputHash  []byte
	PreviousOutputIndex uint32
	ScriptSig           *ScriptSig
	Sequence            uint32
}

type ScriptSig struct{}

type Output struct {
	Value        uint64
	ScriptPubKey *ScriptPubKey
}

type ScriptPubKey struct{}

func NewTransaction(version uint32, inputs []*Input, outputs []*Output, locktime uint32) *Transaction {
	return &Transaction{version, inputs, outputs, locktime}
}

func ParseTransaction(reader io.Reader) (*Transaction, error) {
	var buf []byte

	var version uint32
	buf = make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	version = binary.LittleEndian.Uint32(buf)

	numInputs, err := ParseVarInt(reader)
	if err != nil {
		return nil, err
	}
	inputs := make([]*Input, numInputs)
	for i := 0; i < int(numInputs); i++ {
		inputs[i], err = ParseInput(reader)
	}

	numOutputs, err := ParseVarInt(reader)
	if err != nil {
		return nil, err
	}
	outputs := make([]*Output, numOutputs)
	for i := 0; i < int(numOutputs); i++ {
		outputs[i], err = ParseOutput(reader)
	}

	var locktime uint32
	buf = make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	locktime = binary.LittleEndian.Uint32(buf)

	return &Transaction{version, inputs, outputs, locktime}, nil
}

func (t *Transaction) Serialize() ([]byte, error) {
	var serialized []byte
	binary.LittleEndian.PutUint32(serialized, t.Version)

	numInputs, err := SerializeVarInt(uint64(len(t.Inputs)))
	if err != nil {
		return nil, err
	}
	serialized = append(serialized, numInputs...)
	for _, input := range t.Inputs {
		serialized = append(serialized, input.Serialize()...)
	}

	numOutputs, err := SerializeVarInt(uint64(len(t.Outputs)))
	if err != nil {
		return nil, err
	}
	serialized = append(serialized, numOutputs...)
	for _, output := range t.Outputs {
		serialized = append(serialized, output.Serialize()...)
	}
	binary.LittleEndian.PutUint32(serialized, t.Locktime)
	
	return serialized, nil
}

func ParseInput(reader io.Reader) (*Input, error) {
	var buf []byte

	previousOutputHash := make([]byte, 32)
	if _, err := io.ReadFull(reader, previousOutputHash); err != nil {
		return nil, err
	}

	buf = make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	previousOutputIndex := binary.LittleEndian.Uint32(buf)

	scriptSig, err := ParseScriptSig(reader)
	if err != nil {
		return nil, err
	}

	buf = make([]byte, 4)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	sequence := binary.LittleEndian.Uint32(buf)

	return &Input{previousOutputHash, previousOutputIndex, scriptSig, sequence}, nil
}

func (i *Input) Serialize() []byte {
	var serialized []byte
	serialized = append(serialized, i.PreviousOutputHash...)
	binary.LittleEndian.PutUint32(serialized, i.PreviousOutputIndex)
	serialized = append(serialized, i.ScriptSig.Serialize()...)
	binary.LittleEndian.PutUint32(serialized, i.Sequence)
	return serialized
}

func ParseScriptSig(reader io.Reader) (*ScriptSig, error) {
	return &ScriptSig{}, nil
}

func (s *ScriptSig) Serialize() []byte {
	return nil
}

func ParseOutput(reader io.Reader) (*Output, error) {
	var buf []byte

	buf = make([]byte, 8)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	value := binary.LittleEndian.Uint64(buf)

	scriptPubKey, err := ParseScriptPubKey(reader)
	if err != nil {
		return nil, err
	}

	return &Output{
		Value:        value,
		ScriptPubKey: scriptPubKey,
	}, nil
}

func (o *Output) Serialize() []byte {
	var serialized []byte
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, o.Value)
	serialized = append(serialized, buf...)
	serialized = append(serialized, o.ScriptPubKey.Serialize()...)
	return serialized
}

func ParseScriptPubKey(reader io.Reader) (*ScriptPubKey, error) {
	return &ScriptPubKey{}, nil
}

func (s *ScriptPubKey) Serialize() []byte {
	return nil
}
