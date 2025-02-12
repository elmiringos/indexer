package token

import "context"

type Repository interface {
	SaveToken(ctx context.Context, token *Token) (string, error)
	SaveTokenAddress(ctx context.Context, address string) error
	DeleteTokenAddress(ctx context.Context, address string) error
}
