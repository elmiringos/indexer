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
	"sync"

	"github.com/elmiringos/indexer/producer/config"
	"github.com/elmiringos/indexer/producer/pkg/logger"

	"github.com/ethereum/go-ethereum"
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

const erc20ABI = `[
	{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},
	{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"}
]`

const erc721ABI = `[
	{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"tokenURI","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[{"name":"tokenId","type":"uint256"}],"name":"ownerOf","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},
	{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":true,"name":"tokenId","type":"uint256"}],"name":"Transfer","type":"event"}
]`

const erc1155ABI = `[
	{"constant":true,"inputs":[{"name":"id","type":"uint256"}],"name":"uri","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[{"name":"account","type":"address"},{"name":"id","type":"uint256"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},
	{"constant":true,"inputs":[{"name":"accounts","type":"address[]"},{"name":"ids","type":"uint256[]"}],"name":"balanceOfBatch","outputs":[{"name":"","type":"uint256[]"}],"payable":false,"stateMutability":"view","type":"function"},
	{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"id","type":"uint256"},{"indexed":false,"name":"value","type":"uint256"}],"name":"TransferSingle","type":"event"},
	{"anonymous":false,"inputs":[{"indexed":true,"name":"operator","type":"address"},{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"ids","type":"uint256[]"},{"indexed":false,"name":"values","type":"uint256[]"}],"name":"TransferBatch","type":"event"}
]`

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

func (p *BlockchainProcessor) GenerateHistoricalBlocks(ctx context.Context, configBlockNumber *big.Int, blocks chan<- *types.Block, latestBlock <-chan *types.Block) error {
	select {
	case latestBlockNumber := <-latestBlock:
		for configBlockNumber.Cmp(latestBlockNumber.Number()) == -1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				block, err := p.ethHttpClient.BlockByNumber(ctx, configBlockNumber)
				if err != nil {
					p.log.Error("Failed to get historical block",
						zap.Any("number", configBlockNumber),
						zap.Error(err),
					)
					// Exponential backoff could be added here
					configBlockNumber = configBlockNumber.Add(configBlockNumber, big.NewInt(1))
					continue
				}
				select {
				case blocks <- block:
					configBlockNumber = configBlockNumber.Add(configBlockNumber, big.NewInt(1))
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (p *BlockchainProcessor) ListenNewBlocks(ctx context.Context, blocks chan<- *types.Block, latestBlock chan<- *types.Block) error {
	headers := make(chan *types.Header)
	sentFirstBlock := false

	sub, err := p.ethWSClient.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new blocks: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-sub.Err():
			return fmt.Errorf("subscription error: %w", err)
		case header, ok := <-headers:
			if !ok {
				return fmt.Errorf("headers channel closed unexpectedly")
			}

			block, err := p.ethHttpClient.BlockByNumber(ctx, header.Number)
			if err != nil {
				p.log.Error("Failed to get block by hash",
					zap.Error(err),
					zap.String("hash", header.Hash().String()),
					zap.Any("number", header.Number),
				)
				continue
			}

			select {
			case blocks <- block:
				if !sentFirstBlock {
					fmt.Println("sending first block")
					select {
					case latestBlock <- block:
						sentFirstBlock = true
					case <-ctx.Done():
						return ctx.Err()
					}
				}
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// GenerateBlocks creates a stream of blocks starting from configBlockNumber,
// including both historical blocks and new incoming blocks.
// It returns a channel that will receive blocks in sequential order.
// The caller should provide a context for cancellation.
func (p *BlockchainProcessor) GenerateBlocks(ctx context.Context, configBlockNumber *big.Int) (<-chan *types.Block, error) {
	blocks := make(chan *types.Block, 100)
	latestBlock := make(chan *types.Block, 1)

	go func() {
		var wg sync.WaitGroup

		wg.Add(2)

		// Create error channels for goroutines
		historicalErr := make(chan error, 1)
		newBlocksErr := make(chan error, 1)

		// Start goroutines with error handling
		go func() {
			defer wg.Done()
			historicalErr <- p.GenerateHistoricalBlocks(ctx, configBlockNumber, blocks, latestBlock)
		}()

		go func() {
			defer wg.Done()
			newBlocksErr <- p.ListenNewBlocks(ctx, blocks, latestBlock)
		}()

		go func() {
			wg.Wait()
			close(blocks)
			close(latestBlock)
		}()

		// Wait for either context cancellation or an error
		select {
		case <-ctx.Done():
			p.log.Info("Context cancelled, stopping block generation")
		case err := <-historicalErr:
			if err != nil && err != context.Canceled {
				p.log.Error("Historical blocks error", zap.Error(err))
			}
		case err := <-newBlocksErr:
			if err != nil && err != context.Canceled {
				p.log.Error("New blocks subscription error", zap.Error(err))
			}
		}
	}()

	return blocks, nil
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

func (p *BlockchainProcessor) GetTokenEvents(receipt *types.Receipt, transactionHash common.Hash) []*TokenEvent {
	erc20, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		p.log.Fatal("error in reading ERC-20 token ABI", zap.Error(err))
	}

	erc721, err := abi.JSON(strings.NewReader(erc721ABI))
	if err != nil {
		p.log.Fatal("error in reading ERC-721 token ABI", zap.Error(err))
	}

	erc1155, err := abi.JSON(strings.NewReader(erc1155ABI))
	if err != nil {
		p.log.Fatal("error in reading ERC-1155 token ABI", zap.Error(err))
	}

	// Get event IDs for token transfers
	erc20TransferSig := erc20.Events["Transfer"].ID
	erc721TransferSig := erc721.Events["Transfer"].ID
	erc1155SingleSig := erc1155.Events["TransferSingle"].ID
	// erc1155BatchSig := erc1155.Events["TransferBatch"].ID

	var tokenEvents []*TokenEvent

	// Check if this transaction created a contract
	contractCreated := receipt.ContractAddress != (common.Address{})

	var contractBytecode string
	if contractCreated {
		var err error
		contractBytecode, err = p.getContractBytecode(receipt.ContractAddress)
		if err != nil {
			p.log.Error("failed to get contract bytecode", zap.Error(err))
		}
	}

	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}

		event := &TokenEvent{}

		switch log.Topics[0] {
		case erc20TransferSig:
			if len(log.Topics) >= 3 {
				from := common.BytesToAddress(log.Topics[1].Bytes())
				to := common.BytesToAddress(log.Topics[2].Bytes())
				value := FromBytesToBigInt(log.Data)

				// Get Token Metadata for ERC-20 token
				metadata := p.getERC20Metadata(log.Address)

				event.TransactionHash = transactionHash
				event.LogIndex = log.Index
				event.From = from
				event.To = to
				event.Value = value
				event.IsMint = from == zeroAddress
				event.IsBurn = to == zeroAddress
				event.TokenMetadata = metadata
			}

		case erc721TransferSig:
			if len(log.Topics) >= 4 {
				from := common.BytesToAddress(log.Topics[1].Bytes())
				to := common.BytesToAddress(log.Topics[2].Bytes())
				tokenId := new(big.Int).SetBytes(log.Topics[3].Bytes())

				// Get Token Metadata for ERC-721 token
				metadata := p.getERC721Metadata(log.Address, tokenId)

				event.TransactionHash = transactionHash
				event.LogIndex = log.Index
				event.TokenId = BigInt(*tokenId)
				event.From = from
				event.To = to
				event.IsMint = from == zeroAddress
				event.IsBurn = to == zeroAddress
				event.TokenMetadata = metadata
			}

		case erc1155SingleSig:
			if len(log.Topics) >= 4 && len(log.Data) >= 64 {
				from := common.BytesToAddress(log.Topics[2].Bytes())
				to := common.BytesToAddress(log.Topics[3].Bytes())
				tokenId := new(big.Int).SetBytes(log.Data[:32])
				value := new(big.Int).SetBytes(log.Data[32:64])

				// Get Token Metadata for ERC-1155 token
				metadata := p.getERC1155Metadata(log.Address, tokenId)

				event.TransactionHash = transactionHash
				event.LogIndex = log.Index
				event.TokenId = BigInt(*tokenId)
				event.From = from
				event.To = to
				event.Value = BigInt(*value)
				event.IsMint = from == zeroAddress
				event.IsBurn = to == zeroAddress
				event.TokenMetadata = metadata
			}
			// case erc1155BatchSig:
			// 	if len(log.Topics) >= 4 && len(log.Data) >= 64 {
			// 		from := common.BytesToAddress(log.Topics[2].Bytes())
			// 		to := common.BytesToAddress(log.Topics[3].Bytes())

			// 		// ERC-1155 Batch Transfers contain multiple token IDs & values
			// 		tokenIds, values := parseBatchTransferData(log.Data)

			// 		metadata := p.getERC1155Metadata(log.Address)

			// 		for i := range tokenIds {
			// 			batchEvent := &TokenEvent{
			// 				TransactionHash: transactionHash,
			// 				LogIndex:        log.Index,
			// 				TokenId:         tokenIds[i],
			// 				From:            from,
			// 				To:              to,
			// 				Value:           values[i],
			// 				IsMint:          from == zeroAddress,
			// 				IsBurn:          to == zeroAddress,
			// 				TokenMetadata:   metadata,
			// 			}
			// 			tokenEvents = append(tokenEvents, batchEvent)
			// 		}
			// 	}
		}

		// Fetch contract bytecode only if the transaction created a contract
		if contractCreated && log.Address == receipt.ContractAddress {
			if event.TokenMetadata == nil {
				event.TokenMetadata = make(map[string]interface{})
			}
			event.TokenMetadata["smartcontract_bytecode"] = contractBytecode
		}
		if event.From != (common.Address{}) || event.TokenId.String() != "0" || event.Value.String() != "0" {
			tokenEvents = append(tokenEvents, event)
		}
	}

	return tokenEvents
}

// getContractBytecode fetches the bytecode of a contract at a given address
func (p *BlockchainProcessor) getContractBytecode(address common.Address) (string, error) {
	// Call the eth_getCode method to get the bytecode from the Ethereum client
	result, err := p.ethHttpClient.CodeAt(context.Background(), address, nil)
	if err != nil {
		return "", err
	}

	// Return the bytecode as a hexadecimal string
	return common.Bytes2Hex(result), nil
}

func (p *BlockchainProcessor) getERC20Metadata(tokenAddress common.Address) TokenMetadata {
	metadata := TokenMetadata{}

	erc20, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		p.log.Error("error in reading ERC-20 ABI", zap.Error(err))
		return metadata
	}

	ctx := context.Background()

	// Helper function to safely call a contract method
	callContract := func(methodName string) ([]byte, error) {
		method, exists := erc20.Methods[methodName]
		if !exists {
			return nil, fmt.Errorf("method %s does not exist", methodName)
		}

		// Properly encode function call
		inputData, err := method.Inputs.Pack()
		if err != nil {
			return nil, fmt.Errorf("failed to encode %s: %w", methodName, err)
		}

		msg := ethereum.CallMsg{
			To:   &tokenAddress,
			Data: append(method.ID, inputData...),
		}

		return p.ethHttpClient.CallContract(ctx, msg, nil)
	}

	// Name
	var name string
	if nameData, err := callContract("name"); err == nil && len(nameData) > 0 {
		if err := erc20.UnpackIntoInterface(&name, "name", nameData); err == nil {
			metadata["name"] = name
		} else {
			p.log.Warn("Failed to decode ERC-20 name", zap.Error(err))
		}
	}

	// Symbol
	var symbol string
	if symbolData, err := callContract("symbol"); err == nil && len(symbolData) > 0 {
		if err := erc20.UnpackIntoInterface(&symbol, "symbol", symbolData); err == nil {
			metadata["symbol"] = symbol
		} else {
			p.log.Warn("Failed to decode ERC-20 symbol", zap.Error(err))
		}
	}

	// Decimals
	var decimals uint8
	if decimalsData, err := callContract("decimals"); err == nil && len(decimalsData) > 0 {
		if err := erc20.UnpackIntoInterface(&decimals, "decimals", decimalsData); err == nil {
			metadata["decimals"] = decimals
		} else {
			p.log.Warn("Failed to decode ERC-20 decimals", zap.Error(err))
		}
	}

	return metadata
}

func (p *BlockchainProcessor) getERC721Metadata(tokenAddress common.Address, tokenId *big.Int) TokenMetadata {
	metadata := make(TokenMetadata)

	erc721, err := abi.JSON(strings.NewReader(erc721ABI))
	if err != nil {
		p.log.Error("error in reading ERC-721 ABI", zap.Error(err))
		return metadata
	}

	ctx := context.Background()

	// Helper function to call contract methods safely
	callContract := func(methodName string, args ...interface{}) ([]byte, error) {
		method, exists := erc721.Methods[methodName]
		if !exists {
			return nil, fmt.Errorf("method %s does not exist", methodName)
		}

		data, err := method.Inputs.Pack(args...)
		if err != nil {
			return nil, err
		}

		msg := ethereum.CallMsg{
			To:   &tokenAddress,
			Data: append(method.ID, data...),
		}

		return p.ethHttpClient.CallContract(ctx, msg, nil)
	}

	// Fetch Name
	if nameData, err := callContract("name"); err == nil {
		var name string
		if err := erc721.UnpackIntoInterface(&name, "name", nameData); err == nil {
			metadata["name"] = name
		} else {
			p.log.Warn("Failed to decode ERC-721 name", zap.Error(err))
		}
	}

	// Fetch Symbol
	if symbolData, err := callContract("symbol"); err == nil {
		var symbol string
		if err := erc721.UnpackIntoInterface(&symbol, "symbol", symbolData); err == nil {
			metadata["symbol"] = symbol
		} else {
			p.log.Warn("Failed to decode ERC-721 symbol", zap.Error(err))
		}
	}

	// Fetch Token URI (Requires Token ID)
	if tokenURIData, err := callContract("tokenURI", tokenId); err == nil {
		var tokenURI string
		if err := erc721.UnpackIntoInterface(&tokenURI, "tokenURI", tokenURIData); err == nil {
			metadata["tokenURI"] = tokenURI
		} else {
			p.log.Warn("Failed to decode ERC-721 tokenURI", zap.Error(err))
		}
	}

	return metadata
}

func (p *BlockchainProcessor) getERC1155Metadata(tokenAddress common.Address, tokenId *big.Int) TokenMetadata {
	metadata := make(TokenMetadata)

	erc1155, err := abi.JSON(strings.NewReader(erc1155ABI))
	if err != nil {
		p.log.Error("error in reading ERC-1155 ABI", zap.Error(err))
		return metadata
	}

	ctx := context.Background()

	// Helper function to call contract methods safely
	callContract := func(methodName string, args ...interface{}) ([]byte, error) {
		method, exists := erc1155.Methods[methodName]
		if !exists {
			return nil, fmt.Errorf("method %s does not exist", methodName)
		}

		data, err := method.Inputs.Pack(args...)
		if err != nil {
			return nil, err
		}

		msg := ethereum.CallMsg{
			To:   &tokenAddress,
			Data: append(method.ID, data...),
		}

		return p.ethHttpClient.CallContract(ctx, msg, nil)
	}

	// Fetch Token URI
	if uriData, err := callContract("uri", tokenId); err == nil {
		var uri string
		if err := erc1155.UnpackIntoInterface(&uri, "uri", uriData); err == nil {
			metadata["uri"] = uri
		} else {
			p.log.Warn("Failed to decode ERC-1155 uri", zap.Error(err))
		}
	}

	return metadata
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

	txHash := trace["transactionHash"].([]byte)
	blockHash := trace["blockHash"].([]byte)
	from := action["from"].([]byte)
	to := action["to"].([]byte)

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
		inputData, _ = hex.DecodeString(input[2:])   // Delete "0x"
		selector = hex.EncodeToString(inputData[:4]) // Delete first 4 bytes
	}

	var outputData []byte
	if output, ok := result["output"].(string); ok && len(output) >= 10 {
		outputData, _ = hex.DecodeString(output[2:]) // Delete "0x"
	}

	var internalTx *InternalTransaction

	if traceType == "call" || traceType == "create" || traceType == "selfdestruct" {
		internalTx = &InternalTransaction{
			BlockHash:       common.BytesToHash(blockHash),
			Index:           traceIndex,
			Type:            traceType,
			TransactionHash: common.BytesToHash(txHash),
			Status:          status,
			GasUsed:         gasUsed.Uint64(),
			Gas:             gas.Uint64(),
			Input:           inputData,
			Output:          outputData,
			Value:           value,
			From:            common.BytesToAddress(from),
			To:              common.BytesToAddress(to),
			Timestamp:       block.Time(),
			ErrorMsg:        errorMsg,
		}

		if traceType == "create" && hasResult && result["address"] != nil {
			internalTx.ContractAddress = common.BytesToAddress(result["address"].([]byte))
		}
	}

	var txAction *TransactionAction

	if selector != "" || traceType == "log" {
		txAction = &TransactionAction{
			TransactionHash: common.BytesToHash(txHash),
			Selector:        selector,
			Type:            traceType,
			From:            common.BytesToAddress(from),
			To:              common.BytesToAddress(to),
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

func (p *BlockchainProcessor) GetTransactionSender(tx *types.Transaction) (common.Address, error) {
	// Use LatestSignerForChainID which handles all transaction types
	signer := types.LatestSignerForChainID(tx.ChainId())
	from, err := types.Sender(signer, tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to get transaction sender: %w", err)
	}

	return from, nil
}
