package governance

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getProposalsUseCase struct {
	db     *psql.Store
	client figmentclient.ClientIface
}

func NewGetProposalsUseCase(c figmentclient.ClientIface, db *psql.Store) *getProposalsUseCase {
	return &getProposalsUseCase{
		client: c,
		db:     db,
	}
}

func (uc *getProposalsUseCase) Execute(ctx context.Context, cursor *int64, pageSize *int64) (*ProposalListView, error) {
	var limit int64
	if pageSize == nil {
		limit = 15
	}

	proposals, nextCursor, err := uc.db.GetGovernance().ProposalAgg.All(limit, cursor)
	if err != nil {
		return nil, err
	}

	return ToProposalListView(proposals, nextCursor), nil
}
