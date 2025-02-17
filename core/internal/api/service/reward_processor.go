package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/reward"
	"github.com/elmiringos/indexer/indexer-core/internal/infrastructure/repository"
	"go.uber.org/zap"
)

type RewardProcessor struct {
	blockRepository  *repository.BlockRepository
	rewardRepository *repository.RewardRepository
	log              *zap.Logger
}

var (
	ErrFailedToUnmarshalReward           = errors.New("failed to unmarshal reward")
	ErrFailedToSaveReward                = errors.New("failed to save reward")
	ErrFailedToDeleteBlockHashReward     = errors.New("failed to delete block hash for reward")
	ErrBlockDoesNotExistForReward        = errors.New("block does not exist for reward")
	ErrFailedToCheckBlockExistsForReward = errors.New("failed to check if block exists for reward")
)

func NewRewardProcessor(blockRepository *repository.BlockRepository, rewardRepository *repository.RewardRepository, log *zap.Logger) *RewardProcessor {
	return &RewardProcessor{blockRepository: blockRepository, rewardRepository: rewardRepository, log: log}
}

func (p *RewardProcessor) Process(ctx context.Context, data []byte) error {
	reward := &reward.Reward{}
	if err := json.Unmarshal(data, reward); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUnmarshalReward, err)
	}

	blockExists, err := p.blockRepository.CheckBlockExistsForReward(ctx, reward.BlockHash)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCheckBlockExistsForReward, err)
	}

	if !blockExists {
		return fmt.Errorf("%w: %s", ErrBlockDoesNotExistForReward, reward.BlockHash)
	}

	if err := p.rewardRepository.SaveReward(ctx, reward); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSaveReward, err)
	}

	return nil
}
