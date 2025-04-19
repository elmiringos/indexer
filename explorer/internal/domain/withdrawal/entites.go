package withdrawal

import (
	"github.com/ethereum/go-ethereum/common"
)

type Withdrawal struct {
	Index          uint64         `json:"index"`
	BlockHash      common.Hash    `json:"block_hash"`
	AddressHash    common.Address `json:"address_hash"`
	ValidatorIndex uint64         `json:"validator_index"`
	Amount         uint64         `json:"amount"`
}

func (w *Withdrawal) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"index":           w.Index,
		"block_hash":      w.BlockHash,
		"address_hash":    w.AddressHash,
		"validator_index": w.ValidatorIndex,
		"amount":          w.Amount,
	}
}

func MakeWithdrawalSlice(withdrawals []*Withdrawal) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(withdrawals))
	for i, withdrawal := range withdrawals {
		slices[i] = withdrawal.ToMap()
	}
	return slices
}
