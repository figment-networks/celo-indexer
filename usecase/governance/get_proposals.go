package governance

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store"
)

type getProposalsUseCase struct {
	db *store.Store
	client figmentclient.Client
}

func NewGetProposalsUseCase(c figmentclient.Client, db *store.Store) *getProposalsUseCase {
	return &getProposalsUseCase{
		client: c,
		db: db,
	}
}

func (uc *getProposalsUseCase) Execute(ctx context.Context, cursor *int64, pageSize *int64) (*ProposalListView, error) {
	var limit int64
	if pageSize == nil {
		limit = 15
	}

	proposals, nextCursor, err := uc.db.ProposalAgg.All(limit, cursor)
	if err != nil {
		return nil, err
	}

	return ToProposalListView(proposals, nextCursor), nil
}
