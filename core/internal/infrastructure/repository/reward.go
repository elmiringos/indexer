package repository

import (
	"context"
	"database/sql"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/reward"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
)

type RewardRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewRewardRepository(db *sql.DB, redis *redis.Client) *RewardRepository {
	return &RewardRepository{db: db, redis: redis}
}

func (r *RewardRepository) SaveReward(ctx context.Context, reward *reward.Reward) error {
	query := `insert into reward (block_hash, address, amount) values ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, reward.BlockHash, reward.Address, reward.Amount)

	return err
}
