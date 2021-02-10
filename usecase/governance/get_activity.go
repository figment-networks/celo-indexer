package governance

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getActivityUseCase struct {
	db     *psql.Store
	client figmentclient.ClientIface
}

func NewGetActivityUseCase(c figmentclient.ClientIface, db *psql.Store) *getActivityUseCase {
	return &getActivityUseCase{
		client: c,
		db:     db,
	}
}

func (uc *getActivityUseCase) Execute(ctx context.Context, proposalId uint64, pageSize *int64, cursor *int64) (*ActivityListView, error) {
	var limit int64
	if pageSize == nil {
		limit = 15
	}

	activities, nextCursor, err := uc.db.GetGovernance().GovernanceActivitySeq.FindByProposalId(proposalId, limit, cursor)
	if err != nil {
		return nil, err
	}

	return ToActivityListView(activities, nextCursor), nil
}
