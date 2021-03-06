package indexing

import (
	"context"
	"github.com/figment-networks/celo-indexer/store/psql"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/config"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

type BackfillCmdHandler struct {
	cfg    *config.Config
	db     *psql.Store
	client figmentclient.Client

	useCase *backfillUseCase
}

func NewBackfillCmdHandler(cfg *config.Config, db *psql.Store, c figmentclient.Client) *BackfillCmdHandler {
	return &BackfillCmdHandler{
		cfg:    cfg,
		db:     db,
		client: c,
	}
}

func (h *BackfillCmdHandler) Handle(ctx context.Context, parallel bool, force bool, targetIds []int64) {
	logger.Info("running backfill use case [handler=cmd]")

	useCaseConfig := BackfillUseCaseConfig{
		Parallel:  parallel,
		Force:     force,
		TargetIds: targetIds,
	}
	err := h.getUseCase().Execute(ctx, useCaseConfig)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *BackfillCmdHandler) getUseCase() *backfillUseCase {
	if h.useCase == nil {
		return NewBackfillUseCase(h.cfg, h.db, h.client)
	}
	return h.useCase
}
