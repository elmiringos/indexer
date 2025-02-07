package repository

import "database/sql"

type SmartContractRepository struct {
	db *sql.DB
}

func NewSmartContractRepository(db *sql.DB) *SmartContractRepository {
	return &SmartContractRepository{db: db}
}
