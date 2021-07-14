package indexing

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/types"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

var (
	_ types.WorkerHandler = (*fetchIdentitiesWorkerHandler)(nil)
)

type fetchIdentitiesWorkerHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *fetchIdentitiesUseCase
}

func NewFetchIdentitiesWorkerHandler(db *psql.Store, c figmentclient.Client) *fetchIdentitiesWorkerHandler {
	return &fetchIdentitiesWorkerHandler{
		db:     db,
		client: c,
	}
}

func (h *fetchIdentitiesWorkerHandler) Handle() {
	ctx := context.Background()

	logger.Info("running fetch identities use case [handler=worker]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *fetchIdentitiesWorkerHandler) getUseCase() *fetchIdentitiesUseCase {
	if h.useCase == nil {
		return NewFetchIdentitiesUseCase(h.db, h.client)
	}
	return h.useCase
}
