package api

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/elmiringos/indexer/indexer-core/config"
	"github.com/elmiringos/indexer/indexer-core/internal/api/handler"
	"github.com/elmiringos/indexer/indexer-core/internal/api/pb"
	"github.com/elmiringos/indexer/indexer-core/internal/api/service"
	"github.com/elmiringos/indexer/indexer-core/internal/infrastructure/repository"
	"github.com/elmiringos/indexer/indexer-core/pkg/postgres"
	"github.com/elmiringos/indexer/indexer-core/pkg/rabbitmq"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg        *config.Config
	db         *postgres.Connection
	logger     *zap.Logger
	handler    *handler.CoreHandler
	subscriber *rabbitmq.Subscriber
}

func NewServer(cfg *config.Config, db *postgres.Connection, logger *zap.Logger) *Server {
	// Initialize repositories
	blockRepository := repository.NewBlockRepository(db.GetDb())
	internalTransactionRepository := repository.NewInternalTransactionRepository(db.GetDb())
	rewardRepository := repository.NewRewardRepository(db.GetDb())
	smartContractRepository := repository.NewSmartContractRepository(db.GetDb())
	tokenRepository := repository.NewTokenRepository(db.GetDb())
	transactionRepository := repository.NewTransactionRepository(db.GetDb())

	// Initialize subscriber
	subscriber := rabbitmq.NewSubscriber(cfg.RMQ.URL)

	// Initialize service
	coreService := service.NewCoreService(
		logger,
		blockRepository,
		internalTransactionRepository,
		rewardRepository,
		smartContractRepository,
		tokenRepository,
		transactionRepository,
	)

	// Initialize handler
	coreHandler := handler.NewCoreHandler(coreService, logger)

	return &Server{
		cfg:        cfg,
		db:         db,
		logger:     logger,
		handler:    coreHandler,
		subscriber: subscriber,
	}
}

func (s *Server) Start() {
	address := fmt.Sprintf(":%s", s.cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCoreServiceServer(grpcServer, s.handler)
	reflection.Register(grpcServer)

	s.logger.Info("gRPC Server started on " + address)
	if err := grpcServer.Serve(listener); err != nil {
		s.logger.Fatal("failed to serve", zap.Error(err))
	}

	// Handle OS shutdown signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		s.logger.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	s.logger.Info("gRPC Server started on " + address)
	if err := grpcServer.Serve(listener); err != nil {
		s.logger.Fatal("failed to serve", zap.Error(err))
	}
}
