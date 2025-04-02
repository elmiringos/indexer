package repository

import (
	"context"
	"database/sql"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/token"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
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
		INSERT INTO token (address, name, symbol, decimals, total_supply, fiat_value, circulation_market_cap, holder_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, token.Address, token.Name, token.Symbol, token.Decimals, token.TotalSupply, token.FiatValue, token.CirculationMarketCap, token.HolderCount)

	return err
}

func (r *TokenRepository) SaveTokenInstance(ctx context.Context, token *token.TokenInstance) error {
	query := `
		INSERT INTO token_instance (token_id, token_contract_address_hash, owner_address_hash)
		VALUES ($1, $2, $3)
		ON CONFLICT (token_id) 
		DO NOTHING
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
