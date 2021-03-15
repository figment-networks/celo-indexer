package server

import (
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/metrics"
	"github.com/figment-networks/celo-indexer/usecase"
	"github.com/figment-networks/celo-indexer/utils/logger"
	"github.com/gin-gonic/gin"
)

// Server handles HTTP requests
type Server struct {
	cfg      *config.Config
	handlers *usecase.HttpHandlers

	engine *gin.Engine
}

// New returns a new server instance
func New(cfg *config.Config, handlers *usecase.HttpHandlers) *Server {
	app := &Server{
		cfg:      cfg,
		engine:   gin.Default(),
		handlers: handlers,
	}
	return app.init()
}

// Start starts the server
func (s *Server) Start(listenAdd string) error {
	logger.Info("starting server...", logger.Field("app", "server"))

	go s.startMetricsServer()

	return s.engine.Run(listenAdd)
}

// init initializes the server
func (s *Server) init() *Server {
	logger.Info("initializing server...", logger.Field("app", "server"))

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) startMetricsServer() error {
	return metrics.NewMetricsServer().StartServer(s.cfg.ServerMetricAddr, s.cfg.MetricServerUrl)
}
