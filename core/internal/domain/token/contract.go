package token

import "context"

type Repository interface {
	SaveToken(ctx context.Context, token *Token) error
	SaveTokenInstance(ctx context.Context, tokenInstance *TokenInstance) error
	SaveTokenTransfer(ctx context.Context, tokenInstance *TokenTransfer) error
}
