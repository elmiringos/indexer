package server

import (
	"fmt"
	"sync"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/elmiringos/indexer/producer/internal/blockchain"
	"github.com/elmiringos/indexer/producer/pkg/logger"
	"github.com/elmiringos/indexer/producer/pkg/rabbitmq"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Server struct {
	blockchainProcessor *blockchain.BlockchainProcessor
	publisher           *rabbitmq.Publisher
	config              *config.Config
	log                 *zap.Logger
	quques              map[rabbitmq.QuqueType]*amqp091.Queue
}

func NewServer(blockhainProcessor *blockchain.BlockchainProcessor, publisher *rabbitmq.Publisher, cfg *config.Config) *Server {
	if blockhainProcessor == nil {
		panic("Blockhain processor is nil")
	}

	if publisher == nil {
		panic("Rabbitmq publisher is nil")
	}

	return &Server{
		blockchainProcessor: blockhainProcessor,
		publisher:           publisher,
		config:              cfg,
		log:                 logger.GetLogger(),
	}
}

func (s *Server) aggregateBlock(block *types.Block) error {
	err := s.publisher.PublishBlockMessage(block)
	if err != nil {
		s.log.Error("Error in publishing block message to broker", zap.Any("message", block))
	}

	for _, transaction := range block.Transactions() {
		err = s.publisher.PublishTransactionMessage(transaction)
		if err != nil {
			s.log.Error("Error in puplishing transaction message to broker")
		}
	}

	return nil
}

func (s *Server) setupAllQueues() {
	blockQueue, err := s.publisher.MakeNewQueueAndExchange(rabbitmq.BlockExchange, rabbitmq.BlockRoute, rabbitmq.BlockQuque)
	if err != nil {
		panic(fmt.Errorf("Error in setuping block queue, %v", err))
	}

	transactionQueue, err := s.publisher.MakeNewQueueAndExchange(rabbitmq.TransactionExchange, rabbitmq.TransactionRoute, rabbitmq.TransactionQuque)
	if err != nil {
		panic(fmt.Errorf("Error in setuping transaction queue, %v", err))
	}

	s.quques = make(map[rabbitmq.QuqueType]*amqp091.Queue)
	s.quques[rabbitmq.BlockQuque] = blockQueue
	s.quques[rabbitmq.TransactionQuque] = transactionQueue
}

func (s *Server) worker(id int, blocks <-chan *types.Block, wg *sync.WaitGroup) {
	defer wg.Done()

	for block := range blocks {
		s.log.Info("Worker started processing block", zap.Int("worker", id), zap.Int64("blockHeight", block.Number().Int64()))

		err := s.aggregateBlock(block)
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

func (s *Server) StartBlockchainDataConsuming() {
	s.setupAllQueues()

	var wg sync.WaitGroup

	blocks := s.blockchainProcessor.ListenNewBlocks(0)

	s.startWorkerPool(s.config.WorkerCount, blocks, &wg)
	wg.Wait()

	s.log.Info("Blocks channel is closed, workers have finished processing.")
}
