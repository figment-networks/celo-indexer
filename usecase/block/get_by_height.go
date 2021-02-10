package block

import (
	"context"

	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store/psql"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	db     *psql.Store
	client figmentclient.ClientIface
}

func NewGetByHeightUseCase(db *psql.Store, c figmentclient.ClientIface) *getByHeightUseCase {
	return &getByHeightUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getByHeightUseCase) Execute(ctx context.Context, height *int64) (*DetailsView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.GetCore().Syncables.FindMostRecent()
	if err != nil {
		return nil, err
	}
	lastH := mostRecentSynced.Height

	// Show last synced height, if not provided
	if height == nil {
		height = &lastH
	}

	if *height > lastH {
		return nil, errors.New("height is not indexed yet")
	}

	block, err := uc.client.GetBlockByHeight(ctx, *height)
	if err != nil {
		return nil, err
	}

	return ToDetailsView(block), nil
}
