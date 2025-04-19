package transaction

import (
	"github.com/elmiringos/indexer/explorer/internal/domain"
	"github.com/ethereum/go-ethereum/common"
)

type Transaction struct {
	Hash      common.Hash    `json:"hash"`
	BlockHash common.Hash    `json:"block_hash"`
	Index     int            `json:"index"`
	Status    uint64         `json:"status"`
	Gas       uint64         `json:"gas"`
	GasUsed   uint64         `json:"gas_used"`
	Input     []byte         `json:"input"`
	Value     domain.BigInt  `json:"value"`
	From      common.Address `json:"from"`
	To        common.Address `json:"to"`
	Nonce     uint64         `json:"nonce"`
	Timestamp int64          `json:"timestamp"`
	LogsCount int            `json:"logs_count"`
}

func (t *Transaction) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"hash":       t.Hash,
		"block_hash": t.BlockHash,
		"index":      t.Index,
		"status":     t.Status,
		"gas":        t.Gas,
		"gas_used":   t.GasUsed,
		"input":      t.Input,
		"value":      t.Value,
		"from":       t.From,
		"to":         t.To,
		"timestamp":  t.Timestamp,
		"nonce":      t.Nonce,
	}
}

func MakeSlice(transactions []*Transaction) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(transactions))
	for i, transaction := range transactions {
		slices[i] = transaction.ToMap()
	}
	return slices
}

type TransactionLog struct {
	Address          common.Address `json:"address"`
	Topics           []common.Hash  `json:"topics"`
	TransactionHash  common.Hash    `json:"transactionHash"`
	BlockHash        common.Hash    `json:"blockHash"`
	TransactionIndex uint           `json:"transactionIndex"`
	Index            uint           `json:"logIndex"`
	Data             []byte         `json:"data"`
}

func (t *TransactionLog) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"address":           t.Address,
		"topics":            t.Topics,
		"transaction_hash":  t.TransactionHash,
		"block_has":         t.BlockHash,
		"transaction_index": t.TransactionIndex,
		"index":             t.Index,
		"data":              t.Data,
	}
}

func MakeTransactionLogSlice(transactionLogs []*TransactionLog) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(transactionLogs))
	for i, transactionLog := range transactionLogs {
		slices[i] = transactionLog.ToMap()
	}
	return slices
}

type TransactionAction struct {
	TransactionHash string
	LogIndex        int
	Data            string
	AddressHash     string
	Type            int
}

func (t *TransactionAction) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"transaction_hash": t.TransactionHash,
		"log_index":        t.LogIndex,
		"data":             t.Data,
		"address_hash":     t.AddressHash,
		"type":             t.Type,
	}
}

func MakeTransactionActionSlice(transactionActions []*TransactionAction) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(transactionActions))
	for i, transactionAction := range transactionActions {
		slices[i] = transactionAction.ToMap()
	}
	return slices
}
