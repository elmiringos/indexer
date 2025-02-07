package block

import "time"

type Block struct {
	Hash          string
	Number        int64
	MinerHash     string
	ParentHash    string
	GasLimit      int64
	GasUsed       int64
	Nonce         string
	Size          int64
	Difficulty    int64
	Consensus     bool
	BaseFeePerGas int64
	IsEmpty       bool
	Timestamp     time.Time
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
		"consensus":        b.Consensus,
		"base_fee_per_gas": b.BaseFeePerGas,
		"is_empty":         b.IsEmpty,
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
