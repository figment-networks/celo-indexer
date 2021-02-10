package indexing

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

var (
	_ types.WorkerHandler = (*summarizeWorkerHandler)(nil)
)

type summarizeWorkerHandler struct {
	cfg    *config.Config
	db     *psql.Store
	client figmentclient.ClientIface

	useCase *summarizeUseCase
}

func NewSummarizeWorkerHandler(cfg *config.Config, db *psql.Store, c figmentclient.ClientIface) *summarizeWorkerHandler {
	return &summarizeWorkerHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *summarizeWorkerHandler) Handle() {
	ctx := context.Background()

	logger.Info("running summarize use case [handler=worker]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *summarizeWorkerHandler) getUseCase() *summarizeUseCase {
	if h.useCase == nil {
		return NewSummarizeUseCase(h.cfg, h.db)
	}
	return h.useCase
}
