package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/elmiringos/indexer/indexer-core/internal/domain"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/token"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
	"github.com/ethereum/go-ethereum/common"
)

type TokenRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewTokenRepository(db *sql.DB, redis *redis.Client) *TokenRepository {
	return &TokenRepository{db: db, redis: redis}
}

func (r *TokenRepository) SaveToken(ctx context.Context, token *token.Token) error {
	query := `
		INSERT INTO token (address_hash, name, symbol, decimals, total_supply, fiat_value, circulation_market_cap)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, token.Address, token.Name, token.Symbol, token.Decimals, token.TotalSupply, token.FiatValue, token.CirculationMarketCap)

	return err
}

func (r *TokenRepository) IncreaseTokenSupply(ctx context.Context, addressHash common.Address, addSupply domain.BigInt) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var currentSupply domain.BigInt

	selectQuery := `
		SELECT total_supply FROM token WHERE address_hash = $1 FOR UPDATE;
	`
	row := tx.QueryRowContext(ctx, selectQuery, addressHash)

	if err = row.Scan(&currentSupply); err != nil {
		return err
	}

	newSupply := currentSupply.Sum(addSupply)

	updateQuery := `
		UPDATE token SET total_supply = $1 WHERE address_hash = $2;
	`
	if _, err = tx.ExecContext(ctx, updateQuery, newSupply, addressHash); err != nil {
		return err
	}

	return tx.Commit()
}

func (r *TokenRepository) DecreaseTokenSupply(ctx context.Context, addressHash common.Address, subSupply domain.BigInt) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var currentSupply domain.BigInt

	selectQuery := `
		SELECT total_supply FROM token WHERE address_hash = $1;
	`

	row := tx.QueryRowContext(ctx, selectQuery, addressHash)
	err = row.Scan(&currentSupply)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Ensure the result is not negative
	newSupply := currentSupply.Sub(subSupply)
	if newSupply.Cmp(domain.BigIntZero()) < 0 {
		tx.Rollback()
		return fmt.Errorf("cannot decrease supply below zero")
	}

	updateQuery := `
		UPDATE token SET total_supply = $1 WHERE address_hash = $2;
	`

	_, err = tx.ExecContext(ctx, updateQuery, newSupply, addressHash)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *TokenRepository) SaveOrUpdateTokenInstance(ctx context.Context, token *token.TokenInstance) error {
	query := `
		INSERT INTO token_instance (token_id, token_contract_address_hash, owner_address_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (token_id, token_contract_address_hash)
		DO UPDATE SET owner_address_hash = EXCLUDED.owner_address_hash
	`

	_, err := r.db.ExecContext(ctx, query, token.TokenId, token.TokenContractAddress, token.OwnerAddress)

	return err
}

func (r *TokenRepository) SaveTokenTransfer(ctx context.Context, token *token.TokenTransfer) error {
	query := `
		INSERT INTO token_transfer (transaction_hash, log_index, from_address, to_address, token_contract_address_hash, amount)
		VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx, query,
		token.TransactionHash,
		token.LogIndex,
		token.From,
		token.To,
		token.TokenContractAddress,
		token.Amount,
	)

	return err
}
