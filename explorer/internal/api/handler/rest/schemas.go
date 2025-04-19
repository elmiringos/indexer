package rest

import (
	"github.com/elmiringos/indexer/explorer/internal/domain/block"
)

type BlockResponse struct {
	Hash          string `json:"hash"`
	Number        string `json:"number"`
	MinerHash     string `json:"miner_hash"`
	ParentHash    string `json:"parent_hash"`
	GasLimit      uint64 `json:"gas_limit"`
	GasUsed       uint64 `json:"gas_used"`
	Nonce         uint64 `json:"nonce"`
	Size          uint64 `json:"size"`
	Difficulty    string `json:"difficulty"`
	IsPos         bool   `json:"is_pos"`
	BaseFeePerGas string `json:"base_fee_per_gas"`
	Timestamp     uint64 `json:"timestamp"`
}

func MapBlockToCurrentBlockResponse(block *block.Block) *BlockResponse {
	return &BlockResponse{
		Hash:          block.Hash.String(),
		Number:        block.Number.String(),
		MinerHash:     block.MinerHash.String(),
		ParentHash:    block.ParentHash.String(),
		GasLimit:      block.GasLimit,
		GasUsed:       block.GasUsed,
		Nonce:         block.Nonce,
		Size:          block.Size,
		Difficulty:    block.Difficulty.String(),
		IsPos:         block.IsPos,
		BaseFeePerGas: block.BaseFeePerGas.String(),
		Timestamp:     block.Timestamp,
	}
}
