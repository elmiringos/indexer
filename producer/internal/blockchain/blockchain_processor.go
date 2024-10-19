package blockchain

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/elmiringos/indexer/producer/pkg/logger"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go.uber.org/zap"
)

type BlockchainProcessor struct {
	ethHttpClient *ethclient.Client
	rawHttpClient *rpc.Client
	ethWSClient   *ethclient.Client
	rawWSClient   *rpc.Client
	log           *zap.Logger
}

// basicAuth creates a base64-encoded string for Basic Authentication header
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

// NewBlockchainProcessor initializes BlockchainProcessor with HTTP and WebSocket clients
func NewBlockchainProcessor(cfg *config.Config) *BlockchainProcessor {
	// Create both HTTP and WebSocket clients
	httpClient, err := createRPCClient(cfg.EthNode.HttpURL, cfg.EthNode.ApiKey)
	if err != nil {
		panic(fmt.Errorf("failed to create HTTP client: %v", err))
	}

	wsClient, err := createRPCClient(cfg.EthNode.WsURL, cfg.EthNode.ApiKey)
	if err != nil {
		panic(fmt.Errorf("failed to create WebSocket client: %v", err))
	}

	// Initialize the BlockchainProcessor struct
	blockchainProcessor := &BlockchainProcessor{
		rawHttpClient: httpClient,
		ethHttpClient: ethclient.NewClient(httpClient),
		rawWSClient:   wsClient,
		ethWSClient:   ethclient.NewClient(wsClient),
		log:           logger.GetLogger(),
	}

	return blockchainProcessor
}

// createRPCClient is a helper function to create an RPC client with authentication
func createRPCClient(url, apiKey string) (*rpc.Client, error) {
	client, err := rpc.DialOptions(
		context.Background(),
		url,
		rpc.WithHeader("Authorization", "Basic "+basicAuth("", apiKey)),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating RPC client for URL %s: %v", url, err)
	}
	return client, nil
}

// CloseClients closes both HTTP and WebSocket clients gracefully
func (p *BlockchainProcessor) CloseClients() {
	// Gracefully close HTTP and WebSocket clients
	if p.rawHttpClient != nil {
		p.rawHttpClient.Close()
	}

	if p.ethHttpClient != nil {
		p.ethHttpClient.Close()
	}

	if p.rawWSClient != nil {
		p.rawWSClient.Close()
	}

	if p.ethWSClient != nil {
		p.ethWSClient.Close()
	}
}

func (p *BlockchainProcessor) GetBlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error) {
	block, err := p.ethHttpClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("Error in aggragating block %v: %v", blockNumber, err)
	}
	return block, nil
}

func (p *BlockchainProcessor) ListenNewBlocks(startBlockNumber int, blocks chan<- *types.Block) {
	headers := make(chan *types.Header)
	sub, err := p.ethWSClient.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		p.log.Fatal("Failed to subscribe to new blocks", zap.Error(err))
	}

	startBlockNumberBigInt := big.NewInt(int64(startBlockNumber))

	for {
		select {
		case err := <-sub.Err():
			p.log.Fatal("Error with block subscription", zap.Error(err))
		case header := <-headers:
			if header.Number.Cmp(startBlockNumberBigInt) >= 0 {
				block, err := p.ethWSClient.BlockByHash(context.Background(), header.Hash())
				if err != nil {
					p.log.Fatal("Failed to get block by hash", zap.Error(err))
				}
				blocks <- block
				p.log.Debug("Successfuly get block and send to blocks channel", zap.Any("block", block))
			}
		}
	}
}
