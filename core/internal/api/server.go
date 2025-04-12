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
	"github.com/elmiringos/indexer/indexer-core/pkg/redis"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg      *config.Config
	db       *postgres.Connection
	log      *zap.Logger
	handler  *handler.CoreHandler
	consumer *rabbitmq.Consumer
}

// NewServer initializes the gRPC server
func NewServer(cfg *config.Config, db *postgres.Connection, redis *redis.Client, logger *zap.Logger) *Server {
	// Initialize repositories
	blockRepository := repository.NewBlockRepository(db.GetDb(), redis, logger)
	internalTransactionRepository := repository.NewInternalTransactionRepository(db.GetDb(), redis)
	rewardRepository := repository.NewRewardRepository(db.GetDb(), redis)
	smartContractRepository := repository.NewSmartContractRepository(db.GetDb(), redis)
	tokenRepository := repository.NewTokenRepository(db.GetDb(), redis)
	transactionRepository := repository.NewTransactionRepository(db.GetDb(), redis, logger)
	withdrawalRepository := repository.NewWithdrawalRepository(db.GetDb())

	// Initialize consumer to message broker
	consumer := rabbitmq.NewConsumer(cfg.RMQ.URL)
	initializeQueues(consumer, logger)

	// Initialize queue channels and get messages
	blockMessages := consumer.Consume(consumer.CreateChannel(), rabbitmq.BlockQueue)
	transactionMessages := consumer.Consume(consumer.CreateChannel(), rabbitmq.TransactionQueue)
	transactionLogMessages := consumer.Consume(consumer.CreateChannel(), rabbitmq.TransactionLogQueue)
	rewardMessages := consumer.Consume(consumer.CreateChannel(), rabbitmq.RewardQueue)
	withdrawalMessages := consumer.Consume(consumer.CreateChannel(), rabbitmq.WithdrawalQueue)
	tokenEventMessages := consumer.Consume(consumer.CreateChannel(), rabbitmq.TokenEventQueue)

	// Initialize service
	coreService := service.NewCoreService(
		logger,
		blockRepository,
		internalTransactionRepository,
		rewardRepository,
		smartContractRepository,
		tokenRepository,
		transactionRepository,
		withdrawalRepository,
	)

	// Initialize handler
	coreHandler := handler.NewCoreHandler(coreService, logger)

	// Initialize worker pools and processors
	// Block processor
	blockProcessor := service.NewBlockProcessor(blockRepository, logger)
	blockWorkerPool := service.NewWorkerPool(blockProcessor, logger, cfg.Server.Worker)
	go blockWorkerPool.Start(blockMessages)

	// Transaction processor
	transactionProcessor := service.NewTransactionProcessor(blockRepository, transactionRepository, logger)
	transactionWorkerPool := service.NewWorkerPool(transactionProcessor, logger, cfg.Server.Worker)
	go transactionWorkerPool.Start(transactionMessages)

	// // Transaction Log processor
	transactionLogProcessor := service.NewTransactionLogProcessor(transactionRepository, logger)
	transactionLogWorkerPool := service.NewWorkerPool(transactionLogProcessor, logger, cfg.Server.Worker)
	go transactionLogWorkerPool.Start(transactionLogMessages)

	// Reward processor
	rewardProcessor := service.NewRewardProcessor(blockRepository, rewardRepository, logger)
	rewardWorkerPool := service.NewWorkerPool(rewardProcessor, logger, cfg.Server.Worker)
	go rewardWorkerPool.Start(rewardMessages)

	// Withdrawal processor
	withdrawalProcessor := service.NewWithdrawalProcessor(blockRepository, withdrawalRepository, logger)
	withdrawalWorkerPool := service.NewWorkerPool(withdrawalProcessor, logger, cfg.Server.Worker)
	go withdrawalWorkerPool.Start(withdrawalMessages)

	// Token event processor
	tokenEventProcessor := service.NewTokenProccesor(tokenRepository, smartContractRepository, logger)
	tokenEventWorkerPool := service.NewWorkerPool(tokenEventProcessor, logger, cfg.Server.Worker)
	go tokenEventWorkerPool.Start(tokenEventMessages)

	return &Server{
		cfg:      cfg,
		db:       db,
		log:      logger,
		handler:  coreHandler,
		consumer: consumer,
	}
}

// InitializeQueues initializes the queues for the gRPC server if they don't exist
func initializeQueues(c *rabbitmq.Consumer, log *zap.Logger) {
	if c == nil {
		panic("subscriber is not initialized")
	}

	_, err := c.MakeNewQueueAndExchange(rabbitmq.BlockQueue, rabbitmq.BlockExchange, rabbitmq.BlockRoute)
	if err != nil {
		log.Fatal("failed to make block queue and exchange", zap.Error(err))
	}

	_, err = c.MakeNewQueueAndExchange(rabbitmq.TransactionQueue, rabbitmq.TransactionExchange, rabbitmq.TransactionRoute)
	if err != nil {
		log.Fatal("failed to make transaction queue and exchange", zap.Error(err))
	}

	_, err = c.MakeNewQueueAndExchange(rabbitmq.TransactionLogQueue, rabbitmq.TransactionLogExchange, rabbitmq.TransactionLogRoute)
	if err != nil {
		log.Fatal("failed to make transaction log queue and exchange", zap.Error(err))
	}

	_, err = c.MakeNewQueueAndExchange(rabbitmq.TokenEventQueue, rabbitmq.TokenEventExchange, rabbitmq.TokenEventRoute)
	if err != nil {
		log.Fatal("failed to make token event queue and exchange", zap.Error(err))
	}

	_, err = c.MakeNewQueueAndExchange(rabbitmq.RewardQueue, rabbitmq.RewardExchange, rabbitmq.RewardRoute)
	if err != nil {
		log.Fatal("failed to make reward queue and exchange", zap.Error(err))
	}

	_, err = c.MakeNewQueueAndExchange(rabbitmq.InternalTransactionQueue, rabbitmq.InternalTransactionExchange, rabbitmq.InternalTransactionRoute)
	if err != nil {
		log.Fatal("failed to make internal transaction queue and exchange", zap.Error(err))
	}

	_, err = c.MakeNewQueueAndExchange(rabbitmq.TransactionActionQueue, rabbitmq.TransactionActionExchange, rabbitmq.TransactionActionRoute)
	if err != nil {
		log.Fatal("failed to make transaction action queue and exchange", zap.Error(err))
	}

	_, err = c.MakeNewQueueAndExchange(rabbitmq.WithdrawalQueue, rabbitmq.WithdrawalExchange, rabbitmq.WithdrawalRoute)
	if err != nil {
		log.Fatal("failed to make withdrawal queue and exchange", zap.Error(err))
	}

	log.Info("queues and exchanges initialized")
}

// Start the gRPC server
func (s *Server) Start() {
	address := fmt.Sprintf(":%s", s.cfg.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.log.Fatal("failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCoreServiceServer(grpcServer, s.handler)
	reflection.Register(grpcServer)

	s.log.Info("gRPC Server started on " + address)
	if err := grpcServer.Serve(listener); err != nil {
		s.log.Fatal("failed to serve", zap.Error(err))
	}

	// Handle OS shutdown signals
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan
		s.log.Info("Shutting down gRPC server...")
		grpcServer.GracefulStop()
	}()

	s.log.Info("gRPC Server started on " + address)
	if err := grpcServer.Serve(listener); err != nil {
		s.log.Fatal("failed to serve", zap.Error(err))
	}
}
