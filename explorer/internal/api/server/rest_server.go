package server

import (
	"net/http"

	"go.uber.org/zap"

	resthandler "github.com/elmiringos/indexer/explorer/internal/api/handler/rest"
	"github.com/elmiringos/indexer/explorer/internal/api/service"
)

type HTTPServer struct {
	router http.Handler
}

func NewRESTServer(
	blockService *service.BlockService,
	log *zap.Logger,
) *HTTPServer {
	router := resthandler.NewRouter(blockService, log)

	return &HTTPServer{
		router: router,
	}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
