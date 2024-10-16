package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/elmiringos/indexer/block-producer/config"
	"github.com/elmiringos/indexer/block-producer/pkg/logger"

	"go.uber.org/zap"
)

func main() {
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	sig := <-c
	log.Info("Received signal, shutting down...", zap.Any("signal", sig))

	log.Info("Server exiting")
}

// loadConfig loads the configuration and handles errors
func loadConfig() (*config.Config, error) {
	return config.NewDefaultConfig()
}
