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
	_ types.WorkerHandler = (*purgeWorkerHandler)(nil)
)

type purgeWorkerHandler struct {
	cfg    *config.Config
	db     *psql.Store
	client figmentclient.ClientIface

	useCase *purgeUseCase
}

func NewPurgeWorkerHandler(cfg *config.Config, db *psql.Store, c figmentclient.ClientIface) *purgeWorkerHandler {
	return &purgeWorkerHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *purgeWorkerHandler) Handle() {
	ctx := context.Background()

	logger.Info("running purge use case [handler=worker]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *purgeWorkerHandler) getUseCase() *purgeUseCase {
	if h.useCase == nil {
		return NewPurgeUseCase(h.cfg, h.db)
	}
	return h.useCase
}
