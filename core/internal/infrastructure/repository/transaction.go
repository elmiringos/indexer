package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/transaction"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type TransactionRepository struct {
	db    *sql.DB
	store KVStorage
	log   *zap.Logger
}

func NewTransactionRepository(db *sql.DB, store KVStorage, log *zap.Logger) *TransactionRepository {
	return &TransactionRepository{db: db, store: store, log: log}
}

func (r *TransactionRepository) SaveTransaction(ctx context.Context, tx *transaction.Transaction) error {
	query := `insert into transaction (
		hash,
		block_hash,
		index,
		status,
		gas,
		gas_used,
		input,
		value,
		from_address,
		to_address,
		nonce,
		timestamp
	) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query, tx.Hash, tx.BlockHash, tx.Index, tx.Status, tx.Gas, tx.GasUsed, tx.Input, tx.Value, tx.From, tx.To, tx.Nonce, tx.Timestamp)

	return err
}

func (r *TransactionRepository) SaveTransactionLog(ctx context.Context, txLog *transaction.TransactionLog) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	logQuery := `
		INSERT INTO transaction_log (
			address, transaction_hash, block_hash, transaction_index, log_index, data  
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, logQuery,
		txLog.Address,
		txLog.TransactionHash,
		txLog.BlockHash,
		txLog.TransactionIndex,
		txLog.Index,
		txLog.Data,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	topicQuery := `
		INSERT INTO transaction_log_topic (
			transaction_hash, log_index, topic_index, topic
		) VALUES ($1, $2, $3, $4)`

	for i, topic := range txLog.Topics {
		_, err := tx.ExecContext(ctx, topicQuery,
			txLog.TransactionHash,
			txLog.Index,
			i,
			topic,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *TransactionRepository) SaveTransactionAction(ctx context.Context, txAction *transaction.TransactionAction) error {
	query := `insert into transaction_action (transaction_hash, log_index, data, address, type) values ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, txAction.TransactionHash, txAction.LogIndex, txAction.Data, txAction.AddressHash, txAction.Type)

	return err
}

func makeTransactionLogKey(hash common.Hash) string {
	return fmt.Sprintf("transaction:%s:logs", hash)
}

func makeTransactionActionKey(hash common.Hash) string {
	return fmt.Sprintf("transaction:%s:actions", hash)
}

func (r *TransactionRepository) SaveTransactionHashForLog(ctx context.Context, hash common.Hash, count int) error {
	return r.store.SetInt(ctx, makeTransactionLogKey(hash), count)
}

func (r *TransactionRepository) CheckTransactionExistForLog(ctx context.Context, hash common.Hash) (bool, error) {
	exists, err := r.store.GetInt(ctx, makeTransactionLogKey(hash))
	if err != nil {
		return false, err
	}

	return exists != 0, nil
}

func (r *TransactionRepository) DecrementTransactionLogCount(ctx context.Context, hash common.Hash) error {
	return r.store.DecrementAndMaybeDelete(ctx, makeTransactionLogKey(hash))
}

func (r *TransactionRepository) SaveTransactionHashForAction(ctx context.Context, hash common.Hash, count int) error {
	return r.store.SetInt(ctx, makeTransactionActionKey(hash), count)
}

func (r *TransactionRepository) CheckTransactionExistForAction(ctx context.Context, hash common.Hash) (bool, error) {
	exists, err := r.store.GetInt(ctx, makeTransactionActionKey(hash))
	if err != nil {
		return false, err
	}

	return exists != 0, nil
}

func (r *TransactionRepository) DecrementTransactionActionCount(ctx context.Context, hash common.Hash) error {
	return r.store.DecrementAndMaybeDelete(ctx, makeTransactionActionKey(hash))
}
