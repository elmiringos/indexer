package repository

import (
	"context"
	"database/sql"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
)

type BlockRepository struct {
	db *sql.DB
}

func NewBlockRepository(db *sql.DB) *BlockRepository {
	return &BlockRepository{db: db}
}

func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (*block.Block, error) {
	query := `select * from block order by block_number desc limit 1`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var b block.Block
		err = rows.Scan(&b.Hash, &b.Number, &b.MinerHash, &b.ParentHash, &b.GasLimit, &b.GasUsed, &b.Nonce, &b.Size, &b.Difficulty, &b.Consensus, &b.BaseFeePerGas, &b.IsEmpty, &b.Timestamp)
		if err != nil {
			return nil, err
		}
		return &b, nil
	}

	return nil, ErrNotFound
}
