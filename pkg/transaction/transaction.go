package transaction

import (
	"bufio"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"golang-bitcoin/pkg/script"
	"golang-bitcoin/pkg/utils"
	"io"
	"math/big"
)

type Transaction struct {
	Version  uint32
	Inputs   []*Input
	Outputs  []*Output
	Locktime uint32
	IsSegwit bool
}

type Input struct {
	PreviousOutputHash  []byte
	PreviousOutputIndex uint32
	ScriptSig           *script.Script
	Sequence            uint32
}

type Output struct {
	Value        uint64
	ScriptPubKey *script.Script
}

func NewTransaction(version uint32, inputs []*Input, outputs []*Output, locktime uint32, isSegwit bool) *Transaction {
	return &Transaction{version, inputs, outputs, locktime, isSegwit}
}

func ParseTransaction(reader io.Reader) (*Transaction, error) {
	breader := bufio.NewReader(reader)
	var buf []byte

	var version uint32
	buf = make([]byte, 4)
	if _, err := io.ReadFull(breader, buf); err != nil {
		return nil, fmt.Errorf("error reading version: %v", err)
	}
	version = binary.LittleEndian.Uint32(buf)

	peekBytes, err := breader.Peek(2)
	if err != nil {
		return nil, fmt.Errorf("error peeking for segwit flag: %v", err)
	}
	isSegwit := peekBytes[0] == 0x00 && peekBytes[1] == 0x01
	if isSegwit {
		if _, err := breader.Discard(2); err != nil {
			return nil, fmt.Errorf("error discarding segwit flag: %v", err)
		}
	}

	numInputs, err := utils.ParseVarInt(breader)
	if err != nil {
		return nil, fmt.Errorf("error reading number of inputs: %v", err)
	}
	inputs := make([]*Input, numInputs)
	for i := 0; i < int(numInputs); i++ {
		inputs[i], err = ParseInput(breader)
		if err != nil {
			return nil, fmt.Errorf("error reading input %d: %v", i, err)
		}
	}

	numOutputs, err := utils.ParseVarInt(breader)
	if err != nil {
		return nil, fmt.Errorf("error reading number of outputs: %v", err)
	}
	outputs := make([]*Output, numOutputs)
	for i := 0; i < int(numOutputs); i++ {
		outputs[i], err = ParseOutput(breader)
		if err != nil {
			return nil, fmt.Errorf("error reading output %d: %v", i, err)
		}
	}

	if isSegwit {
		// TODO: Parse witness data
	}

	remainingData, err := io.ReadAll(breader)
	if err != nil {
		return nil, fmt.Errorf("error reading remaining data: %v", err)
	}
	if len(remainingData) < 4 {
		return nil, fmt.Errorf("remaining data is too short to contain locktime")
	}
	locktime := binary.LittleEndian.Uint32(remainingData[len(remainingData)-4:])

	return &Transaction{version, inputs, outputs, locktime, isSegwit}, nil
}

func (t *Transaction) Serialize() ([]byte, error) {
	serialized := make([]byte, 0)

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, t.Version)
	serialized = append(serialized, buf...)

	numInputs, err := utils.SerializeVarInt(uint64(len(t.Inputs)))
	if err != nil {
		return nil, err
	}
	serialized = append(serialized, numInputs...)
	for _, input := range t.Inputs {
		serialized = append(serialized, input.Serialize()...)
	}

	numOutputs, err := utils.SerializeVarInt(uint64(len(t.Outputs)))
	if err != nil {
		return nil, err
	}
	serialized = append(serialized, numOutputs...)
	for _, output := range t.Outputs {
		serialized = append(serialized, output.Serialize()...)
	}

	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, t.Locktime)
	serialized = append(serialized, buf...)

	return serialized, nil
}

func (t *Transaction) ID() (string, error) {
	serialized, err := t.Serialize()
	if err != nil {
		return "", err
	}

	hash256 := utils.Hash256(serialized)

	return hex.EncodeToString(hash256), nil
}

func (t *Transaction) Fee(testnet bool) (uint64, error) {
	var inputSum uint64
	for _, input := range t.Inputs {
		value, err := input.Value(testnet)
		if err != nil {
			return 0, err
		}
		inputSum += value
	}
	var outputSum uint64
	for _, output := range t.Outputs {
		outputSum += output.Value
	}
	return inputSum - outputSum, nil
}

func (t *Transaction) DeepCopy() *Transaction {
	inputs := make([]*Input, len(t.Inputs))
	for i, input := range t.Inputs {
		inputs[i] = &Input{
			PreviousOutputHash:  input.PreviousOutputHash,
			PreviousOutputIndex: input.PreviousOutputIndex,
			ScriptSig:           input.ScriptSig,
			Sequence:            input.Sequence,
		}
	}

	outputs := make([]*Output, len(t.Outputs))
	for i, output := range t.Outputs {
		outputs[i] = &Output{
			Value:        output.Value,
			ScriptPubKey: output.ScriptPubKey,
		}
	}

	return &Transaction{
		Version:  t.Version,
		Inputs:   inputs,
		Outputs:  outputs,
		Locktime: t.Locktime,
	}
}

func (t *Transaction) SigHash(index int, testnet bool) ([]byte, error) {
	txCopy := t.DeepCopy()
	for i := 0; i < len(txCopy.Inputs); i++ {
		if i != index {
			txCopy.Inputs[i].ScriptSig = script.NewScript()
		}
	}
	scriptPubKey, err := txCopy.Inputs[index].ScriptPubKey(testnet)
	if err != nil {
		return nil, err
	}
	txCopy.Inputs[index].ScriptSig = scriptPubKey

	serialized, err := txCopy.Serialize()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, 1)
	serialized = append(serialized, buf...)
	hash := utils.Hash256(serialized)

	return hash, nil
}

func (t *Transaction) VerifyInput(index int, testnet bool) error {
	sigHash, err := t.SigHash(index, testnet)
	if err != nil {
		return err
	}
	z := new(big.Int).SetBytes(sigHash)
	scriptPubkey, err := t.Inputs[index].ScriptPubKey(testnet)
	if err != nil {
		return err
	}
	scriptSig := t.Inputs[index].ScriptSig
	scriptSig.Add(scriptPubkey)
	return scriptSig.Evaluate(z)
}

func (t *Transaction) Verify(testnet bool) error {
	fee, err := t.Fee(testnet)
	if err != nil {
		return err
	}
	if fee < 0 {
		return fmt.Errorf("transaction has negative fee")
	}
	for i := 0; i < len(t.Inputs); i++ {
		err := t.VerifyInput(i, testnet)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewInput(previousOutputHash []byte, previousOutputIndex uint32, scriptSig *script.Script, sequence uint32) *Input {
	return &Input{previousOutputHash, previousOutputIndex, scriptSig, sequence}
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

	scriptSig, err := script.ParseScript(reader)
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

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i.PreviousOutputIndex)
	serialized = append(serialized, buf...)

	serializedScriptSig, err := i.ScriptSig.Serialize()
	if err != nil {
		return nil
	}
	serialized = append(serialized, serializedScriptSig...)

	buf = make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, i.Sequence)
	serialized = append(serialized, buf...)

	return serialized
}

func (i *Input) Value(testnet bool) (uint64, error) {
	fetcher := NewTransactionFetcher(testnet)
	tx, err := fetcher.FetchTransaction(hex.EncodeToString(i.PreviousOutputHash), false)
	if err != nil {
		return 0, err
	}
	return tx.Outputs[i.PreviousOutputIndex].Value, nil
}

func (i *Input) ScriptPubKey(testnet bool) (*script.Script, error) {
	fetcher := NewTransactionFetcher(testnet)
	tx, err := fetcher.FetchTransaction(hex.EncodeToString(i.PreviousOutputHash), false)
	if err != nil {
		return nil, err
	}
	return tx.Outputs[i.PreviousOutputIndex].ScriptPubKey, nil
}

func NewOutput(value uint64, scriptPubKey *script.Script) *Output {
	return &Output{value, scriptPubKey}
}

func ParseOutput(reader io.Reader) (*Output, error) {
	var buf []byte

	buf = make([]byte, 8)
	if _, err := io.ReadFull(reader, buf); err != nil {
		return nil, err
	}
	value := binary.LittleEndian.Uint64(buf)

	scriptPubKey, err := script.ParseScript(reader)
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

	serializedScriptPubKey, err := o.ScriptPubKey.Serialize()
	if err != nil {
		return nil
	}
	serialized = append(serialized, serializedScriptPubKey...)

	return serialized
}
