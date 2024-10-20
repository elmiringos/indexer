package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/elmiringos/indexer/producer/internal/blockchain"
	"github.com/elmiringos/indexer/producer/internal/server"
	"github.com/elmiringos/indexer/producer/pkg/logger"
	"github.com/elmiringos/indexer/producer/pkg/rabbitmq"

	"go.uber.org/zap"
)

func main() {
	// config and logger creation
	cfg, err := loadConfig()
	if err != nil {
		panic(fmt.Errorf("error in reading config: %v", err))
	}

	log := logger.New(cfg)
	defer func() {
		if err := log.Sync(); err != nil && err.Error() != "sync /dev/stdout: inappropriate ioctl for device" {
			log.Fatal("failed to sync logger", zap.Error(err))
		}
	}()

	startProducerService(cfg, log)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	log.Info("Received signal, shutting down...", zap.Any("signal", sig))

	log.Info("Server exiting")
}

func startProducerService(cfg *config.Config, log *zap.Logger) {
	blockchainProcessor := blockchain.NewBlockchainProcessor(cfg)
	defer blockchainProcessor.CloseClients()

	publisher := rabbitmq.NewPublisher(cfg.RMQ.URL)
	defer func() {
		if err := publisher.CloseChannel(); err != nil {
			log.Error("failed to close RabbitMQ channel", zap.Error(err))
		}
		if err := publisher.CloseConnection(); err != nil {
			log.Error("failed to close RabbitMQ connection", zap.Error(err))
		}
	}()

	server := server.NewServer(blockchainProcessor, publisher, cfg)
	server.StartBlockchainDataConsuming()
}

func loadConfig() (*config.Config, error) {
	return config.NewDefaultConfig()
}
