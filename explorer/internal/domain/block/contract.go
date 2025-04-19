package block

import (
	"context"

	"github.com/elmiringos/indexer/explorer/internal/domain"
	"github.com/ethereum/go-ethereum/common"
)

type Repository interface {
	GetCurrentBlock(ctx context.Context) (*Block, error)
	GetBlock(ctx context.Context, blockNumber domain.BigInt, hash common.Hash) (*Block, error)
	GetBlocks(ctx context.Context, fromblockNumber domain.BigInt, toBlockNumber domain.BigInt) ([]*Block, error)
}
