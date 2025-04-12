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
	ErrFailedToUnmarshalTransactionLog        = errors.New("failed to unmarshal transaction log")
	ErrFailedToSaveTransaction                = errors.New("failed to save transaction")
	ErrFailedToSaveTransactionLog             = errors.New("failed to save transaction log")
	ErrFailedToSaveTransactionHash            = errors.New("failed to save transaction hash")
	ErrFailedToDecrementBlockHashTransaction  = errors.New("failed to decrement block hash for transaction")
	ErrFailedToDecrementTransactionHash       = errors.New("failed to decrement transaction log count for transaction")
	ErrBlockDoesNotExistForTransaction        = errors.New("block does not exist for transaction")
	ErrFailedToCheckBlockExistsForTransaction = errors.New("failed to check if block exists for transaction")
	ErrTransactionDoesNotExistForLog          = errors.New("transaction does not exist for transaction log")
	ErrFailedToCheckTransactionExistsForLog   = errors.New("failed to check if transaction exists for trancation log")
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

	if err := p.blockRepository.DecrementBlockHashTransactionCount(ctx, transaction.BlockHash); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToDecrementBlockHashTransaction, err)
	}

	p.log.Info("Successful decremented transaction count for block", zap.Any("block_hash", transaction.BlockHash))

	if err := p.transactionRepository.SaveTransactionHashForLog(ctx, transaction.Hash, transaction.LogsCount); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToSaveTransactionHash, err)
	}

	p.log.Info("Transaction Hash for Transction Log successfully save")

	return nil
}

type TransactionLogProcessor struct {
	transactionRepository transaction.Repository
	log                   *zap.Logger
}

func NewTransactionLogProcessor(
	transactionRepository transaction.Repository,
	log *zap.Logger,
) *TransactionLogProcessor {
	log.Info("Creating new transaction log processor")
	return &TransactionLogProcessor{
		transactionRepository: transactionRepository,
		log:                   log,
	}
}

func (p *TransactionLogProcessor) Process(ctx context.Context, data []byte) error {
	p.log.Info("Data", zap.String("json", string(data)))
	transactionLog := &transaction.TransactionLog{}
	if err := json.Unmarshal(data, transactionLog); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToUnmarshalTransactionLog, err)
	}

	transactionExist, err := p.transactionRepository.CheckTransactionExistForLog(ctx, transactionLog.TransactionHash)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFailedToCheckTransactionExistsForLog, err)
	}

	if !transactionExist {
		return fmt.Errorf("%w: %s", ErrTransactionDoesNotExistForLog, transactionLog.TransactionHash)
	}

	if err := p.transactionRepository.SaveTransactionLog(ctx, transactionLog); err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToSaveTransactionLog, err)
	}

	p.log.Info("Transaction Log saved successfully", zap.Uint("transaction_log_index", transactionLog.Index))

	if err := p.transactionRepository.DecrementTransactionLogCount(ctx, transactionLog.TransactionHash); err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToDecrementTransactionHash, err)
	}

	p.log.Info("Successful decremented transaction log count for transaction", zap.Any("transaction_hash", transactionLog.TransactionHash))

	return nil
}
