package msgfinder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/glifio/go-pools/abigen"
)

const beryxURL = "https://api.zondax.ch/fil/data/v3/mainnet"

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Height    uint64   `json:"height"`
	TipsetCID string   `json:"tipset_cid"`
	BlockCID  string   `json:"block_cid"`
	From      string   `json:"tx_from"`
	To        string   `json:"tx_to"`
	Amount    *big.Int `json:"amount"`
	CID       string   `json:"tx_cid"`
	Status    string   `json:"status"`
	Type      string   `json:"tx_type"`
	SearchID  string   `json:"search_id"`
}

type MethodLookupError struct {
	Err error
}

func (r *MethodLookupError) Error() string {
	return fmt.Sprintf("%v", r.Err)
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
	Params  interface{}
	Return  interface{}
	EthLogs []interface{} `json:"ethLogs"`
}

type TransactionDetail struct {
	Level      int        `json:"level"`
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
	/*
		fmt.Println("Detail:")
		fmt.Println(string(resBody))
	*/

	var detail TransactionDetail

	err = json.Unmarshal(resBody, &detail)
	if err != nil {
		return TransactionDetail{}, err
	}
	// fmt.Printf("Jim: %+v", detail)

	return detail, nil
}

func (td *TransactionDetail) ParseParams() (*abi.Method, map[string]interface{}, error) {
	if td.TxMetaData.Params == nil || td.TxMetaData.Params == "" {
		return nil, nil, nil
	}
	p, ok := td.TxMetaData.Params.(string)
	if !ok {
		return nil, nil, nil
	}

	// fmt.Printf("Params: %+v\n", td.TxMetaData.Params)
	data := common.FromHex(p)
	// fmt.Printf("Params bytes: %+v\n", data)

	if len(data) == 0 {
		return nil, nil, nil
	}

	var paramsBytes []byte = data

	sig := paramsBytes[0:4]

	abi, err := abigen.AgentMetaData.GetAbi()
	if err != nil {
		return nil, nil, err
	}

	// fmt.Printf("Sig: %+v\n", hex.EncodeToString(sig))

	method, err := abi.MethodById(sig)
	if err != nil {
		// Try unpacking the params as CBOR
		reader := bytes.NewReader(data)
		err := cborutil.ReadCborRPC(reader, &paramsBytes)
		if err != nil {
			return nil, nil, &MethodLookupError{Err: err}
		}
		sig = paramsBytes[0:4]
		method, err = abi.MethodById(sig)
		if err != nil {
			return nil, nil, &MethodLookupError{Err: err}
		}
	}

	unpackedMap := make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(unpackedMap, paramsBytes[4:])
	if err != nil {
		return nil, nil, err
	}
	// fmt.Printf("Unpacked Map: %+v\n", unpackedMap)

	return method, unpackedMap, nil
}
