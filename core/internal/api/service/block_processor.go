package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
	"go.uber.org/zap"
)

var errFailedToUnmarshalBlock = errors.New("failed to unmarshal block")
var errFailedToSaveBlock = errors.New("failed to save block")
var errFailedToSaveBlockHash = errors.New("failed to save block hash")

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
		return fmt.Errorf("%s: %w", errFailedToUnmarshalBlock, err)
	}

	if err := p.blockRepository.SaveBlock(ctx, block); err != nil {
		return fmt.Errorf("%s: %w", errFailedToSaveBlock, err)
	}

	p.log.Info("Block saved successfully", zap.Any("block_number", block.Number))

	if err := p.blockRepository.SaveBlockHash(ctx, block.Hash); err != nil {
		return fmt.Errorf("%s: %w", errFailedToSaveBlockHash, err)
	}

	p.log.Info("Block hash saved successfully", zap.Any("block_hash", block.Hash))
	return nil
}
