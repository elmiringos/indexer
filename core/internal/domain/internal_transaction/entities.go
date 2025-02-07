package internal_transaction

import "math/big"

type InternalTransaction struct {
	BlockHash       string
	Index           int
	TransactionHash string
	Status          int
	Gas             uint64
	GasUsed         uint64
	Input           []byte
	Output          []byte
	Value           *big.Int
	From            string
	To              string
	ContractAddress string
	Timestamp       uint64
	ErrorMsg        string
}

func (i *InternalTransaction) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"block_hash": i.BlockHash,
	}
}

func MakeSlice(internalTransactions []*InternalTransaction) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(internalTransactions))
	for i, internalTransaction := range internalTransactions {
		slices[i] = internalTransaction.ToMap()
	}
	return slices
}
