package reward

import (
	"github.com/ethereum/go-ethereum/common"
)

type Reward struct {
	BlockHash common.Hash    `json:"block_hash"`
	Address   common.Address `json:"address"`
	Amount    uint64         `json:"amount"`
}

func (r *Reward) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"block_hash": r.BlockHash,
		"address":    r.Address,
		"amount":     r.Amount,
	}
}

func MakeRewardSlice(rewards []*Reward) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(rewards))
	for i, reward := range rewards {
		slices[i] = reward.ToMap()
	}
	return slices
}
