package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elmiringos/indexer/explorer/config"
	"github.com/elmiringos/indexer/explorer/internal/api/server"
	"github.com/elmiringos/indexer/explorer/internal/api/service"
	"github.com/elmiringos/indexer/explorer/internal/infrastructure/repository"
	"github.com/elmiringos/indexer/explorer/pkg/logger"
	"github.com/elmiringos/indexer/explorer/pkg/postgres"
	"go.uber.org/zap"
)

func Run() error {
	cfg, err := config.NewDefaultConfig()
	if err != nil {
		panic("Error in loading config: " + err.Error())
	}

	// Initialize logger
	log := logger.New(cfg)
	defer log.Sync()

	db := postgres.NewPostgresConnection(cfg, log)

	// Initialize Repositories
	blockRepository := repository.NewBlockRepository(db.GetDb(), log)

	// Initialize services
	blockService := service.NewBlockService(
		blockRepository,
		log,
	)

	// Initialize REST and gRPC servers
	httpServer := server.NewRESTServer(blockService, log)
	grpcServer := server.NewGRPCServer(blockService, log)

	// Initialize listeners
	grpcL, httpL, muxer, err := server.SetupListeners(cfg.Port)
	if err != nil {
		log.Fatal("Failed to setup listeners", zap.Error(err))
	}

	// Graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start servers
	go func() {
		log.Info("Starting gRPC server...")
		if err := grpcServer.Serve(grpcL); err != nil {
			log.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	go func() {
		log.Info("Starting HTTP server...")
		if err := http.Serve(httpL, httpServer); err != nil {
			log.Fatal("HTTP server failed", zap.Error(err))
		}
	}()

	go func() {
		log.Info("Starting connection multiplexer...")
		if err := muxer.Serve(); err != nil {
			log.Fatal("cmux failed", zap.Error(err))
		}
	}()

	// Wait for interrupt
	<-ctx.Done()
	log.Info("Shutdown signal received. Cleaning up...")

	// You may want to gracefully stop gRPC & HTTP here if needed
	return nil
}
