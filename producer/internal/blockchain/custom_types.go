// Package blockchain contains custom types for the blockchain cause of the json marshalling issues
package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Block represents a block in the blockchain
type Block struct {
	Hash          common.Hash    `json:"hash"`
	Number        BigInt         `json:"number"`
	MinerHash     common.Address `json:"miner_hash"`
	ParentHash    common.Hash    `json:"parent_hash"`
	GasLimit      uint64         `json:"gas_limit"`
	GasUsed       uint64         `json:"gas_used"`
	Nonce         uint64         `json:"nonce"`
	Size          uint64         `json:"size"`
	Difficulty    BigInt         `json:"difficulty"`
	IsPos         bool           `json:"is_pos"`
	BaseFeePerGas BigInt         `json:"base_fee_per_gas"`
	Timestamp     uint64         `json:"timestamp"`
}

// ConvertBlockToBlock converts a types.Block to a custom type Block
func ConvertBlockToBlock(block *types.Block) *Block {
	isPoS := block.Difficulty().Cmp(big.NewInt(0)) == 0 && block.Nonce() == 0 && len(block.Extra()) >= 32

	var baseFee *big.Int
	if block.Header().BaseFee != nil {
		baseFee = block.Header().BaseFee
	} else {
		baseFee = big.NewInt(0)
	}

	return &Block{
		Hash:          block.Hash(),
		Number:        BigInt(*block.Number()),
		MinerHash:     block.Coinbase(),
		ParentHash:    block.ParentHash(),
		GasLimit:      block.GasLimit(),
		GasUsed:       block.GasUsed(),
		Nonce:         block.Nonce(),
		Size:          block.Size(),
		Difficulty:    BigInt(*block.Difficulty()),
		IsPos:         isPoS,
		BaseFeePerGas: BigInt(*baseFee),
		Timestamp:     block.Header().Time,
	}
}

// Transaction represents a transaction in the blockchain
type Transaction struct {
	Hash      common.Hash    `json:"hash"`
	BlockHash common.Hash    `json:"block_hash"`
	Index     int            `json:"index"`
	Status    uint64         `json:"status"`
	Gas       uint64         `json:"gas"`
	GasUsed   uint64         `json:"gas_used"`
	Input     []byte         `json:"input"`
	Value     BigInt         `json:"value"`
	From      common.Address `json:"from"`
	To        common.Address `json:"to"`
	Nonce     uint64         `json:"nonce"`
	Timestamp int64          `json:"timestamp"`
}

// ConvertTransactionToTransaction converts a types.Transaction to a custom type Transaction
func (p *BlockchainProcessor) ConvertTransactionToTransaction(transaction *types.Transaction, blockHash common.Hash, receipt *types.Receipt, index int) (*Transaction, error) {
	transactionSender, err := p.GetTransactionSender(transaction)
	if err != nil {
		return nil, err
	}

	transactionMessage := &Transaction{
		Hash:      transaction.Hash(),
		BlockHash: blockHash,
		Index:     index,
		Status:    receipt.Status,
		Gas:       transaction.Gas(),
		GasUsed:   receipt.GasUsed,
		Input:     transaction.Data(),
		Value:     BigInt(*transaction.Value()),
		From:      transactionSender,
		Nonce:     transaction.Nonce(),
		Timestamp: transaction.Time().Unix(),
	}

	toAddress := transaction.To()
	if toAddress == nil {
		transactionMessage.To = common.Address{}
	} else {
		transactionMessage.To = *toAddress
	}

	return transactionMessage, nil
}

// Reward represents a reward for a validator
type Reward struct {
	BlockHash common.Hash    `json:"block_hash"`
	Address   common.Address `json:"address"`
	Amount    uint64         `json:"amount"`
}

// Withdrawal represents a withdrawal from a validator
// Using custom type for tracking block hash
type Withdrawal struct {
	Index          uint64         `json:"index"`
	BlockHash      common.Hash    `json:"block_hash"`
	AddressHash    common.Address `json:"address_hash"`
	ValidatorIndex uint64         `json:"validator_index"`
	Amount         uint64         `json:"amount"`
}

// ConvertWithdrawalToWithdrawal converts a types.Withdrawal to a custom type Withdrawal
func ConvertWithdrawalToWithdrawal(withdrawal *types.Withdrawal, blockHash common.Hash) *Withdrawal {
	return &Withdrawal{
		Index:          withdrawal.Index,
		BlockHash:      blockHash,
		AddressHash:    withdrawal.Address,
		ValidatorIndex: withdrawal.Validator,
		Amount:         withdrawal.Amount,
	}
}

// TokenMetadata represents the metadata of a token
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
