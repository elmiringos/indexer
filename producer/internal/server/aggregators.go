package server

import (
	"encoding/json"

	"github.com/elmiringos/indexer/producer/internal/blockchain"
	"github.com/elmiringos/indexer/producer/pkg/rabbitmq"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Trace struct {
	Type        string                 `json:"type"`
	Action      map[string]interface{} `json:"action"`
	Result      map[string]interface{} `json:"result"`
	Transaction string                 `json:"transactionHash"`
	Block       string                 `json:"blockHash"`
	TraceAddr   []int                  `json:"traceAddress"`
}

// aggregateBlock aggregates a block and publishes the messages to the broker
func (s *Server) aggregateBlock(channel *amqp.Channel, block *types.Block) error {
	// publish block message
	blockMessage := blockchain.ConvertBlockToBlock(block)

	s.log.Debug("publishing block message to broker", zap.Any("block", blockMessage))

	err := s.publisher.PublishMessage(channel, rabbitmq.BlockExchange, rabbitmq.BlockRoute, blockMessage)
	if err != nil {
		s.log.Error("error in publishing block message to broker", zap.Any("block", block), zap.Error(err))
		return err
	}

	s.log.Debug("starting to aggregate transactions", zap.Bool("channel", channel == nil))

	// aggregate transactions and compute reward
	totalGasFees, err := s.aggregateTransactions(channel, block)
	if err != nil {
		s.log.Error("error in aggregating transactions", zap.Error(err))
		return err
	}

	// aggregate reward
	err = s.aggregateReward(channel, totalGasFees, block)
	if err != nil {
		s.log.Error("error in aggregating reward", zap.Error(err))
		return err
	}

	// aggregate withdrawals
	err = s.aggregateWithdrawals(channel, block.Withdrawals(), block.Hash())
	if err != nil {
		s.log.Error("error in aggregating withdrawals", zap.Error(err))
		return err
	}

	return nil
}

func (s *Server) aggregateTransactions(channel *amqp.Channel, block *types.Block) (uint64, error) {
	var totalGasFees uint64

	for index, transaction := range block.Transactions() {
		// get transaction receipt
		transactionReceipt, err := s.blockchainProcessor.GetTransactionReceipt(transaction)
		if err != nil {
			s.log.Error("error in getting transaction receipt")
			return 0, err
		}

		// publish transaction message
		transactionMessage, err := s.blockchainProcessor.ConvertTransactionToTransaction(transaction, block.Hash(), transactionReceipt, index)
		if err != nil {
			s.log.Error("error in converting transaction to custom type", zap.Error(err))
			return 0, err
		}

		s.log.Debug(
			"publishing transaction message to broker",
			zap.String("transaction hash", transactionMessage.Hash.String()),
			zap.String("block hash", transactionMessage.BlockHash.String()),
		)

		err = s.publisher.PublishMessage(channel, rabbitmq.TransactionExchange, rabbitmq.TransactionRoute, transactionMessage)
		if err != nil {
			s.log.Error("error in publishing transaction message to broker", zap.Error(err))
			return 0, err
		}

		// compute reward
		totalGasFees += transactionReceipt.CumulativeGasUsed

		// aggregate transaction logs
		err = s.aggragateTransactionLogs(channel, transactionReceipt.Logs)
		if err != nil {
			s.log.Error("error in aggregating transaction logs", zap.Error(err))
			return 0, err
		}

		// aggregate transaction trace if full node used
		if s.config.EthNode.Trace {
			transactionTrace, err := s.blockchainProcessor.GetTransactionTrace(transaction)
			if err != nil {
				s.log.Error("error in getting transaction trace", zap.Error(err), zap.String("transactionHash", transaction.Hash().String()))
				return 0, err
			}

			err = s.processTransactionTrace(channel, block, transactionTrace)
			if err != nil {
				s.log.Error("error in aggregating transaction trace", zap.Error(err), zap.String("transactionHash", transaction.Hash().String()))
				return 0, err
			}
		}

		// aggregate token events
		tokenEvents := s.blockchainProcessor.GetTokenEvents(transactionReceipt)
		err = s.aggregateTokenEvents(channel, tokenEvents)
		if err != nil {
			s.log.Error("error in aggregating token events", zap.Error(err))
			return 0, err
		}
	}

	return totalGasFees, nil
}

func (s *Server) aggregateWithdrawals(channel *amqp.Channel, withdrawals []*types.Withdrawal, blockHash common.Hash) error {
	for _, withdrawal := range withdrawals {
		withdrawalMessage := blockchain.ConvertWithdrawalToWithdrawal(withdrawal, blockHash)

		err := s.publisher.PublishMessage(channel, rabbitmq.WithdrawalExchange, rabbitmq.WithdrawalRoute, withdrawalMessage)
		if err != nil {
			s.log.Error("error in publishing withdrawal message to broker", zap.Error(err))
			return err
		}

		s.log.Debug("published withdrawal message to broker", zap.Any("withdrawal", withdrawalMessage))
	}

	return nil
}

func (s *Server) aggregateTokenEvents(channel *amqp.Channel, tokenEvents []*blockchain.TokenEvent) error {
	for _, tokenEvent := range tokenEvents {
		tokenEventMessage, err := json.Marshal(tokenEvent)
		if err != nil {
			s.log.Error("error in marshalling token event message", zap.Error(err))
			return err
		}

		err = s.publisher.PublishMessage(channel, rabbitmq.TokenEventExchange, rabbitmq.TokenEventRoute, tokenEventMessage)
		if err != nil {
			s.log.Error("error in publishing token event message to broker", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *Server) aggragateTransactionLogs(channel *amqp.Channel, transactionLogs []*types.Log) error {
	for _, transactionLog := range transactionLogs {
		transactionLogMessage, err := json.Marshal(transactionLog)
		if err != nil {
			s.log.Error("error in marshalling transaction log message", zap.Error(err))
			return err
		}

		err = s.publisher.PublishMessage(channel, rabbitmq.TransactionLogExchange, rabbitmq.TransactionLogRoute, transactionLogMessage)
		if err != nil {
			s.log.Error("error in publishing transaction log message to broker", zap.Error(err))
			return err
		}
	}

	return nil
}

func (s *Server) aggregateReward(channel *amqp.Channel, totalGasFees uint64, block *types.Block) error {
	reward := &blockchain.Reward{
		Address:   block.Coinbase(),
		Amount:    totalGasFees,
		BlockHash: block.Hash(),
	}

	err := s.publisher.PublishMessage(channel, rabbitmq.RewardExchange, rabbitmq.RewardRoute, reward)
	if err != nil {
		s.log.Error("error in publishing reward message to broker", zap.Error(err))
		return err
	}

	s.log.Debug("published reward message to broker", zap.Any("reward", reward))

	return nil
}

func (s *Server) processTransactionTrace(channel *amqp.Channel, block *types.Block, traces []map[string]interface{}) error {
	for index, trace := range traces {
		internalTransaction, transactionAction, err := s.blockchainProcessor.ProcessTransactionTrace(index, trace, block)
		if err != nil {
			s.log.Error("error in processing transaction trace", zap.Error(err))
			return err
		}

		// publish internal transaction message
		internalTransactionMessage, err := json.Marshal(internalTransaction)
		if err != nil {
			s.log.Error("error in marshalling internal transaction message", zap.Error(err))
			return err
		}

		err = s.publisher.PublishMessage(channel, rabbitmq.InternalTransactionExchange, rabbitmq.InternalTransactionRoute, internalTransactionMessage)
		if err != nil {
			s.log.Error("error in publishing internal transaction message to broker", zap.Error(err))
			return err
		}

		// publish transaction action message
		transactionActionMessage, err := json.Marshal(transactionAction)
		if err != nil {
			s.log.Error("error in marshalling transaction action message", zap.Error(err))
			return err
		}

		err = s.publisher.PublishMessage(channel, rabbitmq.TransactionActionExchange, rabbitmq.TransactionActionRoute, transactionActionMessage)
		if err != nil {
			s.log.Error("error in publishing transaction action message to broker", zap.Error(err))
			return err
		}
	}

	return nil
}
