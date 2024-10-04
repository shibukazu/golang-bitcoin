package transaction

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
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

func (t *Transaction) ID() (string, error) {
	serialized, err := t.Serialize()
	if err != nil {
		return "", err
	}
	firstHash := sha256.Sum256(serialized)
	secondHash := sha256.Sum256(firstHash[:])

	return hex.EncodeToString(secondHash[:]), nil
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

func (i *Input) Value(testnet bool) (uint64, error) {
	fetcher := NewTransactionFetcher(testnet)
	tx, err := fetcher.FetchTransaction(hex.EncodeToString(i.PreviousOutputHash), false)
	if err != nil {
		return 0, err
	}
	return tx.Outputs[i.PreviousOutputIndex].Value, nil
}

func (i *Input) ScriptPubKey(testnet bool) (*ScriptPubKey, error) {
	fetcher := NewTransactionFetcher(testnet)
	tx, err := fetcher.FetchTransaction(hex.EncodeToString(i.PreviousOutputHash), false)
	if err != nil {
		return nil, err
	}
	return tx.Outputs[i.PreviousOutputIndex].ScriptPubKey, nil
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

type TransactionFetcher struct {
	url    string
	cached map[string]*Transaction
}

func NewTransactionFetcher(testnet bool) *TransactionFetcher {
	var url string
	if testnet {
		url = "http://testnet.programmingbitcoin.com"
	} else {
		url = "http://mainnet.programmingbitcoin.com"
	}
	return &TransactionFetcher{url, make(map[string]*Transaction)}
}

func (tf *TransactionFetcher) FetchTransaction(txid string, fresh bool) (*Transaction, error) {
	if !fresh && tf.cached[txid] != nil {
		return tf.cached[txid], nil
	}
	url := fmt.Sprintf("%s/tx/%s.hex", tf.url, txid)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("error fetching transaction: %s", resp.Status)
	}

	tx, err := ParseTransaction(resp.Body)
	if err != nil {
		return nil, err
	}

	actualTxid, err := tx.ID()
	if err != nil {
		return nil, err
	}
	if actualTxid != txid {
		return nil, fmt.Errorf("fetched transaction id does not match expected: %s != %s", actualTxid, txid)
	}

	tf.cached[txid] = tx

	return tx, nil
}
