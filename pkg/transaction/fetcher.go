package transaction

import (
	"fmt"
	"net/http"
)

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
