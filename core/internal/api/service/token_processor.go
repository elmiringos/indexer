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

	if tokenEvent.SmartContractDeployed {
		p.log.Debug("Parsing metadata", zap.Any("metadata", metadata))

		decimals, ok := metadata["decimals"].(float64)
		if !ok {
			decimals = 0
		}

		p.log.Debug("token event", zap.Any("d", tokenEvent))

		// Process Token Entity
		tokenEntity := &token.Token{
			Address:     tokenEvent.Address,
			Name:        metadata["name"].(string),
			Symbol:      metadata["symbol"].(string),
			TotalSupply: tokenEvent.Value,
			Decimals:    int(decimals),
		}

		// Save or update Token
		err := p.tokenRepository.SaveToken(ctx, tokenEntity)
		if err != nil {
			return fmt.Errorf("error saving/updating token: %w", err)
		}

		contract := &smartcontract.SmartContract{
			AddressHash:     tokenEvent.Address,
			Name:            metadata["name"].(string),
			CompilerVersion: "not_imlemented",
			SourceCode:      metadata["smartcontract_bytecode"].(string),
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
			TokenContractAddress: tokenEvent.Address,
			OwnerAddress:         tokenEvent.To,
		}

		err := p.tokenRepository.SaveOrUpdateTokenInstance(ctx, tokenInstance)
		if err != nil {
			return fmt.Errorf("error saving token instance: %w", err)
		}
	}

	// Process TokenTransfer Entity
	tokenTransfer := &token.TokenTransfer{
		TransactionHash:      tokenEvent.TransactionHash,
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
