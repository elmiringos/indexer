package jwtsession

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, session *JwtSession) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetAll(ctx context.Context, expired bool) ([]*JwtSession, error)
	RevokeSession(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*JwtSession, error)
	GetActiveByUserEmail(ctx context.Context, email string) (*JwtSession, error)
	UpdateAccessTokenId(ctx context.Context, sessionId, accesTokenId uuid.UUID) error
}
