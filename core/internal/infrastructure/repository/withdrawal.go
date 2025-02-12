package repository

import (
	"context"
	"database/sql"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/withdrawal"
)

type WithdrawalRepository struct {
	db *sql.DB
}

func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{db: db}
}

func (r *WithdrawalRepository) SaveWithdrawal(ctx context.Context, withdrawal *withdrawal.Withdrawal) error {
	query := `
		INSERT INTO withdrawal (
			index,
			block_hash,
			address_hash,
			validator_index,
			amount
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)
	`
	_, err := r.db.ExecContext(ctx, query, withdrawal.Index, withdrawal.BlockHash, withdrawal.AddressHash, withdrawal.ValidatorIndex, withdrawal.Amount)
	return err
}
