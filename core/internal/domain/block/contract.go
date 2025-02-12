package block

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

type Repository interface {
	GetCurrentBlock(ctx context.Context) (*Block, error)
	SaveBlock(ctx context.Context, b *Block) error
	SaveBlockHash(ctx context.Context, hash common.Hash) error
	CheckBlockExists(ctx context.Context, hash common.Hash) (bool, error)
	DeleteBlockHash(ctx context.Context, hash common.Hash) error
}
