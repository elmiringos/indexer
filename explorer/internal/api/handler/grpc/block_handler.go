package grpc

import (
	"context"

	"github.com/elmiringos/indexer/explorer/internal/api/pb"
	"github.com/elmiringos/indexer/explorer/internal/api/service"
	"go.uber.org/zap"
)

// BlockHandler implements the gRPC service for blocks
type BlockHandler struct {
	pb.UnimplementedExplorerServiceServer
	BlockService *service.BlockService
	log          zap.Logger
}

// NewBlockHandler returns a new instance of BlockHandler
func NewBlockHandler(s *service.BlockService, log *zap.Logger) *BlockHandler {
	return &BlockHandler{
		BlockService: s,
		log:          *log,
	}
}

// GetBlock handles the gRPC request to fetch a block by hash
func (h *BlockHandler) GetCurrentBlock(ctx context.Context, req *pb.GetCurrentBlockRequest) (*pb.GetCurrentBlockResponse, error) {
	block, err := h.BlockService.GetCurrentBlock()
	if err != nil {
		return nil, err
	}

	return &pb.GetCurrentBlockResponse{
		Block: &pb.Block{
			Hash:          block.Hash.String(),
			Number:        block.Number.String(),
			ParentHash:    block.ParentHash.String(),
			MinerHash:     block.MinerHash.String(),
			GasLimit:      block.GasLimit,
			GasUsed:       block.GasUsed,
			Nonce:         block.Nonce,
			Size:          block.Size,
			Difficulty:    block.Difficulty.String(),
			IsPos:         block.IsPos,
			BaseFeePerGas: block.BaseFeePerGas.String(),
			Timestamp:     block.Timestamp,
		},
	}, nil
}
