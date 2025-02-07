package repository

import "database/sql"

type InternalTransactionRepository struct {
	db *sql.DB
}

func NewInternalTransactionRepository(db *sql.DB) *InternalTransactionRepository {
	return &InternalTransactionRepository{db: db}
}
