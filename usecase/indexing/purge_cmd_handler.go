package indexing

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

type PurgeCmdHandler struct {
	cfg    *config.Config
	db     *store.Store
	client figmentclient.Client

	useCase *purgeUseCase
}

func NewPurgeCmdHandler(cfg *config.Config, db *store.Store, c figmentclient.Client) *PurgeCmdHandler {
	return &PurgeCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *PurgeCmdHandler) Handle(ctx context.Context) {
	logger.Info("running purge use case [handler=cmd]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *PurgeCmdHandler) getUseCase() *purgeUseCase {
	if h.useCase == nil {
		return NewPurgeUseCase(h.cfg, h.db)
	}
	return h.useCase
}
