package governance

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/theceloclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/figment-networks/celo-indexer/utils/logger"
)

type UpdateProposalsCmdHandler struct {
	db     *psql.Store
	client theceloclient.Client

	useCase *updateProposalsUseCase
}

func NewUpdateProposalsCmdHandler(db *psql.Store, c theceloclient.Client) *UpdateProposalsCmdHandler {
	return &UpdateProposalsCmdHandler{
		db:     db,
		client: c,
	}
}

func (h *UpdateProposalsCmdHandler) Handle(ctx context.Context) {
	logger.Info("running update proposals use case [handler=cmd]")

	err := h.getUseCase().Execute(ctx)
	if err != nil {
		logger.Error(err)
		return
	}
}

func (h *UpdateProposalsCmdHandler) getUseCase() *updateProposalsUseCase {
	if h.useCase == nil {
		return NewUpdateProposalsUseCase(h.client, h.db.GetGovernance().ProposalAgg)
	}
	return h.useCase
}

