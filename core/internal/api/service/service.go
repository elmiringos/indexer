package service

import (
	"context"
	"errors"

	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/internal_transaction"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/reward"
	smartcontract "github.com/elmiringos/indexer/indexer-core/internal/domain/smart_contract"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/token"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/transaction"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/withdrawal"
	"github.com/elmiringos/indexer/indexer-core/internal/infrastructure/repository"

	"go.uber.org/zap"
)

type CoreService struct {
	logger                        *zap.Logger
	BlockRepository               block.Repository
	InternalTransactionRepository internal_transaction.Repository
	RewardRepository              reward.Repository
	SmartContractRepository       smartcontract.Repository
	TokenRepository               token.Repository
	TransactionRepository         transaction.Repository
	WithdrawalRepository          withdrawal.Repository
}

var (
	ErrUnexpectedError = errors.New("unexpected error")
)

func NewCoreService(
	logger *zap.Logger,
	blockRepository block.Repository,
	internalTransactionRepository internal_transaction.Repository,
	rewardRepository reward.Repository,
	smartContractRepository smartcontract.Repository,
	tokenRepository token.Repository,
	transactionRepository transaction.Repository,
	withdrawalRepository withdrawal.Repository,
) *CoreService {
	return &CoreService{
		logger:                        logger,
		BlockRepository:               blockRepository,
		InternalTransactionRepository: internalTransactionRepository,
		RewardRepository:              rewardRepository,
		SmartContractRepository:       smartContractRepository,
		TokenRepository:               tokenRepository,
		TransactionRepository:         transactionRepository,
		WithdrawalRepository:          withdrawalRepository,
	}
}

func (s *CoreService) GetCurrentBlock(ctx context.Context) (*block.Block, error) {
	s.logger.Info("Getting current block")

	currentBlock, err := s.BlockRepository.GetCurrentBlock(ctx)
	if err != nil {
		switch err {
		case repository.ErrNotFound:
			return nil, nil
		default:
			s.logger.Error("error in getting current block", zap.Error(err))
			return nil, ErrUnexpectedError
		}
	}

	return currentBlock, nil
}
