//go:generate protoc --go_out=. --go-grpc_out=. ../proto/core.proto

package handler

import (
	"context"

	"github.com/elmiringos/indexer/indexer-core/internal/api/pb"
	"github.com/elmiringos/indexer/indexer-core/internal/api/service"

	"go.uber.org/zap"
)

type CoreHandler struct {
	coreService *service.CoreService
	logger      *zap.Logger
	pb.UnimplementedCoreServiceServer
}

func NewCoreHandler(coreService *service.CoreService, logger *zap.Logger) *CoreHandler {
	return &CoreHandler{
		coreService: coreService,
		logger:      logger,
	}
}

func (h *CoreHandler) GetCurrentBlock(ctx context.Context, req *pb.GetCurrentBlockRequest) (*pb.GetCurrentBlockResponse, error) {
	block, err := h.coreService.GetCurrentBlock(ctx)
	if err != nil {
		return nil, err
	}

	return MapBlockToCurrentBlockResponse(block), nil
}

func (h *CoreHandler) ResetState(ctx context.Context, req *pb.ResetStateRequest) (*pb.ResetStateResponse, error) {
	state := true
	return &pb.ResetStateResponse{
		Success: state,
	}, nil
}
