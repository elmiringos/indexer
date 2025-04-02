package transaction

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type Repository interface {
	SaveTransaction(ctx context.Context, tx *Transaction) error
	SaveTransactionLog(ctx context.Context, txLog *TransactionLog) error
	SaveTransactionAction(ctx context.Context, txAction *TransactionAction) error
	SaveTransactionHashForAction(ctx context.Context, hash common.Hash, count int) error
	SaveTransactionHashForLog(ctx context.Context, hash common.Hash, count int) error
	CheckTransactionExistForLog(ctx context.Context, hash common.Hash) (bool, error)
	CheckTransactionExistForAction(ctx context.Context, hash common.Hash) (bool, error)
	DecrementTransactionActionCount(ctx context.Context, hash common.Hash) error
	DecrementTransactionLogCount(ctx context.Context, hash common.Hash) error
}
