package blockchain

import (
	"math/big"
)

type TokenMetadata map[string]interface{}

// TokenEvent represents a token event
type TokenEvent struct {
	From          string         `json:"from"`
	To            string         `json:"to"`
	Value         *big.Int       `json:"value"`
	TokenId       *big.Int       `json:"token_id"`
	TokenMetadata *TokenMetadata `json:"token_metadata"`
	IsMint        bool           `json:"is_mint"`
	IsBurn        bool           `json:"is_burn"`
}

// Reward represents a reward for a validator
type Reward struct {
	Address   string   `json:"address"`
	Amount    uint64   `json:"amount"`
	Block     *big.Int `json:"block"`
	BlockHash string   `json:"block_hash"`
}

// InternalTransaction represents an internal transaction
type InternalTransaction struct {
	BlockHash       string   `json:"block_hash"`
	Index           int      `json:"index"`
	Type            string   `json:"type"`
	TransactionHash string   `json:"transaction_hash"`
	Status          int      `json:"status"`
	Gas             uint64   `json:"gas"`
	GasUsed         uint64   `json:"gas_used"`
	Input           []byte   `json:"input"`
	Output          []byte   `json:"output"`
	Value           *big.Int `json:"value"`
	From            string   `json:"from"`
	To              string   `json:"to"`
	ContractAddress string   `json:"contract_address"`
	Timestamp       uint64   `json:"timestamp"`
	ErrorMsg        string   `json:"error_msg"`
}

type TransactionAction struct {
	TransactionHash string   `json:"transaction_hash"`
	Selector        string   `json:"selector"`
	Type            string   `json:"type"`
	From            string   `json:"from"`
	To              string   `json:"to"`
	Value           *big.Int `json:"value"`
	Input           []byte   `json:"input"`
	Status          int      `json:"status"`
}

const zeroAddress = "0x0000000000000000000000000000000000000000"
