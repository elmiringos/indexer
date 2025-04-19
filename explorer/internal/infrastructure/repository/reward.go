package repository

import (
	"database/sql"
)

type RewardRepository struct {
	db *sql.DB
}

func NewRewardRepository(db *sql.DB) *RewardRepository {
	return &RewardRepository{db: db}
}
