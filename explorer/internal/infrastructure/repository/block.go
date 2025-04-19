package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/elmiringos/indexer/explorer/internal/domain"
	"github.com/elmiringos/indexer/explorer/internal/domain/block"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type BlockRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewBlockRepository(db *sql.DB, log *zap.Logger) *BlockRepository {
	return &BlockRepository{db: db, log: log}
}

func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (*block.Block, error) {
	var block block.Block

	query := `
		SELECT 
			hash,
			number,
			miner_hash,
			parent_hash,
			gas_limit,
			gas_used, 
			nonce,
			size,
			difficulty,
			is_pos,
			base_fee_per_gas,
			timestamp
		FROM block 
		ORDER BY number::numeric DESC LIMIT 1`

	row := r.db.QueryRowContext(ctx, query)
	err := row.Scan(
		&block.Hash,
		&block.Number,
		&block.MinerHash,
		&block.ParentHash,
		&block.GasLimit,
		&block.GasUsed,
		&block.Nonce,
		&block.Size,
		&block.Difficulty,
		&block.IsPos,
		&block.BaseFeePerGas,
		&block.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Info("No blocks found")
			return nil, nil
		}
		return nil, err
	}

	return &block, nil
}

func (r *BlockRepository) GetBlock(ctx context.Context, blockNumber domain.BigInt, hash common.Hash) (*block.Block, error) {
	var block block.Block

	blockNumberStr := blockNumber.String()

	if blockNumberStr == "" && hash == (common.Hash{}) {
		return nil, errors.New("block number and hash cannot be empty")
	}

	var query string
	var row *sql.Row

	if blockNumberStr == "" {
		query = `
			SELECT hash, number, miner_hash, parent_hash, gas_limit, gas_used, nonce, size, difficulty, is_pos, base_fee_per_gas, timestamp
			FROM blocks WHERE hash = $1`

		row = r.db.QueryRowContext(ctx, query, hash.String())
	} else if hash == (common.Hash{}) {
		query = `
			SELECT hash, number, miner_hash, parent_hash, gas_limit, gas_used, nonce, size, difficulty, is_pos, base_fee_per_gas, timestamp
			FROM blocks WHERE number = $1`

		row = r.db.QueryRowContext(ctx, query, blockNumberStr)
	}

	err := row.Scan(
		&block.Number,
		&block.Hash,
		&block.ParentHash,
		&block.MinerHash,
		&block.GasLimit,
		&block.GasUsed,
		&block.Nonce,
		&block.Size,
		&block.Difficulty,
		&block.IsPos,
		&block.BaseFeePerGas,
		&block.Timestamp,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Info("No blocks found")
			return nil, nil
		}
		return nil, err
	}

	return &block, nil
}

func (r *BlockRepository) GetBlocks(ctx context.Context, fromBlockNumber domain.BigInt, toBlockNumber domain.BigInt) ([]*block.Block, error) {

	var blocks []*block.Block

	query := `
		SELECT hash, number, miner_hash, parent_hash, gas_limit, gas_used, nonce, size, difficulty, is_pos, base_fee_per_gas, timestamp
		FROM blocks WHERE number::numeric BETWEEN $1 AND $2 ORDER BY number ASC`

	rows, err := r.db.QueryContext(ctx, query, fromBlockNumber.String(), toBlockNumber.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var block block.Block
		err := rows.Scan(
			&block.Number,
			&block.Hash,
			&block.ParentHash,
			&block.MinerHash,
			&block.GasLimit,
			&block.GasUsed,
			&block.Nonce,
			&block.Size,
			&block.Difficulty,
			&block.IsPos,
			&block.BaseFeePerGas,
			&block.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, &block)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		r.log.Info("No blocks found")
		return nil, nil
	}
	return blocks, nil
}
