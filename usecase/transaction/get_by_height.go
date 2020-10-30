package transaction

import (
	"context"
	"github.com/figment-networks/celo-indexer/client/figmentclient"
	"github.com/figment-networks/celo-indexer/store"
	"github.com/pkg/errors"
)

type getByHeightUseCase struct {
	db     *store.Store
	client figmentclient.Client
}

func NewGetByHeightUseCase(db *store.Store, c figmentclient.Client) *getByHeightUseCase {
	return &getByHeightUseCase{
		db:     db,
		client: c,
	}
}

func (uc *getByHeightUseCase) Execute(ctx context.Context, height *int64) (*ListView, error) {
	// Get last indexed height
	mostRecentSynced, err := uc.db.Syncables.FindMostRecent()
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

	txs, err := uc.client.GetTransactionsByHeight(ctx, *height)
	if err != nil {
		return nil, err
	}

	return ToListView(txs), nil
}
