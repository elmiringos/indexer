package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/elmiringos/indexer/indexer-core/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Client struct {
	client *redis.Client
}

func NewClient(cfg *config.Config, log *zap.Logger) *Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.URL,
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Sprintf("Error in creating redis client: %v", err))
	}

	log.Debug("Redis is healthy", zap.String("connection", cfg.Redis.URL))

	return &Client{client: client}
}

func (r *Client) Get(key string) ([]byte, error) {
	val, err := r.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r *Client) Set(key string, val []byte) error {
	return r.client.Set(context.Background(), key, val, 0).Err()
}

func (r *Client) SetWithTTL(key string, val []byte, exp time.Duration) error {
	return r.client.Set(context.Background(), key, val, exp).Err()
}

func (r *Client) Delete(key string) error {
	return r.client.Del(context.Background(), key).Err()
}

func (r *Client) Reset() error {
	return r.client.FlushDB(context.Background()).Err()
}

func (r *Client) Close() error {
	return r.client.Close()
}

func (r *Client) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
