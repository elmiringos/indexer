package blockchain

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/elmiringos/indexer/producer/pkg/logger"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
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

const erc20ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}]`

const erc721ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"}]`

const erc1155ABI = `[{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferSingle","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"ids","type":"uint256[]"},{"indexed":false,"name":"values","type":"uint256[]"}],"name":"TransferBatch","type":"event"}]`

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

// createRPCClient creates an RPC client with authentication
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

// GetBlockByNumber gets a block by number
func (p *BlockchainProcessor) GetBlockByNumber(ctx context.Context, blockNumber *big.Int) (*types.Block, error) {
	block, err := p.ethHttpClient.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("error in aggragating block %v: %v", blockNumber, err)
	}
	return block, nil
}

// ListenNewBlocks listens for new blocks and get old blocks
func (p *BlockchainProcessor) ListenNewBlocks(startBlockNumber int64) <-chan *types.Block {
	headers := make(chan *types.Header)
	blocks := make(chan *types.Block, 100)

	go func() {
		defer close(blocks)

		latestBlock, err := p.ethHttpClient.BlockByNumber(context.Background(), nil)
		if err != nil {
			p.log.Error("Failed to get latest block", zap.Error(err))
			return
		}
		currentBlockNum := startBlockNumber
		latestBlockNum := latestBlock.Number().Int64()

		for currentBlockNum <= latestBlockNum {
			block, err := p.ethHttpClient.BlockByNumber(
				context.Background(),
				big.NewInt(currentBlockNum),
			)
			if err != nil {
				p.log.Error("Failed to get historical block",
					zap.Int64("number", currentBlockNum),
					zap.Error(err),
				)
				currentBlockNum++
				continue
			}
			blocks <- block
			currentBlockNum++
		}

		sub, err := p.ethWSClient.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			p.log.Error("Failed to subscribe to new blocks", zap.Error(err))
			return
		}
		defer sub.Unsubscribe()

		for {
			select {
			case err := <-sub.Err():
				p.log.Error("Error in block subscription", zap.Error(err))
				return
			case header, ok := <-headers:
				if !ok {
					return
				}
				block, err := p.ethWSClient.BlockByHash(context.Background(), header.Hash())
				if err != nil {
					p.log.Error("Failed to get block by hash", zap.Error(err))
					continue
				}
				p.log.Debug("Get block by hash", zap.Any("number", block.Number()))
				blocks <- block
				p.log.Debug("Successful send block to blocks channel", zap.Any("number", block.Number()))
			}
		}
	}()

	return blocks
}

// GetBlockTraces returns the traces of a block
func (p *BlockchainProcessor) GetBlockTraces(blockNumber *big.Int) ([]map[string]interface{}, error) {
	var blockTrace []map[string]interface{}
	err := p.rawHttpClient.Call(&blockTrace, "trace_block", blockNumber.Uint64())
	if err != nil {
		return nil, fmt.Errorf("error in fetching block traces: %v", err)
	}

	return blockTrace, nil
}

// GetTokenEvents returns the token events of a transaction
func (p *BlockchainProcessor) GetTokenEvents(receipt *types.Receipt) []*TokenEvent {
	erc20, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		p.log.Fatal("error in reading erc20 token", zap.Error(err))
	}

	erc721, err := abi.JSON(strings.NewReader(erc721ABI))
	if err != nil {
		p.log.Fatal("error in reading erc721 token", zap.Error(err))
	}

	erc1155, err := abi.JSON(strings.NewReader(erc1155ABI))
	if err != nil {
		p.log.Fatal("error in reading erc1155 token", zap.Error(err))
	}

	// Get event IDs
	erc20TransferSig := erc20.Events["Transfer"].ID
	erc721TransferSig := erc721.Events["Transfer"].ID
	erc1155SingleSig := erc1155.Events["TransferSingle"].ID

	var tokenEvents []*TokenEvent

	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}

		switch log.Topics[0] {
		case erc20TransferSig:
			if len(log.Topics) >= 3 {
				from := strings.ToLower(log.Topics[1].Hex())
				to := strings.ToLower(log.Topics[2].Hex())
				value := new(big.Int).SetBytes(log.Data)

				event := &TokenEvent{
					From:   from,
					To:     to,
					Value:  value,
					IsMint: from == zeroAddress,
					IsBurn: to == zeroAddress,
				}
				tokenEvents = append(tokenEvents, event)
			}

		case erc721TransferSig:
			if len(log.Topics) >= 4 {
				from := strings.ToLower(log.Topics[1].Hex())
				to := strings.ToLower(log.Topics[2].Hex())
				tokenId := new(big.Int).SetBytes(log.Topics[3].Bytes())

				event := &TokenEvent{
					From:    from,
					To:      to,
					TokenId: tokenId,
					IsMint:  from == zeroAddress,
					IsBurn:  to == zeroAddress,
				}
				tokenEvents = append(tokenEvents, event)
			}

		case erc1155SingleSig:
			if len(log.Topics) >= 4 {
				from := strings.ToLower(log.Topics[2].Hex())
				to := strings.ToLower(log.Topics[3].Hex())
				event := &TokenEvent{
					From:    from,
					To:      to,
					TokenId: new(big.Int).SetBytes(log.Data[:32]),
					Value:   new(big.Int).SetBytes(log.Data[32:]),
					IsMint:  from == zeroAddress,
					IsBurn:  to == zeroAddress,
				}
				tokenEvents = append(tokenEvents, event)
			}
		}
	}

	return tokenEvents
}

// GetTransaction returns a transaction by hash
func (p *BlockchainProcessor) GetTransaction(txHash string) (*types.Transaction, error) {
	tx, _, err := p.ethHttpClient.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionReceipt returns the receipt of a transaction
func (p *BlockchainProcessor) GetTransactionReceipt(tx *types.Transaction) (*types.Receipt, error) {
	receipt, err := p.ethHttpClient.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

// GetTransactionTrace returns the trace of a transaction (only for full node)
func (p *BlockchainProcessor) GetTransactionTrace(tx *types.Transaction) ([]map[string]interface{}, error) {
	var trace []map[string]interface{}
	err := p.rawHttpClient.Call(&trace, "trace_transaction", tx.Hash().String())
	if err != nil {
		return nil, err
	}

	return trace, nil
}

// ProcessTransactionTrace processes a transaction trace into internal transaction and transaction actions
func (p *BlockchainProcessor) ProcessTransactionTrace(traceIndex int, trace map[string]interface{}, block *types.Block) (*InternalTransaction, *TransactionAction, error) {
	traceType := trace["type"].(string)
	action := trace["action"].(map[string]interface{})
	result, hasResult := trace["result"].(map[string]interface{})
	errorMsg, hasError := trace["error"].(string)

	txHash := trace["transactionHash"].(string)
	blockHash := trace["blockHash"].(string)
	from := action["from"].(string)
	to := action["to"].(string)

	gas, _ := new(big.Int).SetString(action["gas"].(string)[2:], 16)
	gasUsed := new(big.Int)
	if hasResult {
		gasUsed, _ = new(big.Int).SetString(result["gasUsed"].(string)[2:], 16)
	}
	value, _ := new(big.Int).SetString(action["value"].(string)[2:], 16)

	status := 1
	if hasError || !hasResult || gasUsed.Cmp(gas) >= 0 {
		status = 0
	}

	var inputData []byte
	var selector string
	if input, ok := action["input"].(string); ok && len(input) >= 10 {
		inputData, _ = hex.DecodeString(input[2:])   // Убираем "0x"
		selector = hex.EncodeToString(inputData[:4]) // Первые 4 байта
	}

	var outputData []byte
	if output, ok := result["output"].(string); ok && len(output) >= 10 {
		outputData, _ = hex.DecodeString(output[2:]) // Убираем "0x"
	}

	var internalTx *InternalTransaction

	if traceType == "call" || traceType == "create" || traceType == "selfdestruct" {
		internalTx = &InternalTransaction{
			BlockHash:       blockHash,
			Index:           traceIndex,
			Type:            traceType,
			TransactionHash: txHash,
			Status:          status,
			GasUsed:         gasUsed.Uint64(),
			Gas:             gas.Uint64(),
			Input:           inputData,
			Output:          outputData,
			Value:           value,
			From:            from,
			To:              to,
			Timestamp:       block.Time(),
			ErrorMsg:        errorMsg,
		}

		if traceType == "create" && hasResult && result["address"] != nil {
			internalTx.ContractAddress = result["address"].(string)
		}
	}

	var txAction *TransactionAction

	if selector != "" || traceType == "log" {
		txAction = &TransactionAction{
			TransactionHash: txHash,
			Selector:        selector,
			Type:            traceType,
			From:            from,
			To:              to,
			Value:           value,
			Input:           inputData,
			Status:          status,
		}

		if traceType == "log" {
			topics := trace["topics"].([]string)
			if len(topics) > 0 {
				txAction.Selector = topics[0] // Event signature
			}
		}
	}

	return internalTx, txAction, nil
}

// FetchMetadata retrieves and parses metadata from a token URI
func FetchMetadata(tokenURI string) (*TokenMetadata, error) {
	resp, err := http.Get(tokenURI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var metadata TokenMetadata
	json.NewDecoder(resp.Body).Decode(&metadata)
	return &metadata, nil
}
