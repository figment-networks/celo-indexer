package metrics

import (
	"fmt"
	"net/http"

	"github.com/figment-networks/indexing-engine/metrics"
	"github.com/figment-networks/indexing-engine/metrics/prometheusmetrics"

	"github.com/figment-networks/celo-indexer/utils/logger"
)

// MetricsServer handles HTTP requests
type MetricsServer struct{}

// NewMetricsServer returns a new server instance
func NewMetricsServer() *MetricsServer {
	logger.Info("initializing metrics server...", logger.Field("app", "server"))
	return &MetricsServer{}
}

// StartServer starts the metrics server
func (ms *MetricsServer) StartServer(listenAddr string, url string) error {
	logger.Info(fmt.Sprintf("starting metrics server at %s...", url), logger.Field("app", "indexer"))

	prom := prometheusmetrics.New()

	err := metrics.AddEngine(prom)
	if err != nil {
		return err
	}

	err = metrics.Hotload(prom.Name())
	if err != nil {
		return err
	}

	server := &http.Server{
		Addr:    listenAddr,
		Handler: metrics.Handler(),
	}

	return server.ListenAndServe()
}
