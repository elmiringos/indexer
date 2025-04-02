package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
	"go.uber.org/zap"
)

var (
	ErrFailedToUnmarshalBlock = errors.New("failed to unmarshal block")
	ErrFailedToSaveBlock      = errors.New("failed to save block")
	ErrFailedToSaveBlockHash  = errors.New("failed to save block hash")
)

type BlockProcessor struct {
	blockRepository block.Repository
	log             *zap.Logger
}

func NewBlockProcessor(blockRepository block.Repository, log *zap.Logger) *BlockProcessor {
	log.Info("Creating new block processor")
	return &BlockProcessor{
		blockRepository: blockRepository,
		log:             log,
	}
}

func (p *BlockProcessor) Process(ctx context.Context, data []byte) error {
	block := &block.Block{}
	if err := json.Unmarshal(data, block); err != nil {
		return fmt.Errorf("%s: %w", ErrFailedToUnmarshalBlock, err)
	}

	if err := p.blockRepository.SaveBlock(ctx, block); err != nil {
		return fmt.Errorf("%s: %w", ErrFailedToSaveBlock, err)
	}

	p.log.Info("Block saved successfully", zap.Any("block_number", block.Number))

	// Save block hash for transaction, withdrawal and reward for sync queues
	if err := p.blockRepository.SaveBlockHashForTransaction(ctx, block.Hash, block.TransactionsCount); err != nil {
		return fmt.Errorf("%s: %w", ErrFailedToSaveBlockHash, err)
	}

	if err := p.blockRepository.SaveBlockHashForWithdrawal(ctx, block.Hash, block.WithdrawalsCount); err != nil {
		return fmt.Errorf("%s: %w", ErrFailedToSaveBlockHash, err)
	}

	if err := p.blockRepository.SaveBlockHashForReward(ctx, block.Hash, 1); err != nil {
		return fmt.Errorf("%s: %w", ErrFailedToSaveBlockHash, err)
	}

	p.log.Info("Block hash saved successfully", zap.Any("block_hash", block.Hash))
	return nil
}
