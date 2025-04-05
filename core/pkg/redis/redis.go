package redis

import (
	"context"
	"fmt"
	"strconv"
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

func (r *Client) GetInt(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return intVal, nil
}

func (r *Client) SetInt(ctx context.Context, key string, count int) error {
	return r.client.SetNX(ctx, key, count, 0).Err()
}

var decrAndMaybeDeleteScript = redis.NewScript(`
    local val = redis.call("DECR", KEYS[1])
    if val <= 0 then
        redis.call("DEL", KEYS[1])
    end
    return val
`)

func (r *Client) DecrementAndMaybeDelete(ctx context.Context, key string) error {
	_, err := decrAndMaybeDeleteScript.Run(ctx, r.client, []string{key}).Result()
	return err
}

func (r *Client) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
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
