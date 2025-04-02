package server

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/elmiringos/indexer/producer/internal/blockchain"
	grpccoreclient "github.com/elmiringos/indexer/producer/pkg/grpc_core_client"
	"github.com/elmiringos/indexer/producer/pkg/logger"
	"github.com/elmiringos/indexer/producer/pkg/rabbitmq"

	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/zap"
)

type Server struct {
	blockchainProcessor *blockchain.BlockchainProcessor
	grpcCoreClient      *grpccoreclient.CoreClient
	publisher           *rabbitmq.Publisher
	config              *config.Config
	log                 *zap.Logger
}

func NewServer(
	blockhainProcessor *blockchain.BlockchainProcessor,
	publisher *rabbitmq.Publisher,
	grpcCoreClient *grpccoreclient.CoreClient,
	cfg *config.Config,
) *Server {
	if blockhainProcessor == nil {
		panic("Blockhain processor is nil")
	}

	if publisher == nil {
		panic("Rabbitmq publisher is nil")
	}

	if grpcCoreClient == nil {
		panic("GRPC Core client is nil")
	}

	return &Server{
		blockchainProcessor: blockhainProcessor,
		publisher:           publisher,
		grpcCoreClient:      grpcCoreClient,
		config:              cfg,
		log:                 logger.GetLogger(),
	}
}

func (s *Server) setupAllQueues() {
	_, err := s.publisher.MakeNewQueueAndExchange(rabbitmq.BlockExchange, rabbitmq.BlockRoute, rabbitmq.BlockQueue)
	if err != nil {
		panic(fmt.Errorf("error in setuping block queue, %v", err))
	}

	_, err = s.publisher.MakeNewQueueAndExchange(rabbitmq.TransactionExchange, rabbitmq.TransactionRoute, rabbitmq.TransactionQueue)
	if err != nil {
		panic(fmt.Errorf("error in setuping transaction queue, %v", err))
	}

	_, err = s.publisher.MakeNewQueueAndExchange(rabbitmq.WithdrawalExchange, rabbitmq.WithdrawalRoute, rabbitmq.WithdrawalQueue)
	if err != nil {
		panic(fmt.Errorf("error in setuping withdrawal queue, %v", err))
	}

	_, err = s.publisher.MakeNewQueueAndExchange(rabbitmq.TransactionLogExchange, rabbitmq.TransactionLogRoute, rabbitmq.TransactionLogQueue)
	if err != nil {
		panic(fmt.Errorf("error in setuping transaction log queue, %v", err))
	}

	_, err = s.publisher.MakeNewQueueAndExchange(rabbitmq.InternalTransactionExchange, rabbitmq.InternalTransactionRoute, rabbitmq.InternalTransactionQueue)
	if err != nil {
		panic(fmt.Errorf("error in setuping internal transaction queue, %v", err))
	}

	_, err = s.publisher.MakeNewQueueAndExchange(rabbitmq.TransactionActionExchange, rabbitmq.TransactionActionRoute, rabbitmq.TransactionActionQueue)
	if err != nil {
		panic(fmt.Errorf("error in setuping transaction action queue, %v", err))
	}

	_, err = s.publisher.MakeNewQueueAndExchange(rabbitmq.TokenEventExchange, rabbitmq.TokenEventRoute, rabbitmq.TokenEventQueue)
	if err != nil {
		panic(fmt.Errorf("error in setuping token event queue, %v", err))
	}
}

func (s *Server) worker(id int, blocks <-chan *types.Block, wg *sync.WaitGroup) {
	defer wg.Done()

	channel := s.publisher.CreateChannel()
	s.log.Info("Creating new rabbitmqchannel", zap.Int("worker_id", id))
	defer channel.Close()

	for block := range blocks {
		s.log.Info("Worker started processing block", zap.Int("worker", id), zap.Int64("blockHeight", block.Number().Int64()))

		err := s.aggregateBlock(channel, block)
		if err != nil {
			s.log.Error("Error aggregating block", zap.Error(err))
		}

		s.log.Info("Worker finished processing block", zap.Int("worker", id), zap.Int64("blockHeight", block.Number().Int64()))
	}
}

func (s *Server) startWorkerPool(numWorkers int, blocks <-chan *types.Block, wg *sync.WaitGroup) {
	for id := 1; id <= numWorkers; id++ {
		wg.Add(1)
		go s.worker(id, blocks, wg)
	}
}

func (s *Server) ResetState() {
	messageState, err := s.grpcCoreClient.ResetState()
	if err != nil {
		s.log.Fatal("Error in reseting core service state", zap.Error(err))
	}

	if !messageState.Success {
		s.log.Fatal("Unsuccesful atemt in reseting core service state")
	}
}

func (s *Server) SyncStartingBlock(configBlockStartNumber *big.Int) {
	currentBlock, err := s.grpcCoreClient.GetCurrentBlock()
	if err != nil {
		s.log.Fatal("Error in getting starting block", zap.Error(err))
	}

	currentBlockNumber := big.NewInt(0).SetBytes(currentBlock.BlockNumber)

	if currentBlockNumber.Cmp(big.NewInt(0)) == 0 {
		s.log.Info("No starting block found, starting from block that placed in config.yml")
	} else if currentBlockNumber.Cmp(configBlockStartNumber) < 0 {
		s.log.Warn(
			"Current block number in Database is less than config block number. Some data will be lost.",
			zap.Any("currentBlockNumber", currentBlockNumber),
			zap.Any("startingBlockNumber", configBlockStartNumber),
		)
	} else if currentBlockNumber.Cmp(configBlockStartNumber) > 0 {
		s.log.Warn(
			"Current block number is greater than starting block number. This is not possible. Please check the database and the core service.",
			zap.Any("currentBlockNumber", currentBlockNumber),
			zap.Any("startingBlockNumber", configBlockStartNumber),
		)
	}
}

func (s *Server) StartBlockchainDataConsuming() {
	blockStartNumber, ok := big.NewInt(0).SetString(s.config.Server.BlockStartNumber, 10)
	if !ok {
		s.log.Fatal("Error in setting block start number", zap.String("blockStartNumber", s.config.Server.BlockStartNumber))
	}

	// Sync starting block before starting the workers
	s.SyncStartingBlock(blockStartNumber)

	// Setup all queues (blockQueue, transactionQueue, withdrawalQueue, transactionLogQueue, internalTransactionQueue, transactionActionQueue, tokenEventQueue)
	s.setupAllQueues()

	var wg sync.WaitGroup

	// Listen for new blocks
	blocks, err := s.blockchainProcessor.GenerateBlocks(context.Background(), blockStartNumber)
	if err != nil {
		s.log.Fatal("Error in generating blocks", zap.Error(err))
	}

	// Start the worker pool
	s.startWorkerPool(s.config.WorkerCount, blocks, &wg)
	wg.Wait()

	s.log.Info("Blocks channel is closed, workers have finished processing.")
}
