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

func (r *TokenRepository) SaveToken(ctx context.Context, token *token.Token) (string, error) {
	query := `insert into token (address, name, symbol, decimals, total_supply, fiat_value, circulation_market_cap, holder_count) values ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, token.AddressHash, token.Name, token.Symbol, token.Decimals, token.TotalSupply, token.FiatValue, token.CirculationMarketCap, token.HolderCount)
	if err != nil {
		return "", err
	}

	return token.AddressHash, nil
}

func (r *TokenRepository) SaveTokenAddress(ctx context.Context, address string) error {
	return r.redis.Set(address, []byte(address))
}

func (r *TokenRepository) DeleteTokenAddress(ctx context.Context, address string) error {
	return r.redis.Delete(address)
}
