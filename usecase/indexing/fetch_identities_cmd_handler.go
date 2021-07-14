package indexing

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

type FetchIdentitiesCmdHandler struct {
	db     *psql.Store
	client figmentclient.Client

	useCase *fetchIdentitiesUseCase
}

func NewFetchIdentitiesCmdHandler(db *psql.Store, c figmentclient.Client) *FetchIdentitiesCmdHandler {
	return &FetchIdentitiesCmdHandler{
		db:     db,
		client: c,
	}
}

func (h *FetchIdentitiesCmdHandler) Handle(ctx context.Context) {
	logger.Info("running fetch identities use case [handler=cmd]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *FetchIdentitiesCmdHandler) getUseCase() *fetchIdentitiesUseCase {
	if h.useCase == nil {
		return NewFetchIdentitiesUseCase(h.db, h.client)
	}
	return h.useCase
}
