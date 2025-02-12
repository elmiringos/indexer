package repository

import (
	"context"
	"database/sql"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
	"github.com/ethereum/go-ethereum/common"
)

type BlockRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewBlockRepository(db *sql.DB, redis *redis.Client) *BlockRepository {
	return &BlockRepository{db: db, redis: redis}
}

func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (*block.Block, error) {
	query := `select hash, number, miner_hash, parent_hash, gas_limit, gas_used, nonce, size, difficulty, is_pos, base_fee_per_gas, timestamp from block order by number desc limit 1`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var b block.Block
		err = rows.Scan(&b.Hash, &b.Number, &b.MinerHash, &b.ParentHash, &b.GasLimit, &b.GasUsed, &b.Nonce, &b.Size, &b.Difficulty, &b.IsPos, &b.BaseFeePerGas, &b.Timestamp)
		if err != nil {
			return nil, err
		}
		return &b, nil
	}

	return nil, ErrNotFound
}

func (r *BlockRepository) SaveBlock(ctx context.Context, b *block.Block) error {
	query := `insert into block (hash, number, miner_hash, parent_hash, gas_limit, gas_used, nonce, size, difficulty, is_pos, base_fee_per_gas, timestamp) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err := r.db.ExecContext(ctx, query, b.Hash, b.Number, b.MinerHash, b.ParentHash, b.GasLimit, b.GasUsed, b.Nonce, b.Size, b.Difficulty, b.IsPos, b.BaseFeePerGas, b.Timestamp)
	return err
}

func (r *BlockRepository) SaveBlockHash(ctx context.Context, hash common.Hash) error {
	return r.redis.Set(hash.Hex(), []byte(hash.Hex()))
}

func (r *BlockRepository) DeleteBlockHash(ctx context.Context, hash common.Hash) error {
	return r.redis.Delete(hash.Hex())
}

func (r *BlockRepository) CheckBlockExists(ctx context.Context, hash common.Hash) (bool, error) {
	exists, err := r.redis.Get(hash.Hex())
	if err != nil {
		return false, err
	}

	return exists != nil, nil
}
