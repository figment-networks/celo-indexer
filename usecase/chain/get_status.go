package chain

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store"
)

type getStatusUseCase struct {
	db     *store.Store
	client figmentclient.Client
}

func NewGetStatusUseCase(db *store.Store, c figmentclient.Client) *getStatusUseCase {
	return &getStatusUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getStatusUseCase) Execute(ctx context.Context) (*DetailsView, error) {
	mostRecentSyncable, err := uc.db.Syncables.FindMostRecent()
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}

	chainStatus, err := uc.client.GetChainStatus(ctx)
	if err != nil {
		return nil, err
	}

	return ToDetailsView(mostRecentSyncable, chainStatus), nil
}
