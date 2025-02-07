package block

import (
	"context"
)

type Repository interface {
	GetCurrentBlock(ctx context.Context) (*Block, error)
}
