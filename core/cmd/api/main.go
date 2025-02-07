package main

import (
	"github.com/elmiringos/indexer/indexer-core/config"
	"github.com/elmiringos/indexer/indexer-core/internal/api"
	"github.com/elmiringos/indexer/indexer-core/pkg/logger"
	"github.com/elmiringos/indexer/indexer-core/pkg/postgres"
)

func main() {
	cfg, err := config.NewDefaultConfig()
	if err != nil {
		panic("Error in loading config: " + err.Error())
	}

	logger := logger.New(cfg)

	db := postgres.NewPostgresConnection(cfg, logger)

	server := api.NewServer(cfg, db, logger)
	server.Start()
}
