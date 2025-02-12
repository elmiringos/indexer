package internal_transaction

import "context"

type Repository interface {
	SaveInternalTransaction(ctx context.Context, tx *InternalTransaction) (string, error)
}
