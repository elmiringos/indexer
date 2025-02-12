package main

import (
	"github.com/elmiringos/indexer/indexer-core/config"
	"github.com/elmiringos/indexer/indexer-core/internal/api"
	"github.com/elmiringos/indexer/indexer-core/pkg/logger"
	"github.com/elmiringos/indexer/indexer-core/pkg/postgres"
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"
)

func main() {
	cfg, err := config.NewDefaultConfig()
	if err != nil {
		panic("Error in loading config: " + err.Error())
	}

	logger := logger.New(cfg)

	db := postgres.NewPostgresConnection(cfg, logger)
	redis := redis.NewClient(cfg, logger)

	server := api.NewServer(cfg, db, redis, logger)
	server.Start()
}
