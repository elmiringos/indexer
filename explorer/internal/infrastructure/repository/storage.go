package repository

import "context"

type KVStorage interface {
	GetInt(ctx context.Context, key string) (int, error)
	SetInt(ctx context.Context, key string, value int) error
	Delete(ctx context.Context, key string) error
	DecrementAndMaybeDelete(ctx context.Context, key string) error
}
