package governance

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/theceloclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

var (
	_ types.WorkerHandler = (*updateProposalsWorkerHandler)(nil)
)

type updateProposalsWorkerHandler struct {
	cfg    *config.Config
	db     *store.Store
	client theceloclient.Client

	useCase *updateProposalsUseCase
}

func NewUpdateProposalsWorkerHandler(cfg *config.Config, db *store.Store, c theceloclient.Client) *updateProposalsWorkerHandler {
	return &updateProposalsWorkerHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *updateProposalsWorkerHandler) Handle() {
	ctx := context.Background()

	logger.Info("running update proposals use case [handler=worker]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *updateProposalsWorkerHandler) getUseCase() *updateProposalsUseCase {
	if h.useCase == nil {
		return NewUpdateProposalsUseCase(h.client, h.db)
	}
	return h.useCase
}

