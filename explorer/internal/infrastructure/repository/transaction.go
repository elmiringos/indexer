package repository

import (
	"database/sql"

	"go.uber.org/zap"
)

type TransactionRepository struct {
	db  *sql.DB
	log *zap.Logger
}

func NewTransactionRepository(db *sql.DB, log *zap.Logger) *TransactionRepository {
	return &TransactionRepository{db: db, log: log}
}
