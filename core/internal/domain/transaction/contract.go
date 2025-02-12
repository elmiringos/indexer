package transaction

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type Repository interface {
	SaveTransaction(ctx context.Context, tx *Transaction) error
	SaveTransactionHash(ctx context.Context, hash common.Hash) error
	CheckTransactionExists(ctx context.Context, hash common.Hash) (bool, error)
	DeleteTransactionHash(ctx context.Context, hash common.Hash) error
	SaveTransactionLog(ctx context.Context, txLog *TransactionLog) error
	SaveTransactionLogIndex(ctx context.Context, transactionHash common.Hash, txLogIndex int) error
	DeleteTransactionLogIndex(ctx context.Context, transactionHash common.Hash, txLogIndex int) error
	SaveTransactionAction(ctx context.Context, txAction *TransactionAction) error
}
