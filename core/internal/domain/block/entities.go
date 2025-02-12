package block

import (
	"github.com/elmiringos/indexer/indexer-core/internal/domain"
	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	Hash          common.Hash    `json:"hash"`
	Number        domain.BigInt  `json:"number"`
	MinerHash     common.Address `json:"miner_hash"`
	ParentHash    common.Hash    `json:"parent_hash"`
	GasLimit      uint64         `json:"gas_limit"`
	GasUsed       uint64         `json:"gas_used"`
	Nonce         uint64         `json:"nonce"`
	Size          uint64         `json:"size"`
	Difficulty    domain.BigInt  `json:"difficulty"`
	IsPos         bool           `json:"is_pos"`
	BaseFeePerGas domain.BigInt  `json:"base_fee_per_gas"`
	Timestamp     uint64         `json:"timestamp"`
}

func (b *Block) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"hash":             b.Hash,
		"number":           b.Number,
		"miner_hash":       b.MinerHash,
		"parent_hash":      b.ParentHash,
		"gas_limit":        b.GasLimit,
		"gas_used":         b.GasUsed,
		"nonce":            b.Nonce,
		"size":             b.Size,
		"difficulty":       b.Difficulty,
		"is_pos":           b.IsPos,
		"base_fee_per_gas": b.BaseFeePerGas,
		"timestamp":        b.Timestamp,
	}
}

func MakeSlice(blocks []*Block) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(blocks))
	for i, block := range blocks {
		slices[i] = block.ToMap()
	}
	return slices
}
