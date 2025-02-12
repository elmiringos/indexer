package repository

import (
	"context"
	"database/sql"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/internal_transaction"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
)

type InternalTransactionRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewInternalTransactionRepository(db *sql.DB, redis *redis.Client) *InternalTransactionRepository {
	return &InternalTransactionRepository{db: db, redis: redis}
}

func (r *InternalTransactionRepository) SaveInternalTransaction(ctx context.Context, tx *internal_transaction.InternalTransaction) (string, error) {
	query := `insert into internal_transaction (hash, block_hash, from_address, to_address, value, gas, gas_price, input, nonce, timestamp) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, tx.Hash, tx.BlockHash, tx.From, tx.To, tx.Value, tx.Gas, tx.Gas, tx.Input, tx.Nonce, tx.Timestamp)

	return tx.Hash, err
}
