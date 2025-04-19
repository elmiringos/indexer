package service

import (
	"context"

	"github.com/elmiringos/indexer/explorer/internal/domain/block"
	"go.uber.org/zap"
)

type BlockService struct {
	Blockrepository block.Repository
	logger          *zap.Logger
}

func NewBlockService(blockRepository block.Repository, logger *zap.Logger) *BlockService {
	return &BlockService{
		Blockrepository: blockRepository,
		logger:          logger,
	}
}

func (s *BlockService) GetCurrentBlock() (*block.Block, error) {
	block, err := s.Blockrepository.GetCurrentBlock(context.Background())
	if err != nil {
		s.logger.Error("Failed to get current block", zap.Error(err))
		return nil, err
	}

	return block, nil
}
