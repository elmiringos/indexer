package transaction

import "time"

type Transaction struct {
	Hash      string
	BlockHash string
	Index     int
	Status    int
	Gas       string
	GasUsed   string
	Input     string
	Amount    string
	From      string
	To        string
	Timestamp time.Time
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
		"amount":     t.Amount,
		"from":       t.From,
		"to":         t.To,
		"timestamp":  t.Timestamp,
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
	Index           int
	TransactionHash string
	FirstTopic      string
	SecondTopic     string
	ThirdTopic      string
	FourthTopic     string
	AddressHash     string
}

func (t *TransactionLog) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"index":            t.Index,
		"transaction_hash": t.TransactionHash,
		"first_topic":      t.FirstTopic,
		"second_topic":     t.SecondTopic,
		"third_topic":      t.ThirdTopic,
		"fourth_topic":     t.FourthTopic,
		"address_hash":     t.AddressHash,
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
