package block

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type Repository interface {
	GetCurrentBlock(ctx context.Context) (*Block, error)
	SaveBlock(ctx context.Context, b *Block) error
	SaveBlockHashForTransaction(ctx context.Context, hash common.Hash, transactionCount int) error
	SaveBlockHashForWithdrawal(ctx context.Context, hash common.Hash, withdrawalCount int) error
	SaveBlockHashForReward(ctx context.Context, hash common.Hash, rewardCount int) error
	DeleteBlockHashForTransaction(ctx context.Context, hash common.Hash) error
	DeleteBlockHashForWithdrawal(ctx context.Context, hash common.Hash) error
	DeleteBlockHashForReward(ctx context.Context, hash common.Hash) error
	DecrementBlockHashTransactionCount(ctx context.Context, hash common.Hash) error
	DecrementBlockHashWithdrawalCount(ctx context.Context, hash common.Hash) error
	DecrementBlockHashRewardCount(ctx context.Context, hash common.Hash) error
	CheckBlockExistsForTransaction(ctx context.Context, hash common.Hash) (bool, error)
	CheckBlockExistsForWithdrawal(ctx context.Context, hash common.Hash) (bool, error)
	CheckBlockExistsForReward(ctx context.Context, hash common.Hash) (bool, error)
}
