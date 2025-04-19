package token

import (
	"encoding/json"

	"github.com/elmiringos/indexer/explorer/internal/domain"
	"github.com/ethereum/go-ethereum/common"
)

// TokenMetadata represents the metadata of a token
type TokenMetadata map[string]interface{}

// TokenEvent represents a token event
type TokenEvent struct {
	Address               common.Address `json:"address"`
	TransactionHash       common.Hash    `json:"transaction_hash"`
	LogIndex              int            `json:"log_index"`
	From                  common.Address `json:"from"`
	To                    common.Address `json:"to"`
	Value                 domain.BigInt  `json:"value"`
	TokenId               domain.BigInt  `json:"token_id"`
	TokenMetadata         *TokenMetadata `json:"token_metadata"`
	IsMint                bool           `json:"is_mint"`
	IsBurn                bool           `json:"is_burn"`
	SmartContractDeployed bool           `json:"smart_contract_deployed"`
}

type Token struct {
	Address              common.Address
	Name                 string
	Symbol               string
	TotalSupply          domain.BigInt
	Decimals             int
	FiatValue            domain.BigInt
	CirculationMarketCap domain.BigInt
}

func (t *Token) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"address":                t.Address,
		"name":                   t.Name,
		"symbol":                 t.Symbol,
		"total_supply":           t.TotalSupply,
		"decimals":               t.Decimals,
		"fiat_value":             t.FiatValue,
		"circulation_market_cap": t.CirculationMarketCap,
	}
}

func MakeTokenSlice(tokens []*Token) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(tokens))
	for i, token := range tokens {
		slices[i] = token.ToMap()
	}

	return slices
}

type TokenTransfer struct {
	TransactionHash      common.Hash
	LogIndex             int
	From                 common.Address
	To                   common.Address
	TokenContractAddress common.Address
	Amount               domain.BigInt
}

func (t *TokenTransfer) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"transaction_hash":       t.TransactionHash,
		"log_index":              t.LogIndex,
		"from":                   t.From,
		"to":                     t.To,
		"token_contract_address": t.TokenContractAddress,
		"amount":                 t.Amount,
	}
}

func MakeTokenTransferSlice(tokenTransfers []*TokenTransfer) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(tokenTransfers))
	for i, tokenTransfer := range tokenTransfers {
		slices[i] = tokenTransfer.ToMap()
	}

	return slices
}

type TokenInstance struct {
	TokenId              domain.BigInt
	TokenContractAddress common.Address
	OwnerAddress         common.Address
	Metadata             *json.RawMessage
}

func (t *TokenInstance) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"token_id":               t.TokenId,
		"token_contract_address": t.TokenContractAddress,
		"owner_address":          t.OwnerAddress,
		"metadata":               t.Metadata,
	}
}

func MakeTokenInstanceSlice(tokenInstances []*TokenInstance) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(tokenInstances))
	for i, tokenInstance := range tokenInstances {
		slices[i] = tokenInstance.ToMap()
	}

	return slices
}
