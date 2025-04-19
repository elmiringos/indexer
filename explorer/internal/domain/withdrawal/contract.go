package withdrawal

import "context"

type Repository interface {
	SaveWithdrawal(ctx context.Context, withdrawal *Withdrawal) error
}
