package token

import (
	"context"

	"github.com/elmiringos/indexer/indexer-core/internal/domain"
	"github.com/ethereum/go-ethereum/common"
)

type Repository interface {
	SaveToken(ctx context.Context, token *Token) error
	IncreaseTokenSupply(ctx context.Context, addressHash common.Address, addSupply domain.BigInt) error
	DecreaseTokenSupply(ctx context.Context, addressHash common.Address, subSupply domain.BigInt) error
	SaveOrUpdateTokenInstance(ctx context.Context, tokenInstance *TokenInstance) error
	SaveTokenTransfer(ctx context.Context, tokenInstance *TokenTransfer) error
}
