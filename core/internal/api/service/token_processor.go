package service

import (
	"context"
	"encoding/json"
	"fmt"

	smartcontract "github.com/elmiringos/indexer/indexer-core/internal/domain/smart_contract"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/token"
	"go.uber.org/zap"
)

type TokenProcessor struct {
	tokenRepository         token.Repository
	smartContractRepository smartcontract.Repository
	log                     *zap.Logger
}

func NewTokenProccesor(
	tokenRepository token.Repository,
	smartContractRepository smartcontract.Repository,
	log *zap.Logger,
) *TokenProcessor {
	return &TokenProcessor{
		tokenRepository:         tokenRepository,
		smartContractRepository: smartContractRepository,
		log:                     log,
	}
}

func (p *TokenProcessor) Process(ctx context.Context, data []byte) error {
	// Unmarshal the token event data into the TokenEvent struct
	tokenEvent := &token.TokenEvent{}
	if err := json.Unmarshal(data, tokenEvent); err != nil {
		return err
	}

	metadata := *tokenEvent.TokenMetadata

	if tokenEvent.IsMint {
		// Process Token Entity
		tokenEntity := &token.Token{
			Address:     tokenEvent.From, // Assuming AddressHash is based on the `From` address
			Name:        metadata["name"].(string),
			Symbol:      metadata["symbol"].(string),
			TotalSupply: tokenEvent.Value,
			Decimals:    int(metadata["decimals"].(float64)),
		}

		// Save or update Token
		err := p.tokenRepository.SaveToken(ctx, tokenEntity)
		if err != nil {
			return fmt.Errorf("error saving/updating token: %w", err)
		}

		// Save contract if it's a mint event (first time interacting with this contract)
		contract := &smartcontract.SmartContract{
			AddressHash:     tokenEvent.From.String(),
			Name:            metadata["name"].(string),
			ABI:             metadata["abi"].(string), // Assuming ABI exists in metadata
			CompilerVersion: metadata["compiler_version"].(string),
			SourceCode:      metadata["source_code"].(string),
			VerifiedByEth:   true,
			EvmVersion:      "latest",
		}

		// Save the contract
		err = p.smartContractRepository.SaveSmartContract(ctx, contract)
		if err != nil {
			return fmt.Errorf("error saving contract: %w", err)
		}
	}

	// Process TokenInstance Entity for ERC-721 and ERC-1155 (if applicable)
	if tokenEvent.TokenId.String() != "0" { // ERC-721 or ERC-1155
		tokenInstance := &token.TokenInstance{
			TokenId:              tokenEvent.TokenId,
			TokenContractAddress: tokenEvent.From,
			OwnerAddress:         tokenEvent.To,
		}

		err := p.tokenRepository.SaveTokenInstance(ctx, tokenInstance)
		if err != nil {
			return fmt.Errorf("error saving token instance: %w", err)
		}
	}

	// Process TokenTransfer Entity
	tokenTransfer := &token.TokenTransfer{
		TransactionHash:      tokenEvent.TransactionHash.String(),
		LogIndex:             tokenEvent.LogIndex,
		From:                 tokenEvent.From,
		To:                   tokenEvent.To,
		TokenContractAddress: tokenEvent.From,
		Amount:               tokenEvent.Value,
	}

	err := p.tokenRepository.SaveTokenTransfer(ctx, tokenTransfer)
	if err != nil {
		return fmt.Errorf("error saving token transfer: %w", err)
	}

	return nil
}
