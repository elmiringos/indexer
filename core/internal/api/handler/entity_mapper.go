package handler

import (
	pb "github.com/elmiringos/indexer/indexer-core/internal/api/pb"
	"github.com/elmiringos/indexer/indexer-core/internal/domain/block"
)

func MapBlockToCurrentBlockResponse(block *block.Block) *pb.GetCurrentBlockResponse {
	if block == nil {
		return &pb.GetCurrentBlockResponse{}
	}

	return &pb.GetCurrentBlockResponse{
		BlockNumber: block.Number.Bytes(),
		BlockHash:   block.Hash.Hex(),
	}
}
