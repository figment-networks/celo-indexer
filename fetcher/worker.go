package fetcher

import (
	"net/http"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/indexer"
	"github.com/figment-networks/indexing-engine/datalake"
	"github.com/figment-networks/indexing-engine/worker"
	"github.com/rollbar/rollbar-go"
	"golang.org/x/net/websocket"
)

// Worker represents a fetcher worker
type Worker struct {
	cfg    *config.Config
	client figmentclient.Client
	dl     *datalake.DataLake
}

// NewWorker creates a fetcher worker
func NewWorker(cfg *config.Config, client figmentclient.Client, dl *datalake.DataLake) *Worker {
	return &Worker{
		cfg:    cfg,
		client: client,
		dl:     dl,
	}
}

// Run starts the fetcher worker
func (w *Worker) Run() error {
	server := http.Server{
		Addr:    w.cfg.FetchWorkerListenAddr(),
		Handler: websocket.Handler(w.handleConnection),
	}

	return server.ListenAndServe()
}

func (w *Worker) handleConnection(conn *websocket.Conn) {
	server := worker.NewWebsocketServer(conn)
	loop := worker.NewLoop(server)

	loop.Run(w.handleRequest)
}

func (w *Worker) handleRequest(req worker.Request) error {
	err := indexer.RunFetcherPipeline(req.Height, w.client, w.dl)
	if err != nil {
		rollbar.Error(err)
		return err
	}

	return nil
}
