package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/transaction"

	"go.uber.org/zap"
)

var (
	ErrFailedToUnmarshalTransaction           = errors.New("failed to unmarshal transaction")
	ErrFailedToSaveTransaction                = errors.New("failed to save transaction")
	ErrFailedToSaveTransactionHash            = errors.New("failed to save transaction hash")
	ErrFailedToDeleteBlockHashTransaction     = errors.New("failed to delete block hash for transaction")
	ErrBlockDoesNotExistForTransaction        = errors.New("block does not exist for transaction")
	ErrFailedToCheckBlockExistsForTransaction = errors.New("failed to check if block exists for transaction")
)

type TransactionProcessor struct {
	blockRepository       block.Repository
	transactionRepository transaction.Repository
	log                   *zap.Logger
}

func NewTransactionProcessor(
	blockRepository block.Repository,
	transactionRepository transaction.Repository,
	log *zap.Logger,
) *TransactionProcessor {
	log.Info("Creating new transaction processor")
	return &TransactionProcessor{
		blockRepository:       blockRepository,
		transactionRepository: transactionRepository,
		log:                   log,
	}
}

// TODO: Think about distributed transaction for atomicity to multiple resources (psql, redis)
func (p *TransactionProcessor) Process(ctx context.Context, data []byte) error {
	transaction := &transaction.Transaction{}
	if err := json.Unmarshal(data, transaction); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUnmarshalTransaction, err)
	}

	blockExists, err := p.blockRepository.CheckBlockExistsForTransaction(ctx, transaction.BlockHash)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCheckBlockExistsForTransaction, err)
	}

	if !blockExists {
		return fmt.Errorf("%w: %s", ErrBlockDoesNotExistForTransaction, transaction.BlockHash)
	}

	if err := p.transactionRepository.SaveTransaction(ctx, transaction); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSaveTransaction, err)
	}

	p.log.Info("Transaction saved successfully", zap.Any("transaction_hash", transaction.Hash))

	if err := p.transactionRepository.SaveTransactionHash(ctx, transaction.Hash); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSaveTransactionHash, err)
	}

	p.log.Info("Transaction hash saved successfully", zap.Any("transaction_hash", transaction.Hash.Hex()))

	return nil
}
