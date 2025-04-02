package repository

import (
	"context"
	"database/sql"

	smartcontract "github.com/elmiringos/indexer/indexer-core/internal/domain/smart_contract"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
)

type SmartContractRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewSmartContractRepository(db *sql.DB, redis *redis.Client) *SmartContractRepository {
	return &SmartContractRepository{db: db, redis: redis}
}

func (r *SmartContractRepository) SaveSmartContract(ctx context.Context, smartContract *smartcontract.SmartContract) error {
	query := `insert into smart_contract (hash, block_hash, from_address, to_address, value, gas, gas_price, input, nonce, timestamp) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, smartContract.AddressHash, smartContract.Name, smartContract.CompilerVersion, smartContract.SourceCode, smartContract.ABI, smartContract.CompilerSettings, smartContract.VerifiedByEth, smartContract.EvmVersion)

	return err
}
