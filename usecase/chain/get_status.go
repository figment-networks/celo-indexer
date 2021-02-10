package chain

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
)

type getStatusUseCase struct {
	db     *psql.Store
	client figmentclient.ClientIface
}

func NewGetStatusUseCase(db *psql.Store, c figmentclient.ClientIface) *getStatusUseCase {
	return &getStatusUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getStatusUseCase) Execute(ctx context.Context) (*DetailsView, error) {
	mostRecentSyncable, err := uc.db.GetCore().Syncables.FindMostRecent()
	if err != nil && err != psql.ErrNotFound {
		return nil, err
	}

	chainStatus, err := uc.client.GetChainStatus(ctx)
	if err != nil {
		return nil, err
	}

	return ToDetailsView(mostRecentSyncable, chainStatus), nil
}
