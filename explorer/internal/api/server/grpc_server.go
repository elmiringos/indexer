package server

import (
	grpchandlers "github.com/elmiringos/indexer/explorer/internal/api/handler/grpc"
	"github.com/elmiringos/indexer/explorer/internal/api/pb"
	"github.com/elmiringos/indexer/explorer/internal/api/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(blockService *service.BlockService, log *zap.Logger) *grpc.Server {
	s := grpc.NewServer()

	// Inititalize handlers
	blockHandler := grpchandlers.NewBlockHandler(blockService, log)

	pb.RegisterExplorerServiceServer(s, blockHandler)

	reflection.Register(s)

	return s
}
