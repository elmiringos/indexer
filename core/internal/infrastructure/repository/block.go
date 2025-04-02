package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type BlockRepository struct {
	db    *sql.DB
	store KVStorage
	log   *zap.Logger
}

func NewBlockRepository(db *sql.DB, store KVStorage, log *zap.Logger) *BlockRepository {
	return &BlockRepository{db: db, store: store, log: log}
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

func makeTransactionKey(hash common.Hash) string {
	return fmt.Sprintf("block:%s:transaction", hash.Hex())
}

func makeWithdrawalKey(hash common.Hash) string {
	return fmt.Sprintf("block:%s:withdrawal", hash.Hex())
}

func makeRewardKey(hash common.Hash) string {
	return fmt.Sprintf("block:%s:reward", hash.Hex())
}

func (r *BlockRepository) SaveBlockHashForTransaction(ctx context.Context, hash common.Hash, transactionCount int) error {
	r.log.Debug("Saving block hash for transaction", zap.String("hash", hash.Hex()))
	return r.store.SetInt(ctx, makeTransactionKey(hash), transactionCount)
}

func (r *BlockRepository) SaveBlockHashForWithdrawal(ctx context.Context, hash common.Hash, withdrawalCount int) error {
	r.log.Debug("Saving block hash for withdrawal", zap.String("hash", hash.Hex()))
	return r.store.SetInt(ctx, makeWithdrawalKey(hash), withdrawalCount)
}

func (r *BlockRepository) SaveBlockHashForReward(ctx context.Context, hash common.Hash, rewardCount int) error {
	r.log.Debug("Saving block hash for reward", zap.String("hash", hash.Hex()))
	return r.store.SetInt(ctx, makeRewardKey(hash), rewardCount)
}

func (r *BlockRepository) DeleteBlockHashForTransaction(ctx context.Context, hash common.Hash) error {
	r.log.Debug("Deleting block hash for transaction", zap.String("hash", hash.Hex()))
	return r.store.Delete(ctx, makeTransactionKey(hash))
}

func (r *BlockRepository) DeleteBlockHashForWithdrawal(ctx context.Context, hash common.Hash) error {
	r.log.Debug("Deleting block hash for withdrawal", zap.String("hash", hash.Hex()))
	return r.store.Delete(ctx, makeWithdrawalKey(hash))
}

func (r *BlockRepository) DeleteBlockHashForReward(ctx context.Context, hash common.Hash) error {
	r.log.Debug("Deleting block hash for reward", zap.String("hash", hash.Hex()))
	return r.store.Delete(ctx, makeRewardKey(hash))
}

func (r *BlockRepository) CheckBlockExistsForTransaction(ctx context.Context, hash common.Hash) (bool, error) {
	exists, err := r.store.GetInt(ctx, makeTransactionKey(hash))
	if err != nil {
		return false, err
	}

	return exists != 0, nil
}

func (r *BlockRepository) CheckBlockExistsForWithdrawal(ctx context.Context, hash common.Hash) (bool, error) {
	exists, err := r.store.GetInt(ctx, makeWithdrawalKey(hash))
	if err != nil {
		return false, err
	}

	return exists != 0, nil
}

func (r *BlockRepository) CheckBlockExistsForReward(ctx context.Context, hash common.Hash) (bool, error) {
	exists, err := r.store.GetInt(ctx, makeRewardKey(hash))
	if err != nil {
		return false, err
	}

	return exists != 0, nil
}

func (r *BlockRepository) DecrementBlockHashRewardCount(ctx context.Context, hash common.Hash) error {
	return r.store.DecrementAndMaybeDelete(ctx, makeRewardKey(hash))
}

func (r *BlockRepository) DecrementBlockHashTransactionCount(ctx context.Context, hash common.Hash) error {
	return r.store.DecrementAndMaybeDelete(ctx, makeTransactionKey(hash))
}

func (r *BlockRepository) DecrementBlockHashWithdrawalCount(ctx context.Context, hash common.Hash) error {
	return r.store.DecrementAndMaybeDelete(ctx, makeWithdrawalKey(hash))
}
