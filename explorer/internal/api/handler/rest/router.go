package rest

import (
	"net/http"

	"github.com/elmiringos/indexer/explorer/internal/api/service"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(
	BlockService *service.BlockService,
	logger *zap.Logger,
) *mux.Router {
	r := mux.NewRouter()

	blockHandler := NewBlockHandler(BlockService, logger)

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/block/current", blockHandler.GetCurrentBlock).Methods(http.MethodGet)

	return r
}
