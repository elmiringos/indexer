package rest

import (
	"encoding/json"
	"net/http"

	"github.com/elmiringos/indexer/explorer/internal/api/service"
	"go.uber.org/zap"
)

type BlockHandler struct {
	blockService *service.BlockService
	log          *zap.Logger
}

func NewBlockHandler(blockService *service.BlockService, log *zap.Logger) *BlockHandler {
	return &BlockHandler{
		blockService: blockService,
		log:          log,
	}
}

func (h *BlockHandler) GetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block, err := h.blockService.GetCurrentBlock()
	if err != nil {
		h.log.Error("Failed to get current block", zap.Error(err))
		http.Error(w, "Failed to get current block", http.StatusInternalServerError)
		return
	}

	response := MapBlockToCurrentBlockResponse(block)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("Failed to encode response", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
