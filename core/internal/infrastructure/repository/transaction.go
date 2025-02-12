package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/transaction"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
	"github.com/ethereum/go-ethereum/common"
)

type TransactionRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewTransactionRepository(db *sql.DB, redis *redis.Client) *TransactionRepository {
	return &TransactionRepository{db: db, redis: redis}
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

func (r *TransactionRepository) SaveTransactionHash(ctx context.Context, hash common.Hash) error {
	return r.redis.Set(hash.Hex(), []byte{})
}

func (r *TransactionRepository) DeleteTransactionHash(ctx context.Context, hash common.Hash) error {
	return r.redis.Delete(hash.Hex())
}

func (r *TransactionRepository) SaveTransactionLog(ctx context.Context, txLog *transaction.TransactionLog) error {
	query := `insert into transaction_log (transaction_hash, index, first_topic, second_topic, third_topic, fourth_topic, address_hash) values ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, txLog.TransactionHash, txLog.Index, txLog.FirstTopic, txLog.SecondTopic, txLog.ThirdTopic, txLog.FourthTopic, txLog.AddressHash)

	return err
}

func (r *TransactionRepository) SaveTransactionAction(ctx context.Context, txAction *transaction.TransactionAction) error {
	query := `insert into transaction_action (transaction_hash, log_index, data, address_hash, type) values ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, txAction.TransactionHash, txAction.LogIndex, txAction.Data, txAction.AddressHash, txAction.Type)

	return err
}

func (r *TransactionRepository) SaveTransactionLogIndex(ctx context.Context, transactionHash common.Hash, txLogIndex int) error {
	key := fmt.Sprintf("%s:%d", transactionHash.Hex(), txLogIndex)
	return r.redis.Set(key, []byte(""))
}

func (r *TransactionRepository) DeleteTransactionLogIndex(ctx context.Context, transactionHash common.Hash, txLogIndex int) error {
	key := fmt.Sprintf("%s:%d", transactionHash.Hex(), txLogIndex)
	return r.redis.Delete(key)
}

func (r *TransactionRepository) CheckTransactionExists(ctx context.Context, hash common.Hash) (bool, error) {
	exists, err := r.redis.Get(hash.Hex())
	if err != nil {
		return false, err
	}

	return exists != nil, nil
}
