package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/withdrawal"
	"go.uber.org/zap"
)

var (
	ErrFailedToUnmarshalWithdrawal           = errors.New("failed to unmarshal withdrawal")
	ErrFailedToSaveWithdrawal                = errors.New("failed to save withdrawal")
	ErrFailedToDeleteBlockHashWithdrawal     = errors.New("failed to delete block hash for withdrawal")
	ErrBlockDoesNotExistForWithdrawal        = errors.New("block does not exist for withdrawal")
	ErrFailedToCheckBlockExistsForWithdrawal = errors.New("failed to check if block exists for withdrawal")
)

type WithdrawalProcessor struct {
	blockRepository      block.Repository
	withdrawalRepository withdrawal.Repository
	log                  *zap.Logger
}

func NewWithdrawalProcessor(blockRepository block.Repository, withdrawalRepository withdrawal.Repository, log *zap.Logger) *WithdrawalProcessor {
	return &WithdrawalProcessor{blockRepository: blockRepository, withdrawalRepository: withdrawalRepository, log: log}
}

func (p *WithdrawalProcessor) Process(ctx context.Context, data []byte) error {
	withdrawal := &withdrawal.Withdrawal{}
	if err := json.Unmarshal(data, withdrawal); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUnmarshalWithdrawal, err)
	}

	blockExists, err := p.blockRepository.CheckBlockExistsForWithdrawal(ctx, withdrawal.BlockHash)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCheckBlockExistsForWithdrawal, err)
	}

	if !blockExists {
		return fmt.Errorf("%w: %s", ErrBlockDoesNotExistForWithdrawal, withdrawal.BlockHash)
	}

	if err := p.withdrawalRepository.SaveWithdrawal(ctx, withdrawal); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSaveWithdrawal, err)
	}

	p.log.Info("Withdrawal saved successfully", zap.Any("withdrawal", withdrawal))

	return nil
}
