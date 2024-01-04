package msgfinder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
)

const beryxURL = "https://api.zondax.ch/fil/data/v3/mainnet"

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Height    uint64 `json:"height"`
	TipsetCID string `json:"tipset_cid"`
	BlockCID  string `json:"block_cid"`
	From      string `json:"tx_from"`
	To        string `json:"tx_to"`
	CID       string `json:"tx_cid"`
	Status    string `json:"status"`
	Type      string `json:"tx_type"`
	SearchID  string `json:"search_id"`
}

func GetTransactions(ctx context.Context, agent common.Address) ([]Transaction, error) {
	url := fmt.Sprintf("%s/transactions/address/%v/receiver", beryxURL, agent)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []Transaction{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", os.Getenv("BERYX_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Transaction{}, err
	}

	if res.StatusCode != http.StatusOK {
		return []Transaction{}, fmt.Errorf("bad http status: %v", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return []Transaction{}, err
	}
	// fmt.Println(string(resBody))

	var txs Transactions

	err = json.Unmarshal(resBody, &txs)
	if err != nil {
		return []Transaction{}, err
	}

	return txs.Transactions, nil
}

type TxMetaData struct {
	Params  string
	Return  string
	EthLogs []interface{} `json:"ethLogs"`
}

type TransactionDetail struct {
	TxMetaData TxMetaData `json:"tx_metadata"`
}

func GetTransactionDetail(ctx context.Context, searchID string) (TransactionDetail, error) {
	url := fmt.Sprintf("%s/transactions/id/%s", beryxURL, searchID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return TransactionDetail{}, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", os.Getenv("BERYX_TOKEN")))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TransactionDetail{}, err
	}

	if res.StatusCode != http.StatusOK {
		return TransactionDetail{}, fmt.Errorf("bad http status: %v", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return TransactionDetail{}, err
	}
	// fmt.Println(string(resBody))

	var detail TransactionDetail

	err = json.Unmarshal(resBody, &detail)
	if err != nil {
		return TransactionDetail{}, err
	}

	return detail, nil
}