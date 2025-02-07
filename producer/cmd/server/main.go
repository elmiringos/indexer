package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/elmiringos/indexer/producer/internal/blockchain"
	"github.com/elmiringos/indexer/producer/internal/server"
	grpccoreclient "github.com/elmiringos/indexer/producer/pkg/grpc_core_client"
	"github.com/elmiringos/indexer/producer/pkg/logger"
	"github.com/elmiringos/indexer/producer/pkg/rabbitmq"

	"go.uber.org/zap"
)

func main() {
	// config and logger creation
	cfg, err := config.NewDefaultConfig()
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
		if err := publisher.CloseConnection(); err != nil {
			log.Error("failed to close RabbitMQ channel", zap.Error(err))
		}
		if err := publisher.CloseConnection(); err != nil {
			log.Error("failed to close RabbitMQ connection", zap.Error(err))
		}
	}()

	coreClient, err := grpccoreclient.NewCoreClient(cfg.Server.CoreServiceUrl)
	if err != nil {
		log.Error("failed to create core client", zap.Error(err))
	}
	defer coreClient.Close()

	// tx, err := blockchainProcessor.GetTransaction("0x0b1498d09e7e28e3b95b2ad313ff132e5b9fb503165672151d210a84e67487f3")
	// if err != nil {
	// 	log.Error("failed to get transaction", zap.Error(err))
	// }

	// fmt.Println(tx)

	// trace, err := blockchainProcessor.GetTransactionTrace(tx)
	// if err != nil {
	// 	log.Error("failed to get transaction trace", zap.Error(err))
	// }

	// fmt.Println(trace)

	server := server.NewServer(blockchainProcessor, publisher, coreClient, cfg)
	server.StartBlockchainDataConsuming()
}
