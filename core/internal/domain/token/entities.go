package token

import (
	"encoding/json"
	"math/big"
)

type Token struct {
	AddressHash          string
	Name                 string
	Symbol               string
	TotalSupply          *big.Int
	Decimals             int
	HolderCount          int
	FiatValue            *big.Int
	CirculationMarketCap *big.Int
}

func (t *Token) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"address_hash":           t.AddressHash,
		"name":                   t.Name,
		"symbol":                 t.Symbol,
		"total_supply":           t.TotalSupply,
		"decimals":               t.Decimals,
		"holder_count":           t.HolderCount,
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
	TransactionHash          string
	LogIndex                 int
	From                     string
	To                       string
	TokenContractAddressHash string
	Amount                   *big.Int
}

func (t *TokenTransfer) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"transaction_hash":            t.TransactionHash,
		"log_index":                   t.LogIndex,
		"from":                        t.From,
		"to":                          t.To,
		"token_contract_address_hash": t.TokenContractAddressHash,
		"amount":                      t.Amount,
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
	TokenId                  *big.Int
	TokenContractAddressHash string
	OwnerAddressHash         string
	Metadata                 *json.RawMessage
}

func (t *TokenInstance) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"token_id":                    t.TokenId,
		"token_contract_address_hash": t.TokenContractAddressHash,
		"owner_address_hash":          t.OwnerAddressHash,
		"metadata":                    t.Metadata,
	}
}

func MakeTokenInstanceSlice(tokenInstances []*TokenInstance) []map[string]interface{} {
	slices := make([]map[string]interface{}, len(tokenInstances))
	for i, tokenInstance := range tokenInstances {
		slices[i] = tokenInstance.ToMap()
	}
	return slices
}
