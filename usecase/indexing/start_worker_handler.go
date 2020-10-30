package indexing

import (
	"context"
	"fmt"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

var (
	_ types.WorkerHandler = (*runWorkerHandler)(nil)
)

type runWorkerHandler struct {
	cfg    *config.Config
	db     *store.Store
	client figmentclient.Client

	useCase *startUseCase
}

func NewRunWorkerHandler(cfg *config.Config, db *store.Store, c figmentclient.Client) *runWorkerHandler {
	return &runWorkerHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *runWorkerHandler) Handle() {
	batchSize := h.cfg.DefaultBatchSize
	ctx := context.Background()

	logger.Info(fmt.Sprintf("running indexer use case [handler=worker] [batchSize=%d]", batchSize))

	err := h.getUseCase().Execute(ctx, batchSize)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *runWorkerHandler) getUseCase() *startUseCase {
	if h.useCase == nil {
		return NewStartUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}


